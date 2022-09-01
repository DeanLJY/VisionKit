package app

import (
	"github.com/skyhookml/skyhookml/skyhook"
	"github.com/skyhookml/skyhookml/exec_ops"

	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func ToSkyhookInputDatasets(datasets map[string][]*DBDataset) map[string][]skyhook.Dataset {
	sk := make(map[string][]skyhook.Dataset)
	for name, dslist := range datasets {
		for _, ds := range dslist {
			sk[name] = append(sk[name], ds.Dataset)
		}
	}
	return sk
}

func ToSkyhookOutputDatasets(datasets map[string]*DBDataset) map[string]skyhook.Dataset {
	sk := make(map[string]skyhook.Dataset)
	for name, ds := range datasets {
		sk[name] = ds.Dataset
	}
	return sk
}

// Helper function to compute the keys already computed at a node.
// This only works for incremental nodes, which must produce the same keys across all output datasets.
func (node *DBExecNode) GetComputedKeys() map[string]bool {
	outputDatasets, _ := node.GetDatasets(false)
	outputItems := make(map[string][][]skyhook.Item)
	for name, ds := range outputDatasets {
		if ds == nil {
			return nil
		}
		var skItems []skyhook.Item
		for _, item := range ds.ListItems() {
			skItems = append(skItems, item.Item)
		}
		outputItems[name] = [][]skyhook.Item{skItems}
	}
	groupedItems := exec_ops.GroupItems(outputItems)
	keySet := make(map[string]bool)
	for key := range groupedItems {
		keySet[key] = true
	}
	return keySet
}

type ExecRunOptions struct {
	// If force, we run even if outputs were already available.
	Force bool

	// Whether to try incremental execution at this node.
	// If false, we throw error if parent datasets are not done.
	Incremental bool

	// If set, limit execution to these keys.
	// Only supported by incremental ops.
	LimitOutputKeys map[string]bool
}

// A RunData provides a Run function that executes a Runnable over the specified tasks.
type RunData struct {
	Name string
	Node skyhook.Runnable
	Tasks []skyhook.ExecTask

	// whether we'll be done with the node after running Tasks
	// i.e., whether Tasks contains all pending tasks at this node
	WillBeDone bool

	// job-related things to update
	JobOp *AppJobOp
	ProgressJobOp *ProgressJobOp

	// Saved error if any
	Error error
}

// Create a Job for this RunData and populate JobOp/ProgressJobOp.
func (rd *RunData) SetJob(name string, metadata string) {
	if rd.JobOp != nil {
		return
	}

	// initialize job
	// if the node doesn't provide a custom JobOp, we use "consoleprogress" view
	// otherwise the view for the job is the ExecOp's name
	opImpl := rd.Node.GetOp()
	nodeJobOp, nodeView := opImpl.GetJobOp(rd.Node)
	jobView := "consoleprogress"
	if nodeView != "" {
		jobView = nodeView
	}
	job := NewJob(
		fmt.Sprintf("Exec Node %s", name),
		"execnode",
		jobView,
		metadata,
	)

	rd.ProgressJobOp = &ProgressJobOp{}
	rd.JobOp = &AppJobOp{
		Job: job,
		TailOp: &skyhook.TailJobOp{},
		WrappedJobOps: map[string]skyhook.JobOp{
			"progress": rd.ProgressJobOp,
		},
	}
	if nodeJobOp != nil {
		rd.JobOp.WrappedJobOps["node"] = nodeJobOp
	}
	job.AttachOp(rd.JobOp)
}

// Update the AppJobOp with the saved error.
// We don't call this in RunData.Run by default because it's possible that the
// specified RunData.JobOp is shared across multiple Runs and shouldn't be
// marked as completed.
func (rd *RunData) SetDone() {
	if rd.Error == nil {
		rd.JobOp.SetDone(nil)
	} else {
		rd.JobOp.SetDone(rd.Error)
	}
}

// Prepare to run this node.
// Returns a RunData.
// Or error on error.
// Or nil RunData and error if the node is already done.
func (node *DBExecNode) PrepareRun(opts ExecRunOptions) (*RunData, error) {
	// create datasets for this op if needed
	outputDatasets, _ := node.GetDatasets(true)

	// if force, we clear the datasets first
	// otherwise, check if the datasets are done already
	if opts.Force {
		for _, ds := range outputDatasets {
			ds.Clear()
			ds.SetDone(false)
		}
	} else {
		done := true
		for _, ds := range outputDatasets {
			done = done && ds.Done
		}
		if done {
			return nil, nil
		}
	}

	// get parent datasets
	// for ExecNode parents, get computed dataset
	// in the future, we may need some recursive execution
	parentDatasets := make(map[string][]*DBDataset)
	parentsDone := true // whether parent datasets are fully computed
	for name, plist := range node.Parents {
		parentDatasets[name] = make([]*DBDataset, len(plist))
		for i, parent := range plist {
			if parent.Type == "n" {
				n := GetExecNode(parent.ID)
				dsList, _ := n.GetDatasets(false)
				ds := dsList[parent.Name]
				if ds == nil {
					return nil, fmt.Errorf("dataset for parent node %s[%s] is missing", n.Name, parent.Name)
				} else if !ds.Done && !opts.Incremental {
					return nil, fmt.Errorf("dataset for parent node %s[%s] is not done", n.Name, parent.Name)
				}
				parentDatasets[name][i] = ds
				parentsDone = parentsDone && ds.Done
			} else {
				parentDatasets[name][i] = GetDataset(parent.ID)
			}
		}
	}

	// get items in parent datasets
	items := make(map[string][][]skyhook.Item)
	for name, dslist := range parentDatasets {
		items[name] = make([][]skyhook.Item, len(dslist))
		for i, ds := range dslist {
			var skItems []skyhook.Item
			for _, item := range ds.ListItems() {
				skItems = append(skItems, item.Item)
			}
			items[name][i] = skItems
		}
	}

	// get tasks
	opImpl := node.GetOp()
	vnode := opImpl.Virtualize(no
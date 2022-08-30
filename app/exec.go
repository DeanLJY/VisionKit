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
	nodeJobOp, nodeView := opImpl
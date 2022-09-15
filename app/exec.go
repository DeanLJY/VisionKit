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
	vnode := opImpl.Virtualize(node.ExecNode)
	runnable := vnode.GetRunnable(ToSkyhookInputDatasets(parentDatasets), ToSkyhookOutputDatasets(outputDatasets))
	tasks, err := opImpl.GetTasks(runnable, items)
	if err != nil {
		return nil, err
	}

	// if running incrementally, remove tasks that were already computed
	// this is mostly so that we can see whether we will be done with this node after the current execution
	// (i.e., we are done here if parentsDone and we execute all remaining tasks)
	if opts.Incremental {
		var ntasks []skyhook.ExecTask
		completedKeys := node.GetComputedKeys()
		for _, task := range tasks {
			if completedKeys[task.Key] {
				continue
			}
			ntasks = append(ntasks, task)
		}
		tasks = ntasks
	}

	// limit tasks to LimitOutputKeys if needed
	// also determine whether this current execution will lead to all tasks being completed
	willBeDone := true
	if !parentsDone {
		willBeDone = false
	}
	if opts.LimitOutputKeys != nil {
		var ntasks []skyhook.ExecTask
		for _, task := range tasks {
			if !opts.LimitOutputKeys[task.Key] {
				continue
			}
			ntasks = append(ntasks, task)
		}
		if len(ntasks) != len(tasks) {
			tasks = ntasks
			willBeDone = false
		}
	}

	rd := &RunData{
		Name: node.Name,
		Node: runnable,
		Tasks: tasks,
		WillBeDone: willBeDone,
	}
	rd.SetJob(fmt.Sprintf("Exec Node %s", node.Name), fmt.Sprintf("%d", node.ID))
	return rd, nil
}

func (rd *RunData) Run() error {
	name := rd.Name

	// get container corresponding to rd.Node.Op
	log.Printf("[exec-node %s] [run] acquiring container", name)
	rd.JobOp.Update([]string{"Acquiring worker"})
	if err := AcquireWorker(rd.JobOp); err != nil {
		rd.Error = err
		return err
	}
	defer ReleaseWorker()
	containerInfo, err := AcquireContainer(rd.Node, rd.JobOp)
	if err != nil {
		rd.Error = err
		return err
	}
	log.Printf("[exec-node %s] [run] ... acquired container %s at %s", name, containerInfo.UUID, containerInfo.BaseURL)
	defer ReleaseWorker()

	// we want to de-allocate the container in two cases:
	// (1) when we return from this function
	// (2) if user requests to stop this job
	// we achieve this as follows:
	// - associate cleanup func with the JobOp
	// - on return, call AppJobOp.Cleanup to only de-allocate if it hasn't been de-allocated already
	// this is possible because AppJobOp will take care of unsetting CleanupFunc whenever it's called
	rd.JobOp.SetCleanupFunc(func() {
		err := skyhook.JsonPost(Config.WorkerURL, "/container/end", skyhook.EndRequest{containerInfo.UUID}, nil)
		if err != nil {
			log.Printf("[exec-node %s] [run] error ending exec container: %v", name, err)
		}
	})
	defer rd.JobOp.Cleanup()

	nthreads := containerInfo.Parallelism
	log.Printf("[exec-node %s] [run] running %d tasks in %d threads", name, len(rd.Tasks), nthreads)
	rd.ProgressJobOp.SetTotal(len(rd.Tasks))

	counter := 0
	var applyErr error
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < nthreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for !rd.JobOp.IsStopping() {
				// get next task
				mu.Lock()
				if counter >= len(rd.Tasks) || applyErr != nil {
					mu.Unlock()
					break
				}
				task := rd.Tasks[counter]
				counter++
				mu.Unlock()

				log.Printf("[exec-node %s] [run] apply on %s", name, task.Key)
				err := skyhook.JsonPost(containerInfo.BaseURL, "/exec/task", skyhook.ExecTaskRequest{task}, nil)

				if err != nil {
					mu.Lock()
					applyErr = err
					mu.Unlock()
					break
				}

				mu.Lock()
				rd.ProgressJobOp.Increment()
				rd.JobOp.Update([]string{fmt.Sprintf("finished applying on key [%s]", task.Key)})
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if applyErr != nil {
		rd.Error = applyErr
		return applyErr
	}

	// update dataset states
	if rd.WillBeDone {
		for _, ds := range rd.Node.OutputDatasets {
			(&DBDataset{Dataset: ds}).SetDone(true)
		}
	}

	log.Printf("[exec-node %s] [run] done", name)
	return nil
}

// Get some number of incremental outputs from this node.
type IncrementalOptions struct {
	// Number of random outputs to compute at this node.
	// Only one of Count or Keys should be specified.
	Count int
	// Compute outputs matching these keys.
	Keys []string
	// MultiExecJob to update during incremental execution.
	// For non-incremental ancestors, we pass this JobOp to RunNode.
	JobOp *MultiExecJobOp
}
func (node *DBExecNode) Incremental(opts IncrementalOptions) error {
	isIncremental := func(node *DBExecNode) bool {
		return node.GetOp().IsIncremental()
	}

	if !isIncremental(node) {
		return fmt.Errorf("can only incrementally run incremental nodes")
	} else if node.IsDone() {
		return nil
	}

	log.Printf("[exec-node %s] [incremental] begin execution", node.Name)
	// identify all non-incremental ancestors of this node
	// but stop the search at ExecNodes whose outputs have already been computed
	// we will need to run these ancestors in their entirety
	// note: we do not need to worry about Virtualize here because we assume Virtualize and Incremental are mutually exclusive
	var nonIncremental []*DBExecNode
	incrementalNodes := make(map[int]*DBExecNode)
	q := []*DBExecNode{node}
	seen := map[int]bool{node.ID: true}
	for len(q) > 0 {
		cur := q[len(q)-1]
		q = q[0:len(q)-1]

		if cur.IsDone() {
			continue
		}

		if !isIncremental(cur) {
			nonIncremental = append(nonIncremental, cur)
			continue
		}

		incrementa
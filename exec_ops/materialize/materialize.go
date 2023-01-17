package virtual_debug

// An identity function, but materializes its input datasets.
// If inputs have non-default provider, then we copy them into a new dataset.
// (If they have default provider, we still copy them.)
//

import (
	"github.com/skyhookml/skyhookml/skyhook"
	"github.com/skyhookml/skyhookml/exec_ops"

	"fmt"
)

func init() {
	skyhook.AddExecOpImpl(skyhook.ExecOpImpl{
		Config: skyhook.ExecOpConfig{
			ID: "materialize",
			Name: "Materialize",
			Description: "Materialize input datasets",
		},
		Inputs: []skyhook.ExecInput{{Name: "inputs", Variable: true}},
		GetOutputs: exec_ops.GetOutputsSimilarToInputs,
		Requirements: func(node skyhook.Runnable) map[string]int {
			return nil
		},
		GetTasks: func(node skyhook.Runnable, rawItems map[string][][]skyhook.Item) ([]skyhook.ExecTask, error) {
			items := rawItems["inputs"]
			var tasks []skyhook.ExecTask
			for i, itemList := range items {
				for _, item := range itemList {
					taskItems := make([][]skyhook.Item, len(items))
					taskItems[i] = []skyhook.Item{item}
					tasks = append(tasks, skyhook.ExecTask{
						Key: item.Key,
						Items: map[string
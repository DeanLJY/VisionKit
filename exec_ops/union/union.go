package union

import (
	"github.com/skyhookml/skyhookml/skyhook"

	"fmt"
	"strconv"
	urllib "net/url"
)

func init() {
	skyhook.AddExecOpImpl(skyhook.ExecOpImpl{
		Config: skyhook.ExecOpConfig{
			ID: "union",
			Name: "Union",
			Description: "Create an output dataset that includes all items from all input datasets",
		},
		Inputs: []skyhook.ExecInput{{Name: "inputs", Variable: true}},
		GetOutputs: func(params string, inputTypes map[string][]skyhook.DataType) []skyhook.ExecOutput {
			if len(inputTypes["inputs"]) == 0 {
				return nil
			}
			return []skyhook.ExecOutput{{Name: "output", DataType: inputTypes["inputs"][0]}}
		},
		Requirements: func(node skyhook.Runnable) map[string]int {
			return nil
		},
		GetTasks: func(node skyhook.Runnable, rawItems map[string][][]skyhook.Item) ([]skyhook.ExecTask, error) {
			// Create one task per item.
			// We set the Key of the task to the output key
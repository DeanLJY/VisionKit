package sample

import (
	"github.com/skyhookml/skyhookml/skyhook"
	"github.com/skyhookml/skyhookml/exec_ops"

	"encoding/json"
	"fmt"
	"math/rand"
	urllib "net/url"
)

type Params struct {
	// One of "count", "percentage", or "direct".
	Mode string
	// If mode=="count", the number of items to sample.
	Count int
	// If mode=="percentage", the percentage of items to sample.
	Percentage float64
	// If mode=="direct", the list of keys to sample.
	Keys []string
}

func init() {
	skyhook.AddExecOpImpl(skyhook.ExecOpImpl{
		Config: skyhook.ExecOpConfig{
			ID: "sample",
			Name: "Sample",
			Description: "Sample a subset of items from one or more datasets",
		},
		Inputs: []skyhook.ExecInput{{Name: "inputs", Variable: true}},
		GetOutputs: exec_ops.GetOutputsSimilarToInputs,
		Requirements: func(node skyhook.Runnable) map[string]int {
			return nil
		},
		GetTasks: func(node skyhook.Runnable, allItems map[string][][]skyhook
package concatenate

// Merge all input items into one output item.
// For table inputs, this is like SQL UNION operation, at least within one dataset.

import (
	"github.com/skyhookml/skyhookml/skyhook"
	"github.com/skyhookml/skyhookml/exec_ops"

	"io"
)

func init() {
	skyhook.AddExecOpImpl(skyhook.ExecOpImpl{
		Config: skyhook.ExecOpConfig{
			ID: "concatenate",
			Name: "Concatenate",
			Description: "Merge all items in t
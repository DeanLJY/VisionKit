package pytorch

import (
	"github.com/skyhookml/skyhookml/skyhook"
	"github.com/skyhookml/skyhookml/exec_ops"
	"github.com/skyhookml/skyhookml/exec_ops/python"

	"encoding/json"
	"fmt"
	"strconv"
)

func GetInferOutputs(params skyhook.PytorchInferParams) []skyhook.ExecOutput {
	var outputs []skyhook.ExecOutput
	for i, output := range params.OutputDatasets {
		outputs = append(outputs, skyhook.ExecOutput{
			Name: fmt.Sprintf("%d-%s", i, output.Layer),
			DataType: output.DataType,
		}
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
		})
	}
	return outputs
}

func Prepare(url string, node skyhook.Runnable) (skyhook.ExecOp, error) {
	// check the ArchID just to make sure we have all git repositories
	var params skyhook.PytorchInferParams
	if err := exec_ops.DecodeParams(node, &params, false); err != nil {
		return nil, err
	}
	_, components, err := GetTrainArgs(url, params.ArchID)
	if err != nil {
		return nil, err
	}
	if err := EnsureRepositories(components); err != nil {
		return nil, err
	}

	inputDatasets := node.InputDatasets

	paramsArg := node.Params
	cmd := skyhook.Command(
		fmt.Sprintf("pytorch-exec-%s", node.Name), skyhook.CommandOptions{},
		"python3", "exec_ops/pytorch/run.py",
		strconv.Itoa(inputDatasets["model"][0].ID), paramsArg,
	)

	var flatOutputs []skyhook.Dataset
	for _, output := range GetInferOutputs(params) {
		flatOutputs = append(flatOutputs, node.OutputDatasets[output.Name])
	}

	op, err := python.NewPythonOp(cmd, url, python.Params{}, inputDatasets["inputs"], flatOutputs)
	if err != nil {
		return nil, err
	}

	return op, nil
}

var InferImpl = skyhook.ExecOpImpl{
	Config: skyhook.ExecOpConfig{
		ID: "pytorch_infer",
		Name: "Pytorch (infer)",
		Description: "Pytorch (infer)",
	},
	Inputs: []skyhook.ExecInput{
		{Name: "inputs", Variable: true},
		{Name: "model", DataTypes: []skyhook.DataType{skyhook.FileType}},
	},
	GetOutputs: func(rawParams string, inputTypes map[string][]skyhook.Da
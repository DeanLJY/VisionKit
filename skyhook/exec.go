package skyhook

import (
	"strconv"
	"strings"
)

type ExecParent struct {
	// "n" for ExecNode, "d" for Dataset
	Type string
	ID int

	// name of ExecNode output that is being input
	Name string

	// the data type of this parent
	DataType DataType
}

func (p ExecParent) String() string {
	var parts []string
	parts = append(parts, p.Type)
	parts = append(parts, strconv.Itoa(p.ID))
	if p.Type == "n" {
		parts = append(parts, p.Name)
	}
	return strings.Join(parts, ",")
}

type ExecInput struct {
	Name string
	// nil if input can be any type
	DataTypes []DataType
	// true if this node can accept multiple inputs for this name
	Variable bool
}

type ExecOutput struct {
	Name string
	DataType DataType
}

type ExecNode struct {
	ID int
	Name string
	Op string
	Params string

	// currently configured parents for each input
	Parents map[string][]ExecParent
}

func (node ExecNode) GetOp() ExecOpProvider {
	return GetExecOp(node.Op)
}

func (node ExecNode) GetInputs() []ExecInput {
	return node.GetOp().GetInputs(node.Params)
}

func (node ExecNode) GetInputTypes() map[string][]DataType {
	inputTypes := make(map[string][]DataType)
	for _, input := range node.GetInputs() {
		for _, parent := range node.Parents[input.Name] {
			inputTypes[input.Name] = append(i
package archs

import (
	"github.com/skyhookml/skyhookml/skyhook"
	"github.com/skyhookml/skyhookml/exec_ops"
	"github.com/skyhookml/skyhookml/exec_ops/pytorch"
)

type Impl struct {
	ID string
	Name string
	TrainInputs []skyhook.ExecInput
	InferInputs []skyhook.ExecInput
	InferOutputs []skyhook.ExecOutput
	TrainPrepare func(skyhook.Runnable) (skyhook.PytorchTrainParams, error)
	InferPrepare func(skyhook.Runnable) (sk
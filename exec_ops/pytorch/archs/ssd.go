package archs

import (
	"github.com/skyhookml/skyhookml/skyhook"
	"github.com/skyhookml/skyhookml/exec_ops"
)

func init() {
	type TrainParams struct {
		skyhook.PytorchTrainParams
		Mode string
		ValPercent int
	}

	type InferParams struct {
		ConfidenceThreshold float64
	}

	type ModelParams struct {
		Mode string `json:"mode,omitempty"`
		ConfidenceThreshold float64 `json:"confidence_threshold,omitempty"`
		IouThreshold float64 `json:"iou_threshold,omitempty"`
	}

	AddImpl(Impl{
		ID: "pytorch_ssd",
		Name: "MobileNet+SSD",
		TrainInputs: []skyhook.ExecInput{
			{Name: "images", DataTypes: []skyh
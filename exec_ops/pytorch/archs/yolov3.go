package archs

import (
	"github.com/skyhookml/skyhookml/skyhook"
	"github.com/skyhookml/skyhookml/exec_ops"
)

func init() {
	type TrainParams struct {
		skyhook.PytorchTrainParams
		Mode string
		Resize skyhook.PDDImageOptions
		ValPercent int
	}

	type InferParams struct {
		Resize skyhook.PDDImageOptions
		ConfidenceThreshold float64
	}

	type ModelParams struct {
		Mode string `json:"mode,omitempty"`
		ConfidenceThreshold float64 `json:"confidence_threshold,omitempty"`
		IouThreshold float64 `json:"iou_threshold,omitempty"`
	}

	AddImpl(Impl{
		ID: "pytorch_yolov3",
		Name: "YOLOv3-pytorch",
		TrainInputs: []skyhook.ExecInput{
			{Name: "images", DataTypes: []skyhook.DataType{skyhook.ImageType}},
			{Name: "detections", DataTypes: []skyhook.DataType{skyhook.DetectionType}},
			{Name: "models", DataTypes: []skyhook.DataType{skyhook.FileType}},
		},
		InferInputs: []
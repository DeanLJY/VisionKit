
export default {
	"detection": {
		Name: "Object Detection",
		Help: "Train a model to detect bounding boxes of instances of one or more object categories in images.",
		Inputs: [{
			ID: "images",
			Name: "Images",
			DataType: "image",
			Help: "An image dataset containing example inputs.",
		}, {
			ID: "detections",
			Name: "Detection Labels",
			DataType: "detection",
			Help: "A detection dataset containing bounding box labels corresponding to each input image.",
		}],
		Defaults: {
			Model: 'pytorch_yolov5',
			Mode: 'l',
		},
		Models: {
			'pytorch_yolov3': {
				ID: 'pytorch_yolov3',
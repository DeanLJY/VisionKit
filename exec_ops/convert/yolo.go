
package convert

import (
	"github.com/skyhookml/skyhookml/exec_ops"
	"github.com/skyhookml/skyhookml/skyhook"

	"fmt"
	"path/filepath"
	"strings"
)

// Convert to and from YOLOv3 format.
// Skyhook inputs requires two datasets, one image and one detection.
// This format is a flat FileDataset with paired images and labels stored under same original filename.
// An obj.names file is also created for the category names.

func init() {
	imageSpec := skyhook.DataSpecs[skyhook.ImageType].(skyhook.ImageDataSpec)

	skyhook.AddExecOpImpl(skyhook.ExecOpImpl{
		Config: skyhook.ExecOpConfig{
			ID: "to_yolo",
			Name: "To YOLO",
			Description: "Convert from [image, detection] datasets to YOLO image/txt format",
		},
		Inputs: []skyhook.ExecInput{
			{Name: "images", DataTypes: []skyhook.DataType{skyhook.ImageType}},
			{Name: "detections", DataTypes: []skyhook.DataType{skyhook.DetectionType}},
		},
		Outputs: []skyhook.ExecOutput{{Name: "output", DataType: skyhook.FileType}},
		Requirements: func(node skyhook.Runnable) map[string]int {
			return nil
		},
		GetTasks: func(node skyhook.Runnable, rawItems map[string][][]skyhook.Item) ([]skyhook.ExecTask, error) {
			// we mostly use SimpleTasks, which creates a task for each corresponding image/detection pair between the input datasets
			// but we need to assign one task for writing the "obj.names" output
			// to assign it, we just set the metadata to "obj.names", which applyFunc below will check
			tasks, err := exec_ops.SimpleTasks(node, rawItems)
			if err != nil {
				return nil, err
			}
			if len(tasks) > 0 {
				tasks[0].Metadata =  "obj.names"
			}
			return tasks, nil
		},
		Prepare: func(url string, node skyhook.Runnable) (skyhook.ExecOp, error) {
			var params struct {
				Format string
				Symlink bool
			}
			if err := exec_ops.DecodeParams(node, &params, true); err != nil {
				return nil, err
			}
			if node.Params == "" {
				params.Format = "jpeg"
			}

			outDS := node.OutputDatasets["output"]
			applyFunc := func(task skyhook.ExecTask) error {
				inImageItem := task.Items["images"][0][0]
				inLabelItem := task.Items["detections"][0][0]

				// write the image
				// we produce a symlink if requested by the user and if the output format matches
				// if the output format doesn't match, we have to decode and re-encode the image
				outImageFormat := params.Format
				if outImageFormat == "" {
					outImageFormat = inImageItem.Format
				}
				outImageExt := imageSpec.GetExtFromFormat(outImageFormat)
				if outImageExt == "" {
					// unknown format...? just use the format as ext
					outImageExt = outImageFormat
				}
				outImageMetadata := skyhook.FileMetadata{
					Filename: task.Key+"."+outImageExt,
				}
				outImageItem, err := exec_ops.AddItem(url, outDS, task.Key+"-image", outImageExt, "", outImageMetadata)
				if err != nil {
					return err
				}
				err = inImageItem.CopyTo(outImageItem.Fname(), outImageFormat, params.Symlink)
				if err != nil {
					return err
				}
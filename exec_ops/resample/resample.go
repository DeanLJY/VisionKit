package resample

import (
	"github.com/skyhookml/skyhookml/skyhook"
	"github.com/skyhookml/skyhookml/exec_ops"

	"fmt"
	"runtime"
	"strconv"
	"strings"
	urllib "net/url"
)

type Params struct {
	Fraction string
}

func (params Params) GetFraction() [2]int {
	if !strings.Contains(params.Fraction, "/") {
		x, _ := strconv.Atoi(params.Fraction)
		return [2]int{x, 1}
	}
	parts := strings.Split(params.Fraction, "/")
	numerator, _ := strconv.Atoi(parts[0])
	denominator, _ := strconv.Atoi(parts[1])
	return [2]int{numerator, denominator}
}

type Resample struct {
	URL string
	Params Params
	Datasets map[string]skyhook.Dataset
}

func (e *Resample) Parallelism() int {
	// if we resample video, each ffmpeg runs with two threads
	return runtime.NumCPU()/2
}

func (e *Resample) Apply(task skyhook.ExecTask) error {
	fraction := e.Params.GetFraction()

	process := func(item skyhook.Item, dataset skyhook.Dataset) error {
		if item.Dataset.DataType == skyhook.VideoType {
			// all we need to do is update the framerate in the metadata

			metadata := item.DecodeMetadata().(skyhook.VideoMetadata)
			metadata.Framerate = [2]int{metadata.Framerate[0]*fraction[0], metadata.Framerate[1]*fraction[1]}

			fname := item.Fname()
			if fname != "" {
				// if the filename is available, we can produce output as a reference
				// with modified metadata to the original file
				return skyhook.JsonPostForm(e.URL, fmt.Sprintf("/datasets/%d/items", dataset.ID), urllib.Values{
					"key": {task.Key},
					"ext": {item.Ext},
					"format": {item.Format},
					"metadata": {string(skyhook.JsonMarshal(metadata))},
					"provider": {"reference"},
					"provider_info": {item.Fname()},
				}, nil)
			} else {
				// Filename should always be available since we shouldn't be loading video into m
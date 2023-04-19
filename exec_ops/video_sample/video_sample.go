package video_sample

import (
	"github.com/skyhookml/skyhookml/skyhook"
	"github.com/skyhookml/skyhookml/exec_ops"

	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"runtime"
)

type Params struct {
	Length int
	Count int

	// "random" or "uniform"
	Mode string
}

type VideoSample struct {
	URL string
	Params Params
	Datasets map[string]skyhook.Dataset
}

func (e *VideoSample) Parallelism() int {
	// each ffmpeg runs with two threads
	return runtime.NumCPU()/2
}

func (e *VideoSample) Apply(task skyhook.ExecTask) error {
	// Decode task metadata to get the samples we need to extract.
	var samples [][2]int
	skyhook.JsonUnmarshal([]byte(task.Metadata), &samples)

	log.Printf("extracting %d samples from %s", len(samples), task.Key)

	// Create map of where samples start.
	startToEnd := make(map[int][]int)
	for _, sample := range samples {
		startToEnd[sample[0]] = append(startToEnd[sample[0]], sample[1])
	}

	type ProcessingSample struct {
		Key string
		Start int
		End int
		Writers []skyhook.SequenceWriter
	}

	// Load input items, output datasets, and metadatas.
	var inputs []skyhook.Item
	var outputDatasets []skyhook.Dataset
	var metadatas []skyhook.DataMetadata

	inputs = append(inputs, task.Items["video"][0][0])
	outputDatasets = append(outputDatasets, e.Datasets["samples"])
	metadatas = append(metadatas, inputs[0].DecodeMetadata())

	for i, itemList := range task.Items["others"] {
		inputs = append(inputs, itemList[0])
		outputDatasets = append(outputDatasets, e.Datasets[fmt.Sprintf("others%d", i)])
		metadatas = append(metadatas, itemList[0].DecodeMetadata())
	}

	// Samples where we're currently in the middle of the intervals.
	processing := make(map[string]ProcessingSample)

	err := skyhook.PerFrame(inputs, func(pos int, datas []interface{}) error {
		// add segments that start at this frame to the processing set
		for _, end := range startToEnd[pos] {
			sampleKey := fmt.Sprintf("%s_%d_%d", task.Key, pos, end)
			if _, ok := processing[sampleKey]; ok {
				// duplicate interval
				continue
			}

			sample := ProcessingSample{
				Key: sampleKey,
				Start: pos,
				End: end,
				Writers: make([]skyhook.SequenceWriter, len(inputs)),
			}

			for i, ds := range outputDatasets {
				// Add an item to the dataset first.
				// To do so, we need to know the ext/format/metadata.
				// If input/output type match, then we can copy it from the input.
				// If they don't match (video input, image output), we handle the special case.
				var ext, format string
				var metadata skyhook.DataMetadata
				if inputs[i].Dataset.DataType == skyhook.VideoType && ds.DataType == skyhook.ImageType {
					ext = "jpg"
					format = "jpeg"
					metadata = skyhook.NoMetadata{}
				} else {
					metadata = metadatas[i]
					ext, format = inputs[i].DataSpec().GetDefaultExtAndFormat(datas[i], metadatas[i])
				}
				if ds.DataType == skyhook.VideoType {
					// For video outputs, update the Duration of the metadata so that it matches the sample duration.
					vmeta := metadata.(skyhook.VideoMetadata)
					vmeta.Duration = float64((end-pos)*vmeta.Framerate[1])/float64(vmeta.Framerate[0])
					metadata = vmeta
				}
				item, err := exec_ops.AddItem(e.URL, ds, sampleKey, ext, format, metadata)
				if err != nil {
					return err
				}
				sample
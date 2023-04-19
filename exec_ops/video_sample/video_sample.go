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

	lo
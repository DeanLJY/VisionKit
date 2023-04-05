package mask

import (
	"github.com/skyhookml/skyhookml/skyhook"
	"github.com/skyhookml/skyhookml/exec_ops"

	gomapinfer "github.com/mitroadmaps/gomapinfer/common"

	"fmt"
	"runtime"
)

type Params struct {
	Dims [2]int
	Padding int
}

type Mask struct {
	Params Params
	URL string
	OutputDataset skyhook.Dataset
}

func (e *Mask) Parallelism() int {
	return runtime.NumCPU()
}

// TODO: handle numCategories>256
func (e *Mask) renderFrame(dtype skyhook.DataType, data interface{}, metadata skyhook.DataMetadata, categoryMap map[string]int) ([]byte, error) {
	dims := e.Params.Dims
	padding := e.Params.Padding
	canvas := make([]byte, dims[0]*dims[1])

	fillRectangle := func(sx, sy, ex, ey, cls int) {
		sx = skyh
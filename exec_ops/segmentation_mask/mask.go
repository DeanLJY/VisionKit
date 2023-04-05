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
		sx = skyhook.Clip(sx-padding, 0, dims[0])
		sy = skyhook.Clip(sy-padding, 0, dims[1])
		ex = skyhook.Clip(ex+padding, 0, dims[0])
		ey = skyhook.Clip(ey+padding, 0, dims[1])
		for x := sx; x < ex; x++ {
			for y := sy; y < ey; y++ {
				canvas[y*dims[0] + x] = byte(cls)
			}
		}
	}

	// category string to ID
	getCategoryID := func(name string) int {
		if categoryMap[name] != 0 {
			return categoryMap[name]
		}

		// looks like the category string is not in the category list
		// if we are creating a two-category output, then that's okay, we can just set it to 1
		// otherwise we should return an error
		if len(categoryMap) == 1 {
			return 1
		}
		return -1
	}

	if dtype == skyhook.ShapeType {
		shapes := data.([][]skyhook.Shape)[0]
		shapeDims := metadata.(skyhook.ShapeMetadata).CanvasDims
		if shapeDims[0] == 0 {
			// if no dims set in data, assume it corresponds to output dims
			shapeDims = dims
		}
		for _, shape := range shapes {
			if shape.Type == skyhook.BoxShape {
				bounds := shape.Bounds()
				catID := getCategoryID(shape.Category)
				if catID == -1 {
					return nil, fmt.Errorf("unknown category %s", shape.Category)
				}
				fillRectangle(
					bounds[0]*dims[0]/shapeDims[0],
					bounds[1]*dims[1]/shapeDims[1],
					bounds[2]*dims[0]/shapeDims[0],
					bounds[3]*dims[1]/shapeDims[1],
					catID,
				)
			} else if shape.Type == skyhook.LineShape {
				sx := shape.P
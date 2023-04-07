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
				sx := shape.Points[0][0]*dims[0]/shapeDims[0]
				sy := shape.Points[0][1]*dims[1]/shapeDims[1]
				ex := shape.Points[1][0]*dims[0]/shapeDims[0]
				ey := shape.Points[1][1]*dims[1]/shapeDims[1]
				catID := getCategoryID(shape.Category)
				if catID == -1 {
					return nil, fmt.Errorf("unknown category %s", shape.Category)
				}
				for _, p := range gomapinfer.DrawLineOnCells(sx, sy, ex, ey, dims[0], dims[1]) {
					for ox := -padding; ox < padding; ox++ {
						for oy := -padding; oy < padding; oy++ {
							x := p[0]+ox
							y := p[1]+oy
							if x < 0 || x >= dims[0] || y < 0 || y >= dims[1] {
								continue
							}
							canvas[y*dims[0] + x] = byte(catID)
						}
					}
				}
			} else if shape.Type == skyhook.PolygonShape {
				catID := getCategoryID(shape.Category)
				if catID == -1 {
					return nil, fmt.Errorf("unknown category %s", shape.Category)
				}
				var polygon gomapinfer.Polygon
				for _, point := range shape.Points {
					polygon = append(polygon, gomapinfer.Point{
						float64(point[0]*dims[0]/shapeDims[0]),
						float64(point[1]*dims[1]/shapeDims[1]),
					})
				}
				bounds := shape.Bounds()

				sx := skyhook.Clip(bounds[0]*dims[0]/shapeDims[0], 0, dims[0])
				sy := skyhook.Clip(bounds[1]*dims[1]/shapeDims[1], 0, dims[1])
				ex := skyhook.Clip(bounds[2]*dims[0]/shapeDims[0], 0, dims[0])
				ey := skyhook.Clip(bounds[3]*dims[1]/shapeDims[1], 0, dims[1])
				for x := sx; x < ex; x++ {
					for y := sy; y < ey; y++ {
						if !polygon.Contains(gomapinfer.Point{float64(x), float64(y)}) {
							continue
						}
						canvas[y*dims[0] + x] = byte(catID)
					}
				}
			} else if shape.Type == skyhook.PointShape {
				catID := getCategoryID(shape.Category)
				if catID == -1 {
					return nil, fmt.Errorf("unknown category %s", shape.Category)
				}
				p := [2]int{
					shape.Points[0][0]*dims[0]/shapeDims[0],
					shape.Points[0][1]*dims[1]/shapeDims[1],
				}

				// Draw circle of radius padding centered at p.
				for ox := -padding; ox < padding; ox++ {
					for oy := -padding; oy < padding; oy++ {
						// Check radius.
						d := ox*ox+oy*oy
						if d > padding*padding {
							continue
						}
						// Set pixel.
						x := p[0]+ox
						y := p[1]+oy
						if x < 0 || x >= dims[0] || y < 0 || y >= dims[1] {
							continue
						}
						canvas[y*dims[0] + x] = byte(catID)
					}
				}
			} else {
				panic(fmt.Errorf("mask for shape type %s not implemented", shape.Type))
			}
		}
	} else if dtype == skyhook.DetectionType {
		detections := data.([][]skyhook.Detection)[0]
		detDims := metadata.(skyhook.DetectionMetadata).CanvasDims
		for _, d := range detections {
			if detDims[0] != 0 && detDims != dims {
	
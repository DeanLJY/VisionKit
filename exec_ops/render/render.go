package render

import (
	"github.com/skyhookml/skyhookml/skyhook"
	"github.com/skyhookml/skyhookml/exec_ops"

	"fmt"
	"runtime"
	"strconv"
)

var Colors = [][3]uint8{
	[3]uint8{255, 0, 0},
	[3]uint8{0, 255, 0},
	[3]uint8{0, 0, 255},
	[3]uint8{255, 255, 0},
	[3]uint8{0, 255, 255},
	[3]uint8{255, 0, 255},
	[3]uint8{0, 51, 51},
	[3]uint8{51, 153, 153},
	[3]uint8{102, 0, 51},
	[3]uint8{102, 51, 204},
	[3]uint8{102, 153, 204},
	[3]uint8{102, 255, 204},
	[3]uint8{153, 102, 102},
	[3]uint8{204, 102, 51},
	[3]uint8{204, 255, 102},
	[3]uint8{255, 255, 204},
	[3]uint8{121, 125, 127},
	[3]uint8{69, 179, 157},
	[3]uint8{250, 215, 160},
}


type Render struct {
	URL string
	Dataset skyhook.Dataset
}

func (e *Render) Parallelism() int {
	return runtime.NumCPU()
}

func renderFrame(dtypes []skyhook.DataType, datas []interface{}, metadatas []skyhook.DataMetadata) (skyhook.Image, error) {
	var canvas skyhook.Image
	var canvases []skyhook.Image
	for i, data := range datas {
		if dtypes[i] == skyhook.ImageType || dtypes[i] == skyhook.VideoType {
			canvas = data.([]skyhook.Image)[0].Copy()
			canvases = append(canvases, canvas)
			continue
		}

		if dtypes[i] == skyhook.IntType {
			x := data.([]int)[0]
			var text string
			categories := metadatas[i].(skyhook.IntMetadata).Categories
			if x >= 0 && x < len(categories) {
				text = categories[x]
			} else {
				text = strconv.Itoa(x)
			}
			canvas.DrawText(skyhook.RichText{Text: text})
		} else if dtypes[i] == skyhook.StringType {
			text := data.([]string)[0]
			canvas.DrawText(skyhook.RichText{Text: text})
		} else if dtypes[i] == skyhook.ShapeTy
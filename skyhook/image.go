
package skyhook

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"io"
	"os"
)

type Image struct {
	Width int
	Height int
	Bytes []byte
}

func NewImage(width int, height int) Image {
	return Image{
		Width: width,
		Height: height,
		Bytes: make([]byte, 3*width*height),
	}
}

func ImageFromBytes(width int, height int, bytes []byte) Image {
	return Image{
		Width: width,
		Height: height,
		Bytes: bytes,
	}
}

func ImageFromJPGReader(rd io.Reader) (Image, error) {
	im, err := jpeg.Decode(rd)
	if err != nil {
		return Image{}, err
	}
	return ImageFromGoImage(im), nil
}

func ImageFromPNGReader(rd io.Reader) (Image, error) {
	im, err := png.Decode(rd)
	if err != nil {
		return Image{}, err
	}
	return ImageFromGoImage(im), nil
}

func ImageFromGoImage(im image.Image) Image {
	rect := im.Bounds()
	width := rect.Dx()
	height := rect.Dy()
	bytes := make([]byte, width*height*3)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			r, g, b, _ := im.At(i + rect.Min.X, j + rect.Min.Y).RGBA()
			bytes[(j*width+i)*3+0] = uint8(r >> 8)
			bytes[(j*width+i)*3+1] = uint8(g >> 8)
			bytes[(j*width+i)*3+2] = uint8(b >> 8)
		}
	}
	return Image{
		Width: width,
		Height: height,
		Bytes: bytes,
	}
}

func ImageFromFile(fname string) (Image, error) {
	file, err := os.Open(fname)
	if err != nil {
		return Image{}, err
	}
package skyhook

import (
	"fmt"
	"io"
	"path/filepath"
)


type ImageDataSpec struct{}

func (s ImageDataSpec) DecodeMetadata(rawMetadata string) DataMetadata {
	return NoMetadata{}
}

type ImageStreamHeader struct {
	Width int
	Height int
	Channels int
	Length int
	BytesPerElement int
}

// Image is usually stored as Image but may become []Image since we support
// slice operations (so that operations can process Video/Image in the same way
// through SynchronizedReader).
// So this helper function tries both and returns just the Image.
func (s ImageDataSpec) getImage(data interface{}) Image {
	if image, ok := data.(Image); ok {
		return image
	}
	return data.([]Image)[0]
}

func (s ImageDataSpec) ReadStream(r io.Reader) (interface{}, error) {
	var header ImageStreamHeader
	if err := ReadJsonData(r, &header); err != nil {
		return nil, err
	}
	bytes := make([]byte, header.Width*header.Height*3)
	if _, err := io.ReadFull(r, bytes); err != nil {
		return nil, err
	}
	image := Image{
		Width: header.Width,
		Height: header.Height,
		Bytes: bytes,
	}
	return image, nil
}

func (s ImageDataSpec) WriteStream(data interface{}, w io.Writer) error {
	image := s.getImage(data)
	header := ImageStreamHeader{
		Width: image.Width,
		Height: image.Height,
		Channels: 3,
		Length: 1,
		BytesPerElement: len(image.Bytes),
	}
	if err := WriteJsonData(header, w); err != nil {
		return err
	}
	if _, err := w.Write(image.Bytes); err != nil {
		return err
	}
	return nil
}

func (s ImageDataSpec) Read(format string, metadata DataMetadata, r io.Reader) (data interface{}, err error) {
	var image Image
	if format == "jpeg" {
		image, err = ImageFromJPGReader(r)
	} else if format == "png" {
		image, err = ImageFromPNGReader(r)
	} else {
		err = fmt.Errorf("unknown format %s", format)
	}
	if err != nil {
		return nil, err
	}
	return im
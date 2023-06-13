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
// So 
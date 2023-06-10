package skyhook

import (
	"io"
	"io/ioutil"
	"path/filepath"
)

type FileMetadata struct {
	Filename string `json:",omitempty"`
}

func (m FileMetadata) Update(other DataMetadata) DataMetadata {
	other_ := other.(FileMetadata)
	if other_.Filename != "" {
		m.Filename = other_.Filename
	}
	return m
}

type FileTypeDataSpec struct{}

func (s FileTypeDataSpec) DecodeMetadata(rawMetadata string) DataMetadata {
	if rawMetadata == "" {
		return FileMetadata{}
	}
	var m FileMetadata
	JsonUnmarshal([]byte(rawMetadata), &m)
	return m
}

type FileStreamHeader struct {
	Length int
}

func (s FileTypeDataSpec) ReadStream(r io.Reader) (interface{}, error) {
	var header FileStreamHeader
	if err := ReadJsonData(r, &header); err != nil {
		return nil, err
	}
	bytes := make([]byte, header.Length)
	if _, err := io.ReadFull(r, bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

func (s FileTypeDataSpec) WriteStream(data interface{}, w io.Writer) error {
	bytes := data.([]byte)
	header := FileStreamHeader{
		Length: len(bytes),
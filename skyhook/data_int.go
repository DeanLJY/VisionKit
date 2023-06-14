package skyhook

import (
	"encoding/json"
)

type IntMetadata struct {
	Categories []string `json:",omitempty"`
}

func (m IntMetadata) Update(other DataMetadata) DataMetadata {
	other_ := other.(IntMetadata)
	if len(other_.Categories) > 0 {
		m.Categories = other_.Categories
	}
	return m
}

type IntJsonSpec struct {}

func (s IntJsonSpec) DecodeMetadata(rawMetadata string) DataMetadata {
	if rawMetadata == "" {
		return IntMetadata{}
	}
	var m IntMetadata
	JsonUnmarshal([]byte(rawMetadata), &m)
	return m
}

func (s IntJsonSpec) DecodeData(bytes []byte) (interface{}, error) {
	var data []int
	err := json.Unmarshal(bytes, &data)
	return data, err
}

func (s IntJsonSpec) GetEmptyMetadata() (metadata DataMetadata) {
	return IntMetadata{}
}

func (s IntJsonSpec) Length(data inte
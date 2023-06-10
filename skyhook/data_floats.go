package skyhook

import (
	"encoding/json"
)

type FloatJsonSpec struct {}

func (s FloatJsonSpec) DecodeMetadata(rawMetadata string) DataMetadata {
	return NoMetadata{}
}

func (s FloatJsonSpec) DecodeData(bytes []byte) (interface{}, error) {
	var data [][]float64
	err := json.Unmarshal(bytes, &data)
	return data, err
}

func (s FloatJsonSpec) GetEmptyMetadata() (metadata DataMetadata) {
	return NoMetadata{}
}

func (s FloatJsonSpec) Length(data interface{}) int {
	
package skyhook

import (
	"encoding/json"
	"math"
)

type DetectionMetadata struct {
	CanvasDims [2]int `json:",omitempty"`
	Categories []string `json:",omitempty"`
}

func (m DetectionMetadata) Update(other DataMetadata) DataMetadata {
	other_ := other.(DetectionMetadata)
	if other_.CanvasDims[0] > 0 {
		m.CanvasDims = other_.CanvasDims
	}
	if len(other_.Categories) > 0 {
		m.Categories = other_.Categories
	}
	return m
}

type Detection struct {
	Left int
	Top int
	Right int
	Bottom int

	// Optional
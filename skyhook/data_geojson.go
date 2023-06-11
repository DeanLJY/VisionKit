package skyhook

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/paulmach/go.geojson"
	gomapinfer "github.com/mitroadmaps/gomapinfer/common"
)

type GeoJsonData struct {
	Collection *geojson.FeatureCollection
}

func GetGeometryBbox(g *geojson.Geometry) gomapinfer.Rectangle {
	var bbox gomapinfer.Rectangle = gomapinfer.EmptyRectangle

	handlePointBBox := func(coordinate []float64) {
		p := gomapinfer.Point{coordinate[0], coordinate[1]}
		bbox = bbox.Extend(p)
	}
	handleLineStringBBox := func(coordinates [][]float64) {
		for _, coordinate := range coordinates {
			p := gomapinfer.Point{coordinate[0
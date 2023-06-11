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
			p := gomapinfer.Point{coordinate[0], coordinate[1]}
			bbox = bbox.Extend(p)
		}
	}
	handlePolygonBBox := func(coordinates [][][]float64) {
		// We do not support holes yet, so just use coordinates[0].
		// coordinates[0] is the exterior ring while coordinates[1:] specify
		// holes in the polygon that should be excluded.
		for _, coordinate := range coordinates[0] {
			p := gomapinfer.Point{coordinate[0], coordinate[1]}
			bbox = bbox.Extend(p)
		}
	}

	if g.Type == geojson.GeometryPoint {
		handlePointBBox(g.Point)
	} else if g.Type == geojson.GeometryMultiPoint {
		for _, coordinate := range g.MultiPoint {
			handlePointBBox(coordinate)
		}
	} else if g.Type == geojson.GeometryLineString {
		handleLineStringBBox(g.LineString)
	} else if g.Type == geojson.GeometryMultiLineString {
		for _, coordinates := range g.MultiLineString {
			handleLineStringBBox(coordinates)
		}
	} else if g.Type == geojson.GeometryPolygon {
		handlePolygonBBox(g.Polygon)
	} else if g.Type == geojson.GeometryMultiPolygon {
		for _, coordinates := range g.MultiPolygon {
			handlePolygonBBox(coordinates)
		}
	}

	return bbox
}

type GeoJsonDataSpec struct{}

func (s GeoJsonDataSpec) DecodeMetadata(rawMetadata string) DataMetadata {
	return NoMetadata{}
}

func (s GeoJsonDataSpec) ReadStream(r io.Reader) (interface{}, error) {
	var data *geojson.FeatureColl
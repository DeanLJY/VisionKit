package app

import (
	"github.com/skyhookml/skyhookml/skyhook"

	"fmt"
	"log"
	"math/rand"
	"strings"
)

type DBDataset struct {
	skyhook.Dataset
	Done bool
}
type DBAnnotateDataset struct {
	skyhook.AnnotateDataset
	loaded bool
	InputDatasets []skyhook.Dataset
}
type DBItem struct {
	skyhook.Item
	loaded bool
}

const DatasetQuery = "SELECT id, name, type, data_type, metadata, hash, done FROM datasets"

func datasetListHelper(rows *Rows) []*DBDataset {
	datasets := []*DBDataset{}
	for rows.Next() {
		var ds DBDataset
		rows.Scan(&ds.ID, &ds.Name, &ds.Type, &ds.DataType, &ds.Metadata, &ds.Hash, &ds.Done)
		datasets = append(datasets, &ds)
	}
	return datasets
}

func ListDatasets() []*DBDataset {
	rows := db.Query(DatasetQuery)
	return datasetListHelper(rows)
}

func GetDataset(id int) *DBDataset {
	rows := db.Query(Datase
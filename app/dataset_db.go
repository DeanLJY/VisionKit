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
	rows := db.Query(DatasetQuery + " WHERE id = ?", id)
	datasets := datasetListHelper(rows)
	if len(datasets) == 1 {
		return datasets[0]
	} else {
		return nil
	}
}

func FindDataset(hash string) *DBDataset {
	rows := db.Query(DatasetQuery + " WHERE hash = ?", hash)
	datasets := datasetListHelper(rows)
	if len(datasets) == 1 {
		return datasets[0]
	} else {
		return nil
	}
}

const AnnotateDatasetQuery = "SELECT a.id, d.id, d.name, d.type, d.data_type, a.inputs, a.tool, a.params FROM annotate_datasets AS a LEFT JOIN datasets AS d ON a.dataset_id = d.id"

func annotateDatasetListHelper(rows *Rows) []*DBAnnotateDataset {
	annosets := []*DBAnnotateDataset{}
	for rows.Next() {
		var s DBAnnotateDataset
		var inputsRaw string
		rows.Scan(&s.ID, &s.Dataset.ID, &s.Dataset.Name, &s.Dataset.Type, &s.Dataset.DataType, &inputsRaw, &s.Tool, &s.Params)
		skyhook.JsonUnmarshal([]byte(inputsRaw), &s.Inputs)
		if s.Inputs == nil {
			s.Inputs = []skyhook.ExecParent{}
		}
		annosets = append(annosets, &s)
	}
	return annosets
}

func ListAnnotateDatasets() []*DBAnnotateDataset {
	rows := db.Query(AnnotateDatasetQuery)
	return annotateDatasetListHelper(rows)
}

func GetAnnotateDataset(id int) *DBAnnotateDataset {
	rows := db.Query(AnnotateDatasetQuery + " WHERE a.id = ?", id)
	annosets := annotateDatasetListHelper(rows)
	if len(annosets) == 1 {
		return annosets[0]
	} else {
		return nil
	}
}

func (s *DBAnnotateDataset) Load() {
	if s.loaded {
		return
	}

	s.Dataset = GetDataset(s.Dataset.ID).Dataset
	s.InputDatasets = make([]skyhook.Dataset, len(s.Inputs))
	for i, input := range s.Inputs {
		ds, err := ExecParentToDataset(input)
		if err != nil {
			continue
		}
		s.InputDatasets[i] = ds.Dataset
	}
	s.loaded = true
}

// samples a key that is present in all input datasets but not yet labeled in this annotate dataset
// TODO: have sampler object so that hash tables can be stored in memory instead of loaded from db each time
func (s *DBAnnotateDataset) SampleMissingKey() string {
	var keys map[string]bool
	for _, parent := range s.Inputs {
		ds, err := ExecParentToDataset(parent)
		if err != nil {
			// TODO: probably want to handle this error somehow
			continue
		}
		items := ds.ListItems()
		curKeys := make(map[string]bool)
		for _, item := range items {
			curKeys[item.Key] = true
		}
		if keys == nil {
			keys = curKeys
		} else {
			for key := range keys {
				if !curKeys[key] {
					delete(keys, key)
				}
			}
		}
	}

	items := (&DBDataset{Dataset: s.Dataset}).ListItems()
	for _, item := range items {
		delete(keys, item.Key)
	}

	var keyList []string
	for key := range keys {
		keyList = append(keyList, key)
	}
	if len(keyList) == 0 {
		return ""
	}
	return keyList[rand.Intn(len(keyList))]
}

type AnnotateDatasetUpdate struct {
	Tool *string
	Params *string
}

func (s *DBAnnotateDataset) Update(req AnnotateDatasetUpdate) {
	if req.Tool != nil {
		db.Exec("UPDATE annotate_datasets SET tool = ? WHERE id = ?", *req.Tool, s.ID)
	}
	if req.Params != nil {
		db.Exec("UPDATE annotate_datasets SET params = ? WHERE id = ?", *req.Params, s.ID)
	}
}

func (s *DBAnnotateDataset) Delete() {
	db.Exec("DELETE FROM annotate_datasets WHERE id = ?", s.ID)
}

const ItemQuery = "SELECT k, ext, format, metadata, provider, provider_info FROM items"

func itemListHelper(rows *Rows) []*DBItem {
	var items []*DBItem
	for rows.Next() {
		var item DBItem
		rows.Scan(&item.Key, &item.Ext, &item.Format, &item.Metadata, &item.Provider, &item.ProviderInfo)
		items = append(items, &item)
	}
	return items
}

func (ds *DBDataset) getDB() *Database {
	return GetCachedDB(ds.DBFname(), func(db *Database) {
		db.Exec(`CREATE TABLE IF NOT EXISTS items (
			-- item key
			k TEXT PRIMARY KEY,
			ext TEXT,
			format TEXT,
			metadata TEXT,
			-- set if LoadData call should go through non-default method, else NULL
			provider TEXT,
			provider_info TEXT
		)`)
		db.Exec(`CREATE TABLE IF NOT EXISTS datasets (
			id INTEGER PRIMARY KEY ASC,
			name TEXT,
			-- 'data' or 'computed'
			type TE
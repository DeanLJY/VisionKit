package skyhook

import (
	"crypto/sha256"
	"fmt"
	"os"
)

type Dataset struct {
	ID int
	Name string

	// data or computed
	Type string

	DataType DataType
	Metadata string

	// nil unless Type=computed
	Hash *string
}

type Item struct {
	Dataset Dataset
	Key string
	Ext string
	Format string
	Metadata string

	// nil to use default storage provider for LoadData / UpdateData
	Provider *string
	ProviderInfo *string
}

func (ds Dataset) Dirname() string {
	return fmt.Sprintf("data/items/%d", ds.ID)
}

func (ds Dataset) Mkdir() {
	os.Mkdir(ds.Dirname(), 0755)
}

func (ds Dataset) DataSpec() DataSpec {
	return DataSpecs[ds.DataType]
}

func (ite
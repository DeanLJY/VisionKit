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

func (item Item) DataSpec() DataSpec {
	return item.Dataset.DataSpec()
}

func (item Item) Fname() string {
	provider := item.GetProvider()
	if provider.Fname == nil {
		return ""
	}
	return provider.Fname(item)
}

func (item Item) GetProvider() ItemProvider {
	if item.Provider == nil {
		return DefaultItemProvider
	} else {
		return ItemProviders[*item.Provider]
	}
}

func (item Item) UpdateData(data interface{}, metadata DataMetadata) error {
	provider := item.GetProvider()
	if provider.UpdateData == nil {
		panic(fmt.Errorf("UpdateData not supported in dataset %s", item.Dataset.Name))
	}
	return provider.UpdateData(item, data, metadata)
}

func (item Item) LoadData() (interface{}, DataMetadata, error) {
	return item.GetProvider().LoadData(item)
}

func (item Item) LoadReader() (SequenceReader, DataMetadata) {
	metadata := item.DecodeMetadata()
	spec, ok := item.DataSpec().(SequenceDataSpec)
	if !ok {
		return ErrorSequenceReader{fmt.Errorf("data type %s is not sequence type", item.Dataset.DataType)}, metadata
	}

	fname := item.Fname()
	if fname == "" {
		// Since file is not available, we need to load the data and then return SliceReader.
		data, _, err := item.LoadData()
		if err != nil {
			return ErrorSequenceReader{err}, metadata
		}
		return &SliceReader{
			Data: data,
			Spec: spec,
		}, metadata
	}

	if fileSpec, fileOK := sp
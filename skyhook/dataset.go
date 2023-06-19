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

	if fileSpec, fileOK := spec.(FileSequenceDataSpec); fileOK {
		return fileSpec.FileReader(item.Format, metadata, fname), metadata
	}
	file, err := os.Open(fname)
	if err != nil {
		return ErrorSequenceReader{err}, metadata
	}
	return ClosingSequenceReader{
		Reader: spec.Reader(item.Format, metadata, file),
		ReadCloser: file,
	}, nil
}

func (item Item) LoadWriter() SequenceWriter {
	metadata := item.DecodeMetadata()
	spec, ok := item.DataSpec().(SequenceDataSpec)
	if !ok {
		return ErrorSequenceWriter{fmt.Errorf("data type %s is not sequence type", item.Dataset.DataType)}
	}

	item.Dataset.Mkdir()
	fname := item.Fname()
	if fileSpec, fileOK := spec.(FileSequenceDataSpec); fileOK {
		return fileSpec.FileWriter(item.Format, metadata, fname)
	}
	file, err := os.Create(fname)
	if err != nil {
		return ErrorSequenceWriter{err}
	}
	return ClosingSequenceWriter{
		Writer: spec.Writer(item.Format, metadata, file),
		WriteCloser: file,
	}
}

func (ds Dataset) Remove() {
	os.RemoveAll(fmt.Sprintf("data/items/%d", ds.ID))
}

func (item Item) Remove() {
	fname := item.Fname()
	if fname == "" {
		panic(fmt.Errorf("Remove not supported in dataset %s", item.Dataset.Name))
	}
	os.Remove(fname)
}

func (item Item) DecodeMetadata() DataMetadata {
	spec := item.DataSpec()
	metadata := spec.DecodeMetadata(item.Dataset.Metadata)
	metadata = metadata.Update(spec.DecodeMetadata(item.
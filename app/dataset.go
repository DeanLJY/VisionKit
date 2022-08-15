package app

import (
	"github.com/skyhookml/skyhookml/skyhook"

	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (item *DBItem) Handle(format string, w http.ResponseWriter, r *http.Request) {
	item.Load()

	fname := item.Fname()
	if format == item.Format && fname != "" {
		http.ServeFile(w, r, fname)
		return
	}

	data, metadata, err := item.LoadData()
	if err != nil {
		panic(err)
	}

	if format == "jpeg" {
		w.Header().Set("Content-Type", "image/jpeg")
	} else if format == "png" {
		w.Header().Set("Content-Type", "image/png")
	} else if format == "mp4" {
		w.Header().Set("Content-Type", "video/mp4")
	} else if format == "json" {
		w.Header().Set("Content-Type", "application/json")
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
		var filename string
		if item.Dataset.DataType == skyhook.FileType {
			filename = metadata.(skyhook.FileMetadata).Filename
		} else {
			ext := skyhook.GetExtFromFormat(item.Dataset.DataType, format)
			if ext == "" {
				ext = item.Ext
			}
			filename = item.Key + "." + ext
		}
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=\"%s\"", filename))
	}

	if err := item.DataSpec().Write(data, format, metadata, w); err != nil {
		panic(err)
	}
}

func init() {
	Router.HandleFunc("/datasets", func(w http.ResponseWriter, r *http.Request) {
		skyhook.JsonResponse(w, ListDatasets())
	}).Methods("GET")

	Router.HandleFunc("/datasets", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		name := r.PostForm.Get("name")
		dataType := r.PostForm.Get("data_type")
		ds := NewDataset(name, "data", skyhook.DataType(dataType), nil)
		skyhook.JsonResponse(w, ds)
	}).Methods("POST")

	Router.HandleFunc("/datasets/{ds_id}", func(w http.ResponseWriter, r *http.Request) {
		dsID := skyhoo
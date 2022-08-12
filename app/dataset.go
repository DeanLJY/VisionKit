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

	if fo
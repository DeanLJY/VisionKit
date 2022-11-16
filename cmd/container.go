package main

import (
	"github.com/skyhookml/skyhookml/skyhook"

	_ "github.com/skyhookml/skyhookml/ops"

	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	var coordinatorURL string
	var execOp skyhook.ExecOp

	var bindAddr string = ":8080"
	if len(os.Args) >= 2 {
		bindAddr = os.Args[1]
	}

	http.HandleFunc("/exec/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(404)
			return
		}

		var request skyhook.ExecBeginRequest
		if err := skyhook.ParseJsonRequest(w, r, &request); err !
package app

import (
	"github.com/skyhookml/skyhookml/skyhook"
	"github.com/skyhookml/skyhookml/exec_ops"

	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func ToSkyhookInputDatasets(datasets map[string][]*DBDataset) map[string][]skyhook.Dataset {
	sk := make(map[string][]skyhook.Dataset)
	for name, dslist := range datasets {
		for _, ds := range dslist {
			sk[name] = append(sk[name], ds.Dataset)
		}
	}
	return sk
}

func ToSkyhookOutputDatasets(datasets map[string]*DBDataset) map[string]skyhook.Dataset {
	sk := make(map[string]skyhook.Dataset)
	for name, ds := range datasets {
		sk[name] = ds.Dataset
	}
	return sk
}

// Helper function to compute the keys already computed at a node.
// This only works for incremental nodes, which must produce the same keys across all output datasets.
func (node *DBExecNode) GetComputedKeys() map[string]bool {
	outputDatasets, _ := node.GetDatasets(false)
	outputItems := make(map[string][][]skyhook.Item)
	for name, ds := range outputDatasets {
		if ds == nil {
			return nil
		}
		var skItems []skyhook.Item
		for _, item := range ds.ListItems() {
			skItems = append(skItems, item.Item)
		}
		outputItems[name] = [][]skyhook.It
package app

import (
	"github.com/skyhookml/skyhookml/skyhook"

	"fmt"
	"log"
	"strings"
)

type DBExecNode struct {
	skyhook.ExecNode
	Workspace string
}

const ExecNodeQuery = "SELECT id, name, op, params, parents, workspace FROM exec_nodes"

func execNodeListHelper(rows *Rows) []*DBExecNode {
	nodes := []*DBExecNode{}
	for rows.Next() {
		var node DBExecNode
		var parentsRaw string
		rows.Scan(&node.ID, &node.Name, &node.Op, &node.Params, &parentsRaw, &node.Workspace)
		skyhook.JsonUnmarshal([]byte(parentsRaw), &node.Parents)
		if node.Parents == nil {
			node.Parents = make(map[string][]skyhook.ExecParent)
		}

		// make sure parents list is set for each input
		for _, input := range node.GetInputs() {
			if node.Parents[input.Name] != nil {
				continue
			}
			node.Parents[input.Name] = []skyhook.ExecParent{}
		}

		nodes = append(nodes, &node)
	}
	return nodes
}

func ListExecNodes() []*DBExecNode {
	rows := db.Query(ExecNodeQuery)
	return execNodeListHelper(rows)
}

func (ws DBWorkspace) ListExecNodes() []*DBExecNode {
	rows := db.Query(ExecNodeQuery + " WHERE workspace = ?", ws)
	return execNodeListHelper(rows)
}

func GetExecNode(id int) *DBExecNode {
	rows := db.Query(ExecNodeQuery + " WHERE id = ?", id)
	nodes := execNodeListHelper(rows)
	if len(nodes) == 1 {
		return nodes[0]
	} else {
		return nil
	}
}

func NewExecNode(name string, op string, params string, parents map[string][]skyhook.ExecParent, workspace string) *DBExecNode {
	res := db.Exec(
		"INSERT INTO exec_nodes (name, op, params, parents, workspace) VALUES (?, ?, ?, ?, ?)",
		name, op, params,
		string(skyhook.JsonMarshal(parents)),
		workspace,
	)
	node := GetExecNode(res.LastInsertId())
	return node
}

func (node *DBExecNode) DatasetRefs() []int {
	var ds []int
	rows := db.Query("SELECT dataset_id FROM exec_ds_refs WHERE node_id = ?", node.ID)
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ds = append(ds, id)
	}
	return ds
}

// Get datasets for each output of this node.
// If create=true, creates new datasets to cover missing ones.
// Also returns bool, which is true if all datasets exist.
func (node *DBExecNode) GetDatasets(create bool) (map[string]*DBDataset, bool) {
	nodeHash := node.Hash()

	// remove references to datasets that don't even start with the nodeHash
	existingDS := node.DatasetRefs()
	for _, id := range existingDS {
		ds := GetDataset(id)
		if !strings.HasPrefix(*ds.Hash, nodeHash) {
			ds.DeleteExecRef(node.ID)
		}
	}

	// find datasets that match current hash
	datasets := make(map[string]*DBDataset)
	ok := true
	for _, output := range node.GetOutputs() {
		dsName := fmt.Sprintf("%s[%s]", node.Name, output.Name)
		curHash := fmt.Sprintf("%s[%s]", nodeHash, output.Name)
		ds := FindDataset(curHash)
		if ds == nil {
			ok = false
			if create {
				ds = NewDataset(dsName, "computed", output.DataType, &curHash)
			}
		}

		if ds != nil {
			ds.AddExecRef(node.ID)
			datasets[output.Name] = ds
		} else {
			datasets[output.Name] = nil
		}
	}

	return datasets, ok
}

// Get dataset for a virtual
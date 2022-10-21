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
	node := GetExecNode(res.La
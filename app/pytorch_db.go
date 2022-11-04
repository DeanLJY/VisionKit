package app

import (
	"github.com/skyhookml/skyhookml/skyhook"
)

type DBPytorchComponent struct {skyhook.PytorchComponent}
type DBPytorchArch struct {skyhook.PytorchArch}

const PytorchComponentQuery = "SELECT id, params FROM pytorch_components"

func pytorchComponentListHelper(rows *Rows) []*DBPytorchComponent {
	var comps []*DBPytorchComponent
	for rows.Next() {
		var c DBPytorchComponent
		var paramsRaw string
		rows.Scan(&c.ID, &paramsRaw)
		skyhook.JsonUnmarshal([]byte(paramsRaw), &c.Params)
		comps = append(comps, &c)
	}
	return comps
}

func ListPytorchComponents() []*DBPytorchComponent {
	rows := db.Query(PytorchComponentQuery)
	return pytorchComponentListHelper(rows)
}

func GetPytorchComponent(id string) *DBPytorchComponent {
	rows := db.Query(PytorchComponentQuery + " WHERE id = ?", id)
	comps := pytorchComponentListHelper(rows)
	if len(comps) == 1 {
		return comps[0]
	} else {
		return nil
	}
}

func NewPytorchComponent(id string) *DBPytorchComponent {
	db.Exec("INSERT INTO pytorch_components (id, params) VALUES (?, '{}')", id)
	return GetPytorchComponent(id)
}

type PytorchComponentUpdate struct {
	Params *skyhook.PytorchComponentParams
}

func (c *DBPytorchComponent) Update(req PytorchComponentUpdate) {
	if req.Params != nil {
		db.Exec("UPDATE pytorch_components SET params = ? WHERE id = ?", string(skyhook.JsonMarshal(*req.Params)), c.ID)
	}
}

func (c *DBPytorchComponent) Delete() {
	db.Exec("DE
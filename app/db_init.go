package app

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Initialize the database on startup with cleanup operations.
// If init is true, we also first initialize the schema and populate certain tables.
func InitDB(init bool) {
	if init {
		db.Exec(`CREATE TABLE IF NOT EXISTS kv (
			k TEXT PRIMARY KEY,
			v TEXT
		)`)
		db.Exec(`CREATE TABLE IF NOT EXISTS datasets (
			id INTEGER PRIMARY KEY ASC,
			name TEXT,
			-- 'data' or 'computed'
			type TEXT,
			data_type TEXT,
			metadata TEXT DEFAULT '',
			-- only set if computed

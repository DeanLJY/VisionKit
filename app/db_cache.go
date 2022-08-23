package app

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"
)

// Database cache for per-dataset sqlite3 files.

var dbCache = make(map[string]*Database)
var dbCacheMu sync.Mutex

// Get a cached database connection to the specified sqlite3 file.
// If initFunc is set, we create the sqlite3 if it doesn't already exist, and
//   call initFunc each time a new con
package gostp

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"

	//DB-driver included in handlres
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// WorkDir - is a directory address where program has been launched
var WorkDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

// Db - opened database
var Db, Err = gorm.Open("sqlite3", filepath.Join(WorkDir, "event.db"))

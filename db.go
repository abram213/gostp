package gostp

import (
	"path/filepath"

	"github.com/jinzhu/gorm"

	//DB-driver included in handlres
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Db - opened database
var Db, Err = gorm.Open("sqlite3", filepath.Join(WorkDir, "event.db"))

package gostp

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"

	//DB-driver included in handlres
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Db - gorm db
var Db *gorm.DB

// Err - error variable
var Err error

// FunctionsMap - map of functions
var FunctionsMap = make(map[string]interface{})

// RegexMap - map of regexes
var RegexMap = make(map[string]RegexAndDescription)

// RegexAndDescription struct which contains regexes and description of error
type RegexAndDescription struct {
	Regex       string `yaml:"regex"`       // Regex - regular expression rule to validate string
	Description string `yaml:"description"` // Description - description why string is not valid
}

// Gostp - structure for gostp settings
type Gostp struct {
	WorkDir     string `yaml:"work_dir"`     // WorkDir - is a directory address where program has been launched (default - directory where program stored)
	SigningKey  string `yaml:"signing_key"`  // SigningKey - key for signing JWT (default - "")
	SQLtype     string `yaml:"sql_type"`     // SQLtype - type of gorm SQL (default - "sqlite3")
	SQLfilename string `yaml:"sql_filename"` // SQLfilename - filename of sqlite db (default "app.db")
	SQLhost     string `yaml:"sql_host"`     // SQLhost - host of remote or local db (defatult - "127.0.0.1")
	SQLport     string `yaml:"sql_port"`     // SQLport - port of remote or local db (default - "5432")
	SQLdbname   string `yaml:"sql_dbname"`   // SQLdbname - database name of remote or local db (default - "app")
	SQLuser     string `yaml:"sql_user"`     // SQLuser - database username (default - "admin")
	SQLpassword string `yaml:"sql_password"` // SQLpassword - database password (default - "admin")
	SQLsslmode  string `yaml:"sql_sslmode"`  // SQLsslmode - database sslmode (default - "disabled")
}

// Settings - main settings
var Settings = Gostp{
	CurrentFolder(),
	"",
	"sqlite3",
	"app.db",
	"127.0.0.1",
	"5432",
	"app",
	"admin",
	"admin",
	"disable"}

// Init - initialize of gostp
func Init(functionsMap map[string]interface{}, regexMap map[string]RegexAndDescription) {
	// Functions Map initialization
	FunctionsMap = functionsMap
	// Regex Map initialization
	RegexMap = regexMap

	if FileNotExist(filepath.Join(Settings.WorkDir, "env.yaml")) == nil {
		file, _ := ioutil.ReadFile(filepath.Join(Settings.WorkDir, "env.yaml"))
		err := yaml.Unmarshal(file, &Settings)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	if Settings.SQLtype == "sqlite3" {
		Db, Err = gorm.Open(Settings.SQLtype, filepath.Join(Settings.WorkDir, Settings.SQLfilename))
	} else {
		Db, Err = gorm.Open("postgres", "host="+Settings.SQLhost+" port="+Settings.SQLport+" user="+Settings.SQLuser+" dbname="+Settings.SQLdbname+" password="+Settings.SQLpassword+" sslmode="+Settings.SQLsslmode+"")
	}
}

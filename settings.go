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

var Db *gorm.DB

var Err error

// FunctionsMap - map of functions
var FunctionsMap = make(map[string]interface{})

// RegexMap - map of regexes
var RegexMap = make(map[string]RegexAndDescription)

// RegexAndDescription struct which contains regexes and description of error
type RegexAndDescription struct {
	Regex       string `yaml:"regex"`
	Description string `yaml:"description"`
}

// Gostp - structure for gostp settings
type Gostp struct {
	WorkDir     string `yaml:"work_dir"`    // WorkDir - is a directory address where program has been launched
	SigningKey  string `yaml:"signing_key"` // SigningKey - key for signing JWT
	SQLtype     string `yaml:"sql_type"`    // SQLtype - type of gorm SQL
	SQLfilename string `yaml:"sql_filename"`
	SQLhost     string `yaml:"sql_host"`
	SQLport     string `yaml:"sql_port"`
	SQLdbname   string `yaml:"sql_dbname"`
	SQLuser     string `yaml:"sql_user"`
	SQLpassword string `yaml:"sql_password"`
	SQLssl      string `yaml:"sql_ssl"`
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

// Init - initialize of gospt
func Init(functionsMap map[string]interface{}, regexMap map[string]RegexAndDescription, useYaml bool) {
	// Functions Map initialization
	FunctionsMap = functionsMap
	// Regex Map initialization
	RegexMap = regexMap

	if useYaml {
		file, _ := ioutil.ReadFile(filepath.Join(Settings.WorkDir, Settings.SQLfilename))
		err := yaml.Unmarshal(file, &Settings)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	if Settings.SQLtype == "sqlite3" {
		Db, Err = gorm.Open(Settings.SQLtype, filepath.Join(Settings.WorkDir, "app.db"))
	} else {
		Db, Err = gorm.Open("postgres", "host="+Settings.SQLhost+" port="+Settings.SQLport+" user="+Settings.SQLuser+" dbname="+Settings.SQLdbname+" password="+Settings.SQLpassword+" sslmode="+Settings.SQLssl+"")
	}
}

package gostp

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"

	//DB-driver included in handlres
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Models - all app models
var Models []interface{}

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
	Regex       string // Regex - regular expression rule to validate string
	Description string // Description - description why string is not valid
}

// Settings - structure for gostp settings
var Settings = struct {
	Port                          string `yaml:"port"`
	WorkDir                       string `yaml:"work_dir"`       // WorkDir - is a directory address where program has been launched (default - directory where program stored)
	SigningKey                    string `yaml:"signing_key"`    // SigningKey - key for signing JWT (default - "")
	SSRMillisecondWait            int64  `yaml:"ssr_wait"`       // Time to wait after page loaded in ms (default - 1000 ms)
	SSRexpiration                 int64  `yaml:"ssr_expiration"` // Time after page will be deleted from cache in s (default - 86400 s)
	SSRhost                       string `yaml:"ssr_host"`       // Host for headless chrome rendering (default - "http://localhost:7777")
	SSRdevtools                   string `yaml:"ssr_devtools"`   // Headless Chrome Devtools address (default - "http://localhost:9222")
	SQLtype                       string `yaml:"sql_type"`       // SQLtype - type of gorm SQL (default - "sqlite3")
	SQLfilename                   string `yaml:"sql_filename"`   // SQLfilename - filename of sqlite db (default "app.db")
	SQLhost                       string `yaml:"sql_host"`       // SQLhost - host of remote or local db (defatult - "127.0.0.1")
	SQLport                       string `yaml:"sql_port"`       // SQLport - port of remote or local db (default - "5432")
	SQLdbname                     string `yaml:"sql_dbname"`     // SQLdbname - database name of remote or local db (default - "app")
	SQLuser                       string `yaml:"sql_user"`       // SQLuser - database username (default - "admin")
	SQLpassword                   string `yaml:"sql_password"`   // SQLpassword - database password (default - "admin")
	SQLsslmode                    string `yaml:"sql_sslmode"`    // SQLsslmode - database sslmode (default - "disabled")
	ServerName                    string `yaml:"server_name"`
	ContentType                   string `yaml:"content_type"`
	AccessControlAllowOrigin      string `yaml:"access_control_allow_origin"`
	AccessControlAllowMethods     string `yaml:"access_control_allow_methods"`
	AccessControlAllowHeaders     string `yaml:"access_control_allow_headers"`
	AccessControlAllowCredentials string `yaml:"access_control_allow_credentials"`
}{
	":7777",
	CurrentFolder(),
	"",
	1000,
	86400,
	"http://localhost:7777",
	"http://localhost:9222",
	"sqlite3",
	"app.db",
	"127.0.0.1",
	"5432",
	"app",
	"admin",
	"admin",
	"disable",
	"Gostp",
	"application/json",
	"*",
	"GET,PUT,POST,DELETE,OPTIONS",
	"Accept, Accept-Language, Content-Language, Content-Type, x-xsrf-token, authorization",
	"true"}

// Init - initialize of gostp
func Init(AppRoutes func(r *chi.Mux), functionsMap map[string]interface{}, regexMap map[string]RegexAndDescription, models ...interface{}) {
	r := chi.NewRouter()
	// Functions Map initialization
	FunctionsMap = functionsMap
	// Regex Map initialization
	RegexMap = regexMap
	// Models initialization
	Models = models

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
	// Check arguments from command line
	СheckArguments(&Settings.Port)
	// Start static file server
	filesDir := filepath.Join(Settings.WorkDir, "dist")
	FileServer(r, "/", http.Dir(filesDir))
	// Start all app routes
	AppRoutes(r)
	// Start http listener
	http.ListenAndServe(Settings.Port, r)
}

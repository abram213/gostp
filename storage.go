package gostp

import (
	"os"
	"path/filepath"
)

// RegexAndDescription struct which contains regexes and description of error
type RegexAndDescription struct {
	Regex       string
	Description string
}

// WorkDir - is a directory address where program has been launched
var WorkDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

// FunctionsMap - map of functions
var FunctionsMap map[string]interface{}

// RegexMap - map of regexes
var RegexMap map[string]RegexAndDescription

// InitRegex - initialize all regexes which needed to precheck values
func InitRegex(functionsMap map[string]interface{}, regexMap map[string]RegexAndDescription) {
	// Functions Map initialization
	FunctionsMap = functionsMap
	// Regex Map initialization
	RegexMap = regexMap
}

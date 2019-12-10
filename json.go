package gostp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/fatih/structtag"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func checkSecurity(tagsString string, handlerType, direction string, godMode bool) bool {
	tags := strings.Split(tagsString, ",")
	if !godMode {
		for _, tag := range tags {
			if tag == "protected" && direction == "in" {
				return true
			}

			if tag == "create_only" && handlerType != "create" && direction == "in" {
				return true
			}

			if tag == "update_only" && handlerType != "update" && direction == "in" {
				return true
			}

			if tag == "hidden_out" && direction == "out" {
				return true
			}

		}
	}
	return false
}

func deepInspection(model interface{}, parendJSON string, deletions []string, regexTagsMap, functionsTagsMap map[string]string, handlerType, direction string, godMode bool) {
	values := reflect.ValueOf(model)
	fields := reflect.TypeOf(model)
	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		tags, _ := structtag.Parse(string(field.Tag))
		jsonTag, _ := tags.Get("json")
		securityTag, _ := tags.Get("security")
		regexTag, _ := tags.Get("regex")
		functionTag, _ := tags.Get("function")
		valueValue := values.Field(i).Kind()
		var childJSON string
		if jsonTag != nil {
			childJSON = "." + jsonTag.Name
		} else {
			childJSON = ""
		}

		if valueValue == reflect.Struct && securityTag == nil {
			deepInspection(values.Field(i).Interface(), parendJSON+childJSON, deletions, regexTagsMap, functionsTagsMap, handlerType, direction, godMode)

		} else if valueValue == reflect.Struct && securityTag != nil {
			deleted := checkSecurity(securityTag.Name, handlerType, direction, godMode)
			if !deleted {
				deepInspection(values.Field(i).Interface(), parendJSON+childJSON, deletions, regexTagsMap, functionsTagsMap, handlerType, direction, godMode)
			} else {
				path := parendJSON + childJSON
				deletions = append(deletions, path[1:len(path)])
			}
		} else {
			if jsonTag != nil && jsonTag.Name != "-" {
				if securityTag != nil {
					deleted := checkSecurity(securityTag.Name, handlerType, direction, godMode)
					if deleted {
						path := parendJSON + childJSON
						deletions = append(deletions, path[1:len(path)])
					}
				}

				if regexTag != nil {
					path := parendJSON + childJSON
					regexTagsMap[path[1:len(path)]] = regexTag.Name
				}

				if functionTag != nil {
					//fmt.Println("Found function tag: ", functionTag.Name)
					path := parendJSON + childJSON
					functionsTagsMap[path[1:len(path)]] = functionTag.Name
				}
			}
		}

	}
}

// FillModelSafely - eee2
func FillModelSafely(r *http.Request, model interface{}, modelToFill interface{}, handlerType string, godMode bool) string {
	errors := ""
	var deletions []string                 // deletions pathes
	var regexTagsMap map[string]string     // map of regex tags
	var functionsTagsMap map[string]string // map of functions
	regexTagsMap = make(map[string]string)
	functionsTagsMap = make(map[string]string)
	direction := "in"
	data, ok := r.Context().Value("body").([]byte)
	if !ok {
		fmt.Println(ok)
	}
	deepInspection(model, "", deletions, regexTagsMap, functionsTagsMap, handlerType, direction, godMode)
	stringData := string(data)
	// Checking for regex error
	for k, v := range regexTagsMap {
		value := gjson.Get(stringData, k)
		var re = regexp.MustCompile(regexMap[v].regex)
		if !re.MatchString(value.Str) {
			errors = errors + "; " + regexMap[v].description
		}
	}
	if errors != "" {
		return errors[2:len(errors)]
	}
	// Check functions
	for k, v := range functionsTagsMap {
		value := gjson.Get(stringData, k)
		valueBefore := value.Str
		functionsMap[v].(func(*string))(&valueBefore)
		stringData, _ = sjson.Set(stringData, k, valueBefore)
	}
	// Checking for deletions
	for _, deletion := range deletions {
		stringData, _ = sjson.Delete(stringData, deletion)
	}
	fmt.Println(stringData)
	_ = json.Unmarshal([]byte(stringData), &modelToFill)

	return errors
}

// EncodeModelSafely - encodes given model by security rules
func EncodeModelSafely(model interface{}, modelToFill interface{}) []byte {
	byteJSON, _ := json.Marshal(modelToFill)
	encodedJSON := string(byteJSON)
	var deletions []string                 // deletions pathes
	var regexTagsMap map[string]string     // map of regex tags
	var functionsTagsMap map[string]string // map of functions
	regexTagsMap = make(map[string]string)
	functionsTagsMap = make(map[string]string)
	deepInspection(model, "", deletions, regexTagsMap, functionsTagsMap, "", "out", false)
	for _, deletion := range deletions {
		encodedJSON, _ = sjson.Delete(encodedJSON, deletion)
	}
	return []byte(encodedJSON)
}

package gostp

import (
	"fmt"
	"net/http"
	"reflect"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
)

// CheckAccess gets user by him token and checks accesss by struct fieldnames
func CheckAccess(r *http.Request, accesses []string, accessStruct interface{}) bool {
	token, _ := jwtmiddleware.FromAuthHeader(r)
	userID, _ := GetUserIDClaim(token)
	typeOfAccessStruct := reflect.TypeOf(accessStruct)
	newInterfaceOfAccessStruct := reflect.New(typeOfAccessStruct).Interface()
	if Db.Where("user_id = ?", userID).First(newInterfaceOfAccessStruct).RecordNotFound() {
		fmt.Println("user's access struct not found", userID)
		return false
	}
	accessStructValue := reflect.ValueOf(newInterfaceOfAccessStruct).Elem()
	// If user has Admin access - it function returns true all the time
	if accessStructValue.FieldByName("Admin").Bool() {
		return true
	}
	for _, access := range accesses {
		fmt.Println(access+" value : ", accessStructValue.FieldByName(access).Bool())
		if !accessStructValue.FieldByName(access).Bool() {
			return false
		}
	}
	return true
}

// CheckBelonging checks if user's some struct belogns to another through several structs.
func CheckBelonging(r *http.Request, path []string, models ...interface{}) bool {
	token, _ := jwtmiddleware.FromAuthHeader(r)
	userID, _ := GetUserIDClaim(token)
	previousID := uint64(userID)
	for index, model := range models {
		typeOfModel := reflect.TypeOf(model)
		newInterfaceOfModel := reflect.New(typeOfModel).Interface()
		if Db.Where(path[index]+" = ?", previousID).First(newInterfaceOfModel).RecordNotFound() {
			fmt.Println("user's access struct not found", userID)
			return false
		}
	}
	return true
}

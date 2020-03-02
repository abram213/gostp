package gostp

import (
	"net/http"
	"reflect"
	"strconv"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
)

// CheckAccess gets user by him token and checks accesses by struct fieldnames
func CheckAccess(r *http.Request, accesses []string, accessStruct interface{}, isAdmin *bool) bool {
	token, errTokenGet := jwtmiddleware.FromAuthHeader(r)
	if errTokenGet != nil {
		return false
	}
	userID, errUserGet := GetUserIDClaim(token)
	if errUserGet != nil {
		return false
	}
	typeOfAccessStruct := reflect.TypeOf(accessStruct)
	newInterfaceOfAccessStruct := reflect.New(typeOfAccessStruct).Interface()
	if Db.Where("user_id = ?", userID).First(newInterfaceOfAccessStruct).RecordNotFound() {
		return false
	}
	accessStructValue := reflect.ValueOf(newInterfaceOfAccessStruct).Elem()
	// If user has Admin access - this function returns true all the time
	if accessStructValue.FieldByName("Admin").Bool() {
		if isAdmin != nil {
			*isAdmin = true
		}
		return true
	}
	for _, access := range accesses {
		if !accessStructValue.FieldByName(access).Bool() {
			return false
		}
	}
	return true
}

// CheckBelonging checks if user's some struct belogns to another through several structs.
func CheckBelonging(r *http.Request, target string, isAdmin bool, path []string, models ...interface{}) bool {
	if isAdmin {
		return true
	}
	targetID, err := strconv.ParseUint(target, 10, 64)
	if err == nil {
		for index, model := range models {
			typeOfModel := reflect.TypeOf(model)
			newInterfaceOfModel := reflect.New(typeOfModel).Interface()
			if Db.Where("id = ?", targetID).First(newInterfaceOfModel).RecordNotFound() {
				return false
			}
			targetID = reflect.ValueOf(newInterfaceOfModel).Elem().FieldByName(path[index]).Uint()
		}
		token, errTokenGet := jwtmiddleware.FromAuthHeader(r)
		if errTokenGet != nil {
			return false
		}
		userID, errUserGet := GetUserIDClaim(token)
		if errUserGet != nil {
			return false
		}
		currentUserID := uint64(userID)
		if currentUserID == targetID {
			return true
		}
	}
	return false
}

// CheckCurrentUser - check current user by id from url
func CheckCurrentUser(r *http.Request, URLUserID string) bool {
	token, _ := jwtmiddleware.FromAuthHeader(r)
	userID, _ := GetUserIDClaim(token)
	currentUserID := uint64(userID)
	targetID, _ := strconv.ParseUint(URLUserID, 10, 64)
	return currentUserID == targetID
}

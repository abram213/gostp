package gostp

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns hashed and salted password
func HashPassword(password *string, functionError *string) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.MinCost)
	if err != nil {
		*functionError = err.Error()
	} else {
		*password = string(bytes)
	}
}

// CheckPasswordHash returns if hash is valid
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

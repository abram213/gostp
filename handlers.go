package gostp

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/buger/jsonparser"
)

// Login handles login attempts
func Login(w http.ResponseWriter, r *http.Request) *AppError {
	CommonHeader(w)
	var user User

	data, ok := r.Context().Value("body").([]byte)
	if !ok {
		return &AppError{Err, Err.Error(), 405}
	}

	email, _ := jsonparser.GetString(data, "email")
	password, _ := jsonparser.GetString(data, "password")

	if Db.Where("email = ?", email).First(&user).RecordNotFound() {
		return &AppError{Err, errors.New(`{"error": "Wrong email or password"}`).Error(), 405}
	}
	if CheckPasswordHash(password, user.Password) {
		var userTokens UserTokens
		RefreshUserTokens(user, &userTokens)
		json.NewEncoder(w).Encode(userTokens)
		return nil
	}
	return &AppError{Err, errors.New(`{"error": "Wrong email or password"}`).Error(), 405}
}

// RefreshTokens handles refresh token attempt
func RefreshTokens(w http.ResponseWriter, r *http.Request) *AppError {
	CommonHeader(w)
	data, ok := r.Context().Value("body").([]byte)
	if !ok {
		return &AppError{Err, Err.Error(), 401}
	}

	refreshToken, err := jsonparser.GetString(data, "refresh_token")
	token, err := JWTParse(refreshToken)
	if err != nil {
		return &AppError{err, err.Error(), 401}
	}

	userID, err := token.GetUserID()
	if err != nil {
		return &AppError{err, err.Error(), 401}
	}
	var user User
	if Db.Preload("Token").Where("id = ?", userID).First(&user).RecordNotFound() {
		return &AppError{err, errors.New(`{"error": "No such user"}`).Error(), 401}
	}

	if token.IsExpired() {
		return &AppError{err, err.Error(), 401}
	}

	var userTokens UserTokens
	RefreshUserTokens(user, &userTokens)
	json.NewEncoder(w).Encode(userTokens)
	return nil
}

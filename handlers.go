package gostp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AlecAivazis/survey"
	"github.com/tidwall/gjson"
)

// Login handles login attempts
func Login(w http.ResponseWriter, r *http.Request) *AppError {
	var user User

	data, ok := r.Context().Value("body").([]byte)
	if !ok {
		return &AppError{Err, Err.Error(), 405}
	}

	email := gjson.Get(string(data), "email").Str
	password := gjson.Get(string(data), "password").Str

	if Db.Where("email = ?", email).First(&user).RecordNotFound() {
		return &AppError{Err, errors.New(`{"error": "Wrong email or password"}`).Error(), 401}
	}
	if CheckPasswordHash(password, user.Password) {
		var userTokens UserTokens
		RefreshUserTokens(user, &userTokens)
		json.NewEncoder(w).Encode(userTokens)
		return nil
	}
	return &AppError{Err, errors.New(`{"error": "Wrong email or password"}`).Error(), 401}
}

// RefreshTokens handles refresh token attempt
func RefreshTokens(w http.ResponseWriter, r *http.Request) *AppError {
	data, ok := r.Context().Value("body").([]byte)
	if !ok {
		return &AppError{Err, Err.Error(), 401}
	}

	refreshToken := gjson.Get(string(data), "refresh_token").Str
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

// GenerateUser - generates super user
func GenerateUser() {
	var qs = []*survey.Question{
		{
			Name:     "email",
			Prompt:   &survey.Input{Message: "Enter email:"},
			Validate: survey.Required,
		},
		{
			Name:   "password",
			Prompt: &survey.Password{Message: "Enter password:"},
		},
		{
			Name: "role",
			Prompt: &survey.Select{
				Message: "Choose role:",
				Options: []string{"admin", "user"},
				Default: "admin",
			},
		},
	}
	// the answers will be written to this struct
	answers := struct {
		Email    string
		Password string
		Role     string
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var user User
	user.Email = answers.Email
	hashingError := ""
	HashPassword(&answers.Password, &hashingError)
	if hashingError == "" {
		user.Password = answers.Password
		Db.Create(&user)
		fmt.Printf("User with email: %s created with role: %s.\n", answers.Email, answers.Role)
	} else {
		fmt.Println("Hashing error:", hashingError)
	}
}

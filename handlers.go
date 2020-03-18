package gostp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

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

	username := gjson.Get(string(data), "username").Str
	password := gjson.Get(string(data), "password").Str

	if Db.Where("username = ?", username).First(&user).RecordNotFound() {
		return &AppError{Err, errors.New(`{"error": "wrong username or password"}`).Error(), 401}
	}
	if CheckPasswordHash(password, user.Password) {
		json.NewEncoder(w).Encode(RefreshUserTokens(user.ID))
		return nil
	}
	return &AppError{Err, errors.New(`{"error": "wrong username or password"}`).Error(), 401}
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

	if token.IsExpired() {
		return &AppError{err, err.Error(), 401}
	}

	json.NewEncoder(w).Encode(RefreshUserTokens(uint(userID)))
	return nil
}

// GenerateUser - generates super user
func GenerateUser() {
	var qs = []*survey.Question{
		{
			Name:     "username",
			Prompt:   &survey.Input{Message: "Enter username:"},
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
		Username string
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
	user.Username = answers.Username
	hashingError := ""
	HashPassword(&answers.Password, &hashingError)
	if hashingError == "" {
		user.Password = answers.Password
		Db.Create(&user)

		var access Access
		access.UserID = user.ID
		switch answers.Role {
		case "admin":
			access.Admin = true
		case "user":
			access.User = true
		}
		Db.Create(&access)

		fmt.Printf("User with username: %s created with role: %s.\n", answers.Username, answers.Role)
	} else {
		fmt.Println("Hashing error:", hashingError)
	}
}

func sendSW(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile(filepath.Join(Settings.WorkDir, "dist/service-worker.js"))
	if err != nil {
		http.Error(w, "Couldn't read file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Write(data)
}

func sendManifest(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile(filepath.Join(Settings.WorkDir, "dist/manifest.json"))
	if err != nil {
		http.Error(w, "Couldn't read file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

package system

import (
	"log"
	"net/http"
)

// ErrorHandler - handles http error
type ErrorHandler func(w http.ResponseWriter, r *http.Request) *AppError

// AppError error struct
type AppError struct {
	Error   error
	Message string
	Code    int
}

func (ah ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := ah(w, r); err != nil {
		log.Printf("Error - %v | %v %v%v from %v | %v", err.Message, r.Method, r.Host, r.URL.String(), r.RemoteAddr, err.Code)
		http.Error(w, err.Message, err.Code)
	}
}

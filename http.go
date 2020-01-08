package gostp

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
)

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", Settings.ServerName)
		w.Header().Set("Access-Control-Allow-Origin", Settings.AccessControlAllowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", Settings.AccessControlAllowMethods)
		w.Header().Set("Access-Control-Allow-Headers", Settings.AccessControlAllowHeaders)
		w.Header().Set("Access-Control-Allow-Credentials", Settings.AccessControlAllowCredentials)
		if _, err := os.Stat(fmt.Sprintf("%s", root) + r.RequestURI); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(Settings.WorkDir, "dist/index.html"))
		} else {
			fs.ServeHTTP(w, r)
		}
	}))
}

// Header sets header to all of handlers
func Header(w http.ResponseWriter) {
	w.Header().Set("Server", Settings.ServerName)
	w.Header().Set("Content-Type", Settings.ContentType)
	w.Header().Set("Access-Control-Allow-Origin", Settings.AccessControlAllowOrigin)
	w.Header().Set("Access-Control-Allow-Methods", Settings.AccessControlAllowMethods)
	w.Header().Set("Access-Control-Allow-Headers", Settings.AccessControlAllowHeaders)
	w.Header().Set("Access-Control-Allow-Credentials", Settings.AccessControlAllowCredentials)
}

// SendOptions to client
func SendOptions(w http.ResponseWriter, r *http.Request) {
	Header(w)
	w.WriteHeader(200)
}

// RequestBodyToByte converts request body to byte
func RequestBodyToByte(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		var key interface{}
		key = "body"
		ctx := context.WithValue(r.Context(), key, body)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetBodyFromContext gets body from content
func GetBodyFromContext(r *http.Request) ([]byte, error) {
	data, ok := r.Context().Value("body").([]byte)
	if !ok {
		return []byte(""), errors.New("No body in context")
	}
	return data, nil
}

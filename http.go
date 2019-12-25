package gostp

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
		w.Header().Set("Server", "Gostp")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Language, Content-Type, x-xsrf-token, authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if _, err := os.Stat(fmt.Sprintf("%s", root) + r.RequestURI); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(Settings.WorkDir, "dist/index.html"))
		} else {
			fs.ServeHTTP(w, r)
		}
	}))
}

// OkHeader sets header to all of handlers
func OkHeader(w http.ResponseWriter) {
	w.Header().Set("Server", "Gostp")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Language, Content-Type, x-xsrf-token, authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

// SendOptions to client
func SendOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Gostp")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Language, Content-Type, x-xsrf-token, authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
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
		ctx := context.WithValue(r.Context(), "body", body)
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

// GetImageFromRequest saves image to images folder
func GetImageFromRequest(r *http.Request, formName string) (string, error) {
	file, header, err := r.FormFile(formName)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer file.Close()

	//Generate sha1
	h := sha1.New()
	h.Write([]byte(header.Filename + time.Now().String()))
	serverFileNameHash := hex.EncodeToString(h.Sum(nil))
	fileName := serverFileNameHash + filepath.Ext(header.Filename)

	f, err := os.OpenFile(filepath.Join(Settings.WorkDir, "dist/images/"+fileName), os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)

	return fileName, nil
}

package gostp

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

// DeleteFile deletes file from disk
func DeleteFile(filename string) error {
	err := os.Remove(filepath.Join(Settings.WorkDir, filename))
	return err
}

// FileNotExist checks if file exist on disk
func FileNotExist(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		return nil
	} else if os.IsNotExist(err) {
		return err
	} else {
		return err
	}
}

// CurrentFolder shows folder where binary file of program located
func CurrentFolder() string {
	workDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return workDir
}

// GetFileFromRequest saves image to images folder
func GetFileFromRequest(w http.ResponseWriter, r *http.Request, formName string, sizeBytesLimit int64, path string, allowedExtensions []string) (string, error) {
	if sizeBytesLimit > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, sizeBytesLimit)
	}
	file, header, err := r.FormFile(formName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Generate sha1
	h := sha1.New()
	h.Write([]byte(header.Filename + time.Now().String()))
	serverFileNameHash := hex.EncodeToString(h.Sum(nil))
	fileName := serverFileNameHash + filepath.Ext(header.Filename)
	// Get token
	userID := 0.0
	token, err := jwtmiddleware.FromAuthHeader(r)
	if err == nil {
		userID, _ = GetUserIDClaim(token)
	}
	userIDString := strconv.FormatFloat(userID, 'f', -1, 64)

	userFolder := "dist/" + path + "/" + userIDString + "/"
	os.MkdirAll(filepath.Join(Settings.WorkDir, userFolder), os.ModePerm)
	f, err := os.OpenFile(filepath.Join(Settings.WorkDir, userFolder+fileName), os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)

	if allowedExtensions != nil {
		var allowedExtensionsTypes []types.Type
		for _, allowedExtension := range allowedExtensions {
			allowedExtensionsTypes = append(allowedExtensionsTypes, filetype.GetType(allowedExtension))
		}
		validFile := false
		checkedType, errFileType := filetype.MatchFile(filepath.Join(Settings.WorkDir, userFolder+fileName))
		if errFileType != nil {
			return "", errFileType
		}
		for _, allowedExtensionType := range allowedExtensionsTypes {
			if checkedType == allowedExtensionType {
				validFile = true
			}
		}
		if !validFile {
			DeleteFile(userFolder + fileName)
			return "", errors.New("file type is not valid")
		}
	}

	return "/" + path + "/" + userIDString + "/" + fileName, nil
}

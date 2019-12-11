package gostp

import (
	"os"
	"path/filepath"
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

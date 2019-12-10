package system

import (
	"os"
	"path/filepath"
)

func DeleteFile(filename string) error {
	err := os.Remove(filepath.Join(WorkDir, filename))
	return err
}

func FileNotExist(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		return nil
	} else if os.IsNotExist(err) {
		return err
	} else {
		return err
	}
}

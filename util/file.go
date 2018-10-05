package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

//SaveFile saves the given bytes into a fresh new file specified by its folder
//and name.
//
//If the file already exists then it will be replaced.
func SaveFile(logger *log.Logger, folder FolderPath, name string, b []byte) (string, error) {
	l := filepath.Join(folder.Path(), name)
	logger.Printf(LOG_SAVING, l)
	os.Remove(l)
	if _, err := os.Stat(name); os.IsNotExist(err) {
		e := os.MkdirAll(folder.Path(), 0700)
		if e != nil {
			logger.Printf(ERROR, e.Error())
			return l, e
		}

		logger.Printf(LOG_CREATING_FILE, l)

		f, e := os.Create(l)
		if e != nil {
			logger.Printf(ERROR, e.Error())
			return l, fmt.Errorf(ERROR_CREATING_CONFIG_FILE, name, e.Error())
		}
		defer f.Close()
		_, e = f.Write(b)
		if e != nil {
			return l, e
		}
	}
	return l, nil
}

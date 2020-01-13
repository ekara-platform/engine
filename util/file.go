package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//SaveFile saves the given bytes into a fresh new file specified by its folder
//and name.
//
//If the file already exists then it will be replaced.
func SaveFile(folder FolderPath, name string, b []byte) (string, error) {
	l := filepath.Join(folder.Path(), name)
	os.Remove(l)
	if _, err := os.Stat(name); os.IsNotExist(err) {
		e := os.MkdirAll(folder.Path(), 0700)
		if e != nil {
			return l, e
		}
		f, e := os.Create(l)
		if e != nil {
			return l, fmt.Errorf("error creating the configuration file: %s, %s", name, e.Error())
		}
		defer f.Close()
		_, e = f.Write(b)
		if e != nil {
			return l, e
		}
	}
	return l, nil
}

//FileRead return the file content as []byte]
func FileRead(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

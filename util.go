package engine

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func GetCwdUrl() (*url.URL, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cwd, err = filepath.Abs(cwd)
	if err != nil {
		return nil, err
	}
	cwd = filepath.ToSlash(cwd)
	if strings.HasPrefix(cwd, "/") {
		cwd = "file://" + cwd + "/"
	} else {
		// On windows, absolute paths don't start with /
		cwd = "file:///" + cwd + "/"
	}
	return url.Parse(cwd)
}

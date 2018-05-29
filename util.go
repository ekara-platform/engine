package engine

import (
	"github.com/lagoon-platform/model"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func PathToUrl(path string) (*url.URL, error) {
	if _, err := os.Stat(path); err == nil {
		path = "file://" + filepath.ToSlash(path)
	}
	u, e := url.Parse(path)
	if e != nil {
		return nil, e
	}
	return model.NormalizeUrl(u), nil
}

func EnsurePathSuffix(u *url.URL, suffix string) *url.URL {
	res := model.NormalizeUrl(u)
	if strings.HasSuffix(res.Path, suffix) {
		return res
	} else {
		if strings.HasSuffix(res.Path, "/") {
			res.Path = res.Path + suffix
		} else {
			res.Path = res.Path + "/" + suffix
		}
	}
	return res
}

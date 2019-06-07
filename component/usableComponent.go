package component

import (
	"os"
	"path/filepath"
)

type (
	UsableComponent interface {
		Name() string
		Templated() bool
		Release()
		RootPath() string
		ContainsFile(name string) (bool, MatchingPath)
		ContainsDirectory(name string) (bool, MatchingPath)
	}

	usable struct {
		release   func()
		path      string
		component componentDef
		cm        *context
		templated bool
	}
)

func (u usable) Name() string {
	return u.component.component.Id
}

func (u usable) Release() {
	u.release()
}

func (u usable) RootPath() string {
	return u.path
}

func (u usable) Templated() bool {
	return u.templated
}

func (u usable) ContainsFile(path string) (bool, MatchingPath) {
	return u.contains(false, path)
}

func (u usable) ContainsDirectory(path string) (bool, MatchingPath) {
	return u.contains(true, path)
}

func (u usable) contains(isFolder bool, path string) (bool, MatchingPath) {
	res := mPath{
		comp: u,
	}
	filePath := filepath.Join(u.path, path)
	if info, err := os.Stat(filePath); err == nil && (isFolder == info.IsDir()) {
		res.relativePath = path
		return true, res
	}
	return false, res
}

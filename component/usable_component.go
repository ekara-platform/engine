package component

import (
	"github.com/ekara-platform/model"
	"os"
	"path/filepath"
)

type (
	//UsableComponent Represent a component which can be used physically
	UsableComponent interface {
		//Name returns the name of the component
		Name() string
		//Templated returns true is the component content has been templated
		Templated() bool
		//Release deletes the templated content.
		Release()
		//RootPath returns the absolute path of the, eventually templated, component
		RootPath() string
		//ContainsFile returns the matching path of the searched file
		ContainsFile(name string) (bool, MatchingPath)
		//ContainsDirectory returns the matching path of the searched directory
		ContainsDirectory(name string) (bool, MatchingPath)
		//EnvVars return the environment variables necessary for this component to work
		EnvVars() model.EnvVars
	}

	usable struct {
		release   func()
		path      string
		component model.Component
		envVars   model.EnvVars
		templated bool
	}
)

func (u usable) EnvVars() model.EnvVars {
	return u.envVars
}

func (u usable) Name() string {
	return u.component.Id
}

func (u usable) Release() {
	if u.release != nil {
		u.release()
	}
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

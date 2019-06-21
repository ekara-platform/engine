package scm

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ekara-platform/engine/component/scm/file"
	"github.com/ekara-platform/engine/component/scm/git"
	"github.com/ekara-platform/model"
)

//scmHandler is the common definition of all SCM handlers used to acces
// to component repositories
type scmHandler interface {
	//Matches return true if a repository has already be fetched into the path and if its
	// remote configuration is the same than the desired  one
	Matches(repository model.Repository, path string) bool
	//Fetch fetches the repository content into the given path.
	Fetch(repository model.Repository, path string, auth model.Parameters) error
	//Update updates the repository content into the given path.
	Update(path string, auth model.Parameters) error
	//Switch executes a checkout to the desired reference
	Switch(path string, ref string) error
}

//Handler allows to fetch a component.
//
//If the component repository has already been fetched and if it matches
// then it will be updated, if not it will be fetched.
//
type Handler func() (FetchedComponent, error)

//GetHandler returns an handler able to fetch a component
func GetHandler(l *log.Logger, dir string, c model.Component) (Handler, error) {

	if c.Repository.Url.UpperScheme() == model.SchemeFile {
		return fetchThroughSccm(file.FileScmHandler{Logger: l}, dir, c, l), nil
	}
	switch c.Repository.Scm {
	case model.GitScm:
		return fetchThroughSccm(git.GitScmHandler{Logger: l}, dir, c, l), nil
	}
	return fetchThroughSccm(file.FileScmHandler{Logger: l}, dir, c, l), fmt.Errorf("Unsupported source control management : %s", c.Repository.Scm)

}

func fetchThroughSccm(scm scmHandler, dir string, c model.Component, l *log.Logger) func() (FetchedComponent, error) {
	return func() (FetchedComponent, error) {
		fc := FetchedComponent{
			ID:         c.Id,
			Descriptor: c.Repository.DescriptorName,
		}
		cPath := filepath.Join(dir, c.Id)
		fc.LocalPath = cPath
		if _, err := os.Stat(cPath); err == nil {
			if scm.Matches(c.Repository, cPath) {
				err := scm.Update(cPath, c.Authentication)
				if err != nil {
					return fc, err
				}
			} else {
				l.Println("directory " + cPath + " already exists but doesn't match component source, deleting it")
				err := os.RemoveAll(cPath)
				if err != nil {
					return fc, err
				}
				err = scm.Fetch(c.Repository, cPath, c.Authentication)
				if err != nil {
					return fc, err
				}
			}
		} else {
			err := scm.Fetch(c.Repository, cPath, c.Authentication)
			if err != nil {
				return fc, err
			}
		}
		err := scm.Switch(cPath, c.Ref)
		if err != nil {
			return fc, err
		}

		u, err := model.CreateUrl(fc.LocalPath)
		if err != nil {
			return fc, err
		}
		fc.LocalUrl = u

		du, err := model.CreateUrl(filepath.Join(fc.LocalPath, fc.Descriptor))
		if err != nil {
			return fc, err
		}
		fc.DescriptorUrl = du
		return fc, nil
	}
}

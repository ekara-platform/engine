package engine

import (
	"github.com/lagoon-platform/model"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

type ScmHandler interface {
	Matches(source *url.URL, path string) bool
	Fetch(source *url.URL, dest string) error
	Update(dest string) error
	Switch(path string, tag string) error
}

type ComponentManager interface {
	RegisterComponent(c model.Component)
	ComponentPath(id string) string
	ComponentsPaths() map[string]string

	Fetch(location string, version string) (*url.URL, error)
	Ensure() error
}

type componentManager struct {
	logger     *log.Logger
	directory  string
	components map[string]model.Component
	paths      map[string]string
}

func createComponentManager(ctx *context) ComponentManager {
	return &componentManager{
		logger:     ctx.logger,
		directory:  filepath.Join(ctx.directory, "components"),
		components: map[string]model.Component{},
		paths:      map[string]string{}}
}

func (cm *componentManager) RegisterComponent(c model.Component) {
	cm.logger.Println("Registering component " + c.Repository.String() + "@" + c.Version.String())
	cm.components[c.Id] = c
}

func (cm *componentManager) ComponentPath(id string) string {
	return filepath.Join(cm.directory, id)
}

func (cm *componentManager) ComponentsPaths() map[string]string {
	panic("implement me")
}

func (cm *componentManager) Fetch(location string, tag string) (*url.URL, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	baseUrl, err := PathToUrl(cwd)
	if err != nil {
		return nil, err
	}

	cId, cUrl, err := model.ResolveRepositoryInfo(baseUrl, location)
	if err != nil {
		return nil, err
	}

	path, e := cm.fetchComponent(cId, cUrl, tag)
	if e != nil {
		return nil, e
	}
	return PathToUrl(path)
}

func (cm *componentManager) Ensure() error {
	for cName, c := range cm.components {
		cm.logger.Println("Ensuring that component " + cName + " is available")
		path, err := cm.fetchComponent(c.Id, c.Repository, c.Version.String())
		if err != nil {
			return err
		}
		cm.paths[c.Id] = path
	}
	return nil
}

func (cm *componentManager) fetchComponent(cId string, cUrl *url.URL, tag string) (path string, err error) {
	scm := GitScmHandler{logger: cm.logger} // TODO dynamically select proper handler
	cPath := filepath.Join(cm.directory, cId)
	if _, err := os.Stat(cPath); err == nil {
		if scm.Matches(cUrl, cPath) {
			err = scm.Update(cPath)
			if err != nil {
				return "", err
			}
		} else {
			cm.logger.Println("directory " + cPath + " already exists but doesn't match component source, deleting it")
			err = os.RemoveAll(cPath)
			if err != nil {
				return "", err
			}
			err = scm.Fetch(cUrl, cPath)
			if err != nil {
				return "", err
			}
		}
	} else {
		err = scm.Fetch(cUrl, cPath)
		if err != nil {
			return "", err
		}
	}
	if tag != "" {
		err = scm.Switch(cPath, tag)
		if err != nil {
			return "", err
		}
	}
	return cPath, nil
}

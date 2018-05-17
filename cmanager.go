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
	RegisterComponent(component model.Component)
	ComponentPath(id string) string
	ComponentsPaths() map[string]string

	Fetch(repository string, version string) (string, error)
	Ensure() error
}

type componentManager struct {
	logger     *log.Logger
	directory  string
	components map[string]model.Component
}

func createComponentManager(ctx *context) (cm ComponentManager, err error) {
	cm = &componentManager{
		logger:     ctx.logger,
		directory:  filepath.Join(ctx.baseDir, "components"),
		components: make(map[string]model.Component)}
	return
}

func (cm *componentManager) RegisterComponent(c model.Component) {
	cm.components[c.Id] = c
}

func (cm *componentManager) ComponentPath(id string) string {
	return filepath.Join(cm.directory, id)
}

func (cm *componentManager) ComponentsPaths() map[string]string {
	panic("implement me")
}

func (cm *componentManager) Fetch(location string, version string) (path string, err error) {
	cId, cUrl, err := model.ResolveRepositoryInfo(&url.URL{}, location)
	if err != nil {
		return
	}

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
	err = scm.Switch(cPath, "v"+version)
	if err != nil {
		return "", err
	}

	return cPath, nil
}

func (cm *componentManager) Ensure() error {
	panic("implement me")
}

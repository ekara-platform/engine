package engine

import (
	"errors"
	"github.com/lagoon-platform/model"
	"log"
	"net/url"
	"path/filepath"
)

type ComponentManager interface {
	Fetch(repository string, version string) (string, error)
	RegisterComponent(component model.Component) error
	ComponentPath(id string) (string, error)
	ComponentsPaths() map[string]string

	Ensure() error
}

type componentManager struct {
	logger     *log.Logger
	directory  string
	components map[string]model.Component
	paths      map[string]string
}

func createComponentManager(ctx *context) (cm ComponentManager, err error) {
	absBaseDir, err := filepath.Abs(ctx.baseDir)
	if err != nil {
		return
	}
	cm = &componentManager{
		logger:     ctx.logger,
		directory:  absBaseDir,
		components: make(map[string]model.Component)}
	return
}

func (cm *componentManager) Fetch(repository string, version string) (string, error) {
	repoUrl, e := model.ResolveRepositoryUrl(&url.URL{}, repository)
	if e != nil {
		return "", e
	}
	return repoUrl.String(), nil
}

func (cm *componentManager) RegisterComponent(c model.Component) error {
	cm.components[c.Id] = c
	return nil
}

func (cm *componentManager) ComponentPath(id string) (string, error) {
	if path, ok := cm.paths[id]; ok {
		return path, nil
	}
	return "", errors.New("component not available: " + id)
}

func (cm *componentManager) ComponentsPaths() map[string]string {
	panic("implement me")
}

func (cm *componentManager) Ensure() error {
	for _, c := range cm.components {
		switch c.Scm {
		case model.Git:
			//fetchGitComponent(c.)
		}
	}
	return nil
}

package engine

import (
	"log"
	"github.com/lagoon-platform/model"
	"errors"
)

type ComponentManager interface {
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

func createComponentManager(logger *log.Logger, baseDir string) ComponentManager {
	return &componentManager{logger: logger, components: make(map[string]model.Component)}
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

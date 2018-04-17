package engine

import (
	"log"
	"github.com/lagoon-platform/model"
)

type ComponentManager interface {
	Ensure()
}

type componentManager struct {
	logger     *log.Logger
	components []model.Component
}

func createComponentManager(logger *log.Logger, env *model.Environment) ComponentManager {
	cm := componentManager{logger: logger}

	for _, provider := range env.Providers {
		cm.components = append(cm.components, provider.Component)
	}

	for _, stack := range env.Stacks {
		cm.components = append(cm.components, stack.Component)
	}

	return cm
}

func (componentManager) Ensure() {
	panic("implement me")
}

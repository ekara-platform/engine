package engine

import (
	"log"

	"github.com/lagoon-platform/model"
)

type Lagoon interface {
	Environment() model.Environment
}

type context struct {
	// Global state
	logger *log.Logger

	// Environment info
	location    string
	environment *model.Environment

	// Subsystems
	componentManager ComponentManager
}

func (c context) Environment() model.Environment {
	return *c.environment
}

// Create creates an environment descriptor based on the provider location.
//
// The location can be an URL over http or https
func Create(logger *log.Logger, location string) (Lagoon, error) {
	ctx := context{logger: logger, location: location}

	env, err := model.Parse(logger, location)
	if err != nil {
		return nil, err
	}
	ctx.environment = &env
	ctx.componentManager = createComponentManager(logger, &env)

	return ctx, nil
}

func CreateLagoon(logger *log.Logger, l Lagoon) {

}

func UpdateLagoon(logger *log.Logger, l Lagoon) {

}

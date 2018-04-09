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
}

func (c context) Environment() model.Environment {
	return *c.environment
}

// Create creates an environment descriptor based on the provider location.
//
// The location can be an URL over http or https or even a file system location.
func Create(logger *log.Logger, location string) (Lagoon, error, model.ValidationErrors) {
	ctx := context{logger: logger, location: location}

	env, err, vErrs := model.Parse(logger, location)
	if err != nil || vErrs.HasErrors() {
		return nil, err, vErrs
	}
	ctx.environment = &env
	return ctx, nil, vErrs
}

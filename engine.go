package engine

import (
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/lagoon-platform/model"
)

type Lagoon interface {
	Environment() model.Environment
	ComponentManager() ComponentManager
}

type context struct {
	logger      *log.Logger
	workDir     string
	environment model.Environment

	// Subsystems
	componentManager ComponentManager
}

// Create creates an environment descriptor based on the provider location.
//
// The location can be an URL over http or https or even a file system location.
func Create(logger *log.Logger, baseDir string, location string, tag string) (engine Lagoon, err error) {
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return
	}

	ctx := context{
		logger:  logger,
		workDir: absBaseDir}

	// Create component manager
	ctx.componentManager, err = createComponentManager(&ctx)
	if err != nil {
		return
	}

	// Fetch the main component
	envPath, err := ctx.componentManager.Fetch(location, tag)
	if err != nil {
		return
	}

	// Parse the environment descriptor from the main component
	ctx.environment, err = model.Parse(logger, filepath.Join(envPath, DescriptorFileName))
	if err != nil {
		switch err.(type) {
		case model.ValidationErrors:
			err.(model.ValidationErrors).Log(ctx.logger)
			if err.(model.ValidationErrors).HasErrors() {
				return
			}
		default:
			return
		}
	}

	// Register all environment components
	for pName, pComp := range ctx.environment.Providers {
		ctx.logger.Println("Registering provider " + pName)
		ctx.componentManager.RegisterComponent(pComp.Component)
	}

	// Use context as Lagoon facade
	engine = &ctx

	return
}

// BuildDescriptorUrl builds the url of environment descriptor based on the
// url received has parameter
//
func BuildDescriptorUrl(url url.URL) url.URL {
	if strings.HasSuffix(url.Path, "/") {
		url.Path = url.Path + DescriptorFileName
	} else {
		url.Path = url.Path + "/" + DescriptorFileName
	}
	return url
}

func (c *context) Environment() model.Environment {
	return c.environment
}

func (c *context) ComponentManager() ComponentManager {
	return c.componentManager
}

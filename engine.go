package engine

import (
	"log"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/lagoon-platform/model"
)

type Lagoon interface {
	Environment() model.Environment
	ComponentManager() ComponentManager
}

type context struct {
	logger *log.Logger

	// Subsystems
	componentManager ComponentManager

	// Environment
	baseDir     string
	environment model.Environment
}

// Create creates an environment descriptor based on the provider location.
//
// The location can be an URL over http or https or even a file system location.
func Create(logger *log.Logger, baseDir string, repository string, version string) (lagoon Lagoon, err error) {
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return
	}

	ctx := context{
		logger:  logger,
		baseDir: absBaseDir}

	// Create component manager
	ctx.componentManager, err = createComponentManager(&ctx)
	if err != nil {
		return
	}

	// Fetch the main component
	envPath, err := ctx.componentManager.Fetch(repository, version)
	if err != nil {
		return
	}
	ctx.componentManager.Ensure()

	// Parse the environment descriptor from the main component
	ctx.environment, err = model.Parse(logger, path.Join(envPath, DescriptorFileName))
	if err != nil {
		return
	}

	// Use context as Lagoon facade
	lagoon = &ctx

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

package engine

import (
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/lagoon-platform/model"
	"gopkg.in/yaml.v2"
)

type ScmHandler interface {
	Matches(source *url.URL, path string) bool
	Fetch(source *url.URL, path string) error
	Update(path string) error
	Switch(path string, ref string) error
}

type ComponentManager interface {
	RegisterComponent(c model.Component)
	ComponentPath(id string) string
	ComponentsPaths() map[string]string
	SaveComponentsPaths(log *log.Logger, e model.Environment, dest FolderPath) error
	Fetch(location string, ref string) (*url.URL, error)
	Ensure() error
}

type componentManager struct {
	logger     *log.Logger
	directory  string
	components map[string]model.Component
	paths      map[string]string
}

// FileMap is used to Marshal the map of downloaded componebts
type FileMap struct {
	File map[string]string `yaml:"component_path"`
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
	return cm.paths
}

func (cm *componentManager) SaveComponentsPaths(log *log.Logger, e model.Environment, dest FolderPath) error {
	err := cm.Ensure()
	if err != nil {
		return err
	}
	fMap := FileMap{}
	fMap.File = make(map[string]string)
	// Adding the provider components
	for _, p := range e.Providers {
		repName := cm.ComponentPath(p.Component.Id)
		fMap.File[p.Name] = repName
	}

	repCoreName := cm.ComponentPath(e.LagoonPlatform.Component.Id)
	fMap.File["core"] = repCoreName

	orchestratorName := cm.ComponentPath(e.Orchestrator.Component.Id)
	fMap.File["orchestrator"] = orchestratorName

	b, err := yaml.Marshal(&fMap)
	if err != nil {
		return err
	}
	err = SaveFile(log, dest, ComponentPathsFileName, b)
	if err != nil {
		return err
	}
	return nil
}

func (cm *componentManager) Fetch(location string, ref string) (*url.URL, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	baseUrl, err := model.PathToUrl(cwd)
	if err != nil {
		return nil, err
	}

	cId, cUrl, err := model.ResolveRepositoryInfo(baseUrl, location)
	if err != nil {
		return nil, err
	}

	path, e := cm.fetchComponent(cId, cUrl, ref)
	if e != nil {
		return nil, e
	}
	return model.PathToUrl(path)
}

func (cm *componentManager) Ensure() error {
	for cName, c := range cm.components {
		cm.logger.Println("Ensuring that component " + cName + " is available")
		path, err := cm.fetchComponent(c.Id, c.Repository, c.Version.String())
		if err != nil {
			return err
		}
		cm.logger.Printf("Paths added: \"%s=%s\"", c.Id, path)
		c := cm.components[c.Id]
		cm.paths[c.Id] = path
	}
	return nil
}

func (cm *componentManager) fetchComponent(cId string, cUrl *url.URL, ref string) (path string, err error) {
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

	err = scm.Switch(cPath, ref)
	if err != nil {
		return "", err
	}
	return cPath, nil
}

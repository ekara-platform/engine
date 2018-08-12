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
	ComponentPath(cId string) string
	ComponentsPaths() map[string]string
	SaveComponentsPaths(log *log.Logger, e model.Environment, dest FolderPath) error
	Ensure() error
}

type componentManager struct {
	logger      *log.Logger
	directory   string
	components  map[string]model.Component
	paths       map[string]string
	environment *model.Environment
	data        map[string]interface{}
}

// FileMap is used to Marshal the map of downloaded componebts
type FileMap struct {
	File map[string]string `yaml:"component_path"`
}

func createComponentManager(ctx *context) ComponentManager {
	return &componentManager{
		logger:      ctx.logger,
		directory:   filepath.Join(ctx.directory, "components"),
		components:  map[string]model.Component{},
		paths:       map[string]string{},
		environment: ctx.environment,
		data:        ctx.data}
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
	fMap.File = cm.ComponentsPaths()
	b, err := yaml.Marshal(&fMap)
	if err != nil {
		return err
	}
	_, err = SaveFile(log, dest, ComponentPathsFileName, b)
	if err != nil {
		return err
	}
	return nil
}

func (cm *componentManager) Ensure() error {
	for cName, c := range cm.components {
		cm.logger.Println("Ensuring that component " + cName + " is available")
		cPath, err := cm.fetchComponent(c.Id, c.Repository, c.Version.String())
		if err != nil {
			return err
		}
		cEnv, err := cm.parseComponentDescriptor(cPath)
		if err != nil {
			return err
		}
		if cEnv != nil {
			cm.logger.Printf("Merging component " + cName + " descriptor")
			cm.environment.Merge(cEnv)
		}
		cm.logger.Printf("Paths added: \"%s=%s\"", c.Id, cPath)
		c := cm.components[c.Id]
		cm.paths[c.Id] = cPath
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

func (cm *componentManager) parseComponentDescriptor(cPath string) (*model.Environment, error) {
	cDescriptor := filepath.Join(cPath, DescriptorFileName)
	if _, err := os.Stat(cDescriptor); err == nil {
		locationUrl, err := url.Parse(cDescriptor)
		if err != nil {
			return nil, err
		}
		locationUrl, err = model.NormalizeUrl(locationUrl)
		if err != nil {
			return nil, err
		}
		return model.ParseWithData(cm.logger, locationUrl, cm.data)
	} else {
		return nil, nil
	}
}

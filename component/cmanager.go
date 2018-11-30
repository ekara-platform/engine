package component

import "C"
import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
	"gopkg.in/yaml.v2"
)

const maxFetchIterations = 9

type ScmHandler interface {
	Matches(source *url.URL, path string) bool
	Fetch(source *url.URL, path string, auth model.Parameters) error
	Update(path string, auth model.Parameters) error
	Switch(path string, ref string) error
}

type ComponentManager interface {
	RegisterComponent(c model.Component, imports ... string)
	MatchingDirectories(dirName string) []string
	ComponentPath(cId string) string
	ComponentsPaths() map[string]string
	SaveComponentsPaths(log *log.Logger, dest util.FolderPath) error
	Ensure() error
	Environment() model.Environment
}

type componentDef struct {
	component model.Component
	imports   []string
}

type context struct {
	// Common to all environments
	logger *log.Logger
	data   map[string]interface{}

	// Local to one environment (in the case multiple environments will be supported)
	directory   string
	components  map[string]componentDef
	paths       map[string]string
	environment *model.Environment
}

// FileMap is used to Marshal the map of downloaded components
type fileMap struct {
	File map[string]string `yaml:"component_path"`
}

func CreateComponentManager(logger *log.Logger, data map[string]interface{}, baseDir string) ComponentManager {
	return &context{
		logger:      logger,
		environment: nil,
		directory:   filepath.Join(baseDir, "components"),
		components:  map[string]componentDef{},
		paths:       map[string]string{},
		data:        data,
	}
}

func (cm *context) RegisterComponent(c model.Component, imports ... string) {
	if _, ok := cm.components[c.Id]; !ok {
		cm.logger.Println("registering component " + c.Repository.String() + "@" + c.Version.String())
		cm.components[c.Id] = componentDef{
			component: c,
			imports:   imports}
	}
}

func (cm *context) ComponentPath(id string) string {
	return filepath.Join(cm.directory, id)
}

func (cm *context) ComponentsPaths() map[string]string {
	return cm.paths
}

func (cm *context) MatchingDirectories(dirName string) []string {
	result := make([]string, 0, 10)
	for cPath := range cm.paths {
		subDir := filepath.Join(cm.directory, cPath, dirName)
		if _, err := os.Stat(subDir); err == nil {
			result = append(result, subDir)
		}
	}
	return result
}

func (cm *context) SaveComponentsPaths(log *log.Logger, dest util.FolderPath) error {
	err := cm.Ensure()
	if err != nil {
		return err
	}
	fMap := fileMap{}
	fMap.File = cm.ComponentsPaths()
	b, err := yaml.Marshal(&fMap)
	if err != nil {
		return err
	}
	_, err = util.SaveFile(log, dest, util.ComponentPathsFileName, b)
	if err != nil {
		return err
	}
	return nil
}

func (cm *context) Ensure() error {
	for i := 0; i < maxFetchIterations && cm.isFetchNeeded(); i++ {
		// Fetch all known components
		for cId, c := range cm.components {
			if cm.isComponentFetchNeeded(cId) {
				// Fetch component
				cPath, err := cm.fetchComponent(c.component)
				if err != nil {
					return err
				}
				// Parse default descriptor
				err = cm.parseComponentDescriptor(cId, cPath, util.DescriptorFileName)
				if err != nil {
					return err
				}
				// Parse external imports
				for _, imp := range c.imports {
					if filepath.Clean(imp) != util.DescriptorFileName {
						err = cm.parseComponentDescriptor(cId, cPath, imp)
						if err != nil {
							return err
						}
					}
				}
				cm.logger.Printf("component %s is available in %s", c.component.Id, cPath)
				c := cm.components[c.component.Id]
				cm.paths[c.component.Id] = cPath
			}
		}

		// Register additionally discovered components
		if cm.environment != nil {
			cm.RegisterComponent(cm.environment.Ekara.Distribution)

			or, err := cm.environment.Orchestrator.Component.Resolve()
			if err != nil {
				return err
			}
			cm.RegisterComponent(or)
			for _, pComp := range cm.environment.Providers {
				pr, err := pComp.Component.Resolve()
				if err != nil {
					return err
				}
				cm.RegisterComponent(pr)
			}
			for _, sComp := range cm.environment.Stacks {
				sr, err := sComp.Component.Resolve()
				if err != nil {
					return err
				}
				cm.RegisterComponent(sr)
			}
			for _, tComp := range cm.environment.Tasks {
				tr, err := tComp.Component.Resolve()
				if err != nil {
					return err
				}
				cm.RegisterComponent(tr)
			}
		}
	}
	if cm.isFetchNeeded() {
		return errors.New(fmt.Sprintf("not all components have been fetched after %d iterations, check for import loops in descriptors", maxFetchIterations))
	} else {
		return nil
	}
}

func (cm *context) Environment() model.Environment {
	return *cm.environment
}

func (cm *context) isFetchNeeded() bool {
	for id := range cm.components {
		if cm.isComponentFetchNeeded(id) {
			return true
		}
	}
	return false
}

func (cm *context) isComponentFetchNeeded(id string) bool {
	_, present := cm.paths[id]
	return !present
}

func (cm *context) fetchComponent(c model.Component) (path string, err error) {
	scm := GitScmHandler{logger: cm.logger} // TODO dynamically select proper handler
	cPath := filepath.Join(cm.directory, c.Id)
	if _, err := os.Stat(cPath); err == nil {
		if scm.Matches(c.Repository, cPath) {
			err = scm.Update(cPath, c.Authentication)
			if err != nil {
				return "", err
			}
		} else {
			cm.logger.Println("directory " + cPath + " already exists but doesn't match component source, deleting it")
			err = os.RemoveAll(cPath)
			if err != nil {
				return "", err
			}
			err = scm.Fetch(c.Repository, cPath, c.Authentication)
			if err != nil {
				return "", err
			}
		}
	} else {
		err = scm.Fetch(c.Repository, cPath, c.Authentication)
		if err != nil {
			return "", err
		}
	}
	err = scm.Switch(cPath, c.Version.String())
	if err != nil {
		return "", err
	}
	return cPath, nil
}

func (cm *context) parseComponentDescriptor(cName string, cPath string, descriptor string) error {
	cDescriptor := filepath.Join(cPath, descriptor)
	if _, err := os.Stat(cDescriptor); err == nil {
		// Calculating descriptor path
		if strings.HasPrefix(cDescriptor, "/") {
			cDescriptor = "file://" + filepath.ToSlash(cDescriptor)
		} else {
			cDescriptor = "file:///" + filepath.ToSlash(cDescriptor)
		}
		locationUrl, err := url.Parse(cDescriptor)
		if err != nil {
			return err
		}
		locationUrl, err = model.NormalizeUrl(locationUrl)
		if err != nil {
			return err
		}

		// Parsing the descriptor
		cEnv, err := model.CreateEnvironment(locationUrl, cm.data)
		if err != nil {
			return err
		}

		// Merge the resulting environment into the global one
		if cm.environment == nil {
			cm.environment = &cEnv
		} else {
			err = cm.environment.Merge(cEnv)
			if err != nil {
				return err
			}
		}

		// Recursively parse descriptor internal imports
		for _, imp := range cEnv.Imports {
			cm.parseComponentDescriptor(cName, cPath, imp)
		}
	}
	return nil
}

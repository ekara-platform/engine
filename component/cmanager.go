package component

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
	Fetch(source *url.URL, path string) error
	Update(path string) error
	Switch(path string, ref string) error
}

type ComponentManager interface {
	RegisterComponent(c model.Component, descriptor string)
	MatchingDirectories(dirName string) []string
	ComponentPath(cId string) string
	ComponentsPaths() map[string]string
	SaveComponentsPaths(log *log.Logger, dest util.FolderPath) error
	Ensure() error
	Environment() model.Environment
}

type componentDef struct {
	component  model.Component
	descriptor string
}

type context struct {
	logger      *log.Logger
	environment *model.Environment
	data        map[string]interface{}

	directory  string
	components map[string]componentDef
	paths      map[string]string
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

func (cm *context) RegisterComponent(c model.Component, descriptor string) {
	if _, ok := cm.components[c.Id]; !ok {
		cm.logger.Println("Registering component " + c.Repository.String() + "@" + c.Version.String())
		cm.components[c.Id] = componentDef{
			component:  c,
			descriptor: descriptor}
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
		for cName, c := range cm.components {
			cm.logger.Println("Ensuring that component " + cName + " is available")
			cPath, err := cm.fetchComponent(c.component.Id, c.component.Repository, c.component.Version.String())
			if err != nil {
				return err
			}
			err = cm.parseComponentDescriptor(cName, cPath, c.descriptor)
			if err != nil {
				return err
			}
			cm.logger.Printf("Component %s has been downloaded in %s", c.component.Id, cPath)
			c := cm.components[c.component.Id]
			cm.paths[c.component.Id] = cPath
		}

		// Register additionally discovered components
		if cm.environment != nil {
			coreComponent, e := cm.environment.Ekara.Component.Resolve()
			if e != nil {
				return e
			}
			cm.RegisterComponent(coreComponent, util.DescriptorFileName)
			orchestratorComponent, e := cm.environment.Orchestrator.Component.Resolve()
			if e != nil {
				return e
			}
			cm.RegisterComponent(orchestratorComponent, util.DescriptorFileName)
			for _, provider := range cm.environment.Providers {
				providerComponent, e := provider.Component.Resolve()
				if e != nil {
					return e
				}
				cm.RegisterComponent(providerComponent, util.DescriptorFileName)
			}
			for _, stack := range cm.environment.Stacks {
				stackComponent, e := stack.Component.Resolve()
				if e != nil {
					return e
				}
				cm.RegisterComponent(stackComponent, util.DescriptorFileName)
			}
			for _, task := range cm.environment.Tasks {
				taskComponent, e := task.Component.Resolve()
				if e != nil {
					return e
				}
				cm.RegisterComponent(taskComponent, util.DescriptorFileName)
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
		if _, ok := cm.paths[id]; !ok {
			return true
		}
	}
	return false
}

func (cm *context) fetchComponent(cId string, cUrl *url.URL, ref string) (path string, err error) {
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

func (cm *context) parseComponentDescriptor(cName string, cPath string, descriptors ... string) (error) {
	for _, descriptor := range descriptors {
		cDescriptor := filepath.Join(cPath, descriptor)
		if _, err := os.Stat(cDescriptor); err == nil {
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
			cm.logger.Printf("Parsing descriptor %s from component %s", descriptor, cName)
			cEnv, err := model.CreateEnvironment(cm.logger, locationUrl, cm.data)
			if err != nil {
				return err
			}
			if cm.environment == nil {
				cm.environment = &cEnv
			} else {
				return cm.environment.Merge(cEnv)
			}
		}
	}
	return nil
}

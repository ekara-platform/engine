package component

import (
	"fmt"
	"log"

	"os"
	"path/filepath"

	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
	"gopkg.in/yaml.v2"
)

const maxFetchIterations = 9

//ComponentManager represents the common definition of all Component Manager
type ComponentManager interface {
	RegisterComponent(c model.Component)
	MatchingDirectories(dirName string) []string
	ComponentPath(cID string) string
	ComponentsPaths() map[string]string
	SaveComponentsPaths(log *log.Logger, dest util.FolderPath) error
	Ensure() error
	Environment() model.Environment
}

type componentDef struct {
	component model.Component
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

func (cm *context) RegisterComponent(c model.Component) {
	if _, ok := cm.components[c.Id]; !ok {
		cm.logger.Println("registering component " + c.Repository.Url.String() + "@" + c.Repository.Ref)
		cm.components[c.Id] = componentDef{
			component: c,
		}
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
		for cID, c := range cm.components {
			if cm.isComponentFetchNeeded(cID) {
				err := fetchComponent(cm, c.component)
				if err != nil {
					return err
				}

				// Registering the distribution
				if cm.environment != nil && cm.environment.Ekara != nil && cm.environment.Ekara.Distribution.Repository.Url != nil {
					d := model.Component(cm.environment.Ekara.Distribution)
					if _, ok := cm.components[d.Id]; !ok {
						cm.logger.Printf("registering a distribution")
						cm.RegisterComponent(d)
						err := fetchComponent(cm, d)
						if err != nil {
							return err
						}
						continue
					} else {
						cm.logger.Printf("a distribution has already been registered")
					}
				}
			}
		}

		// Register additionally discovered and used components
		if cm.environment != nil && cm.environment.Ekara != nil && len(cm.environment.Ekara.UsedComponents) > 0 {
			cm.logger.Printf("registering used components")
			for _, c := range cm.environment.Ekara.UsedComponents {
				cr, err := c.Resolve()
				if err != nil {
					return err
				}
				cm.RegisterComponent(cr)
			}
		}
	}
	if cm.isFetchNeeded() {
		return fmt.Errorf("not all components have been fetched after %d iterations, check for import loops in descriptors", maxFetchIterations)
	} else {
		return nil
	}
}

func fetchComponent(cm *context, c model.Component) error {
	cm.logger.Printf("fetching component %s ", c.Id)
	fComp, err := fetchThroughSccm(cm, c)
	if err != nil {
		return err
	}
	// Parse default descriptor
	err = cm.parseComponentDescriptor(fComp)
	if err != nil {
		return err
	}
	cm.logger.Printf("component %s is available in %s", c.Id, fComp.localPath)
	// TODO Change cm.paths to a map[string]FetchedComponent
	cm.paths[c.Id] = fComp.localPath
	return nil
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

func (cm *context) parseComponentDescriptor(fComp FetchedComponent) error {
	if fComp.hasDescriptor() {
		// Parsing the descriptor
		cEnv, err := model.CreateEnvironment(fComp.descriptorUrl, cm.data)
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
	}
	return nil
}

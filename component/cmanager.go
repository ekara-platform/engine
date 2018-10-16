package component

import (
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
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
	SaveComponentsPaths(log *log.Logger, dest util.FolderPath) error
	Ensure() error
}

type context struct {
	logger      *log.Logger
	environment *model.Environment
	data        map[string]interface{}

	directory  string
	components map[string]model.Component
	paths      map[string]string
}

// FileMap is used to Marshal the map of downloaded components
type fileMap struct {
	File map[string]string `yaml:"component_path"`
}

func CreateComponentManager(logger *log.Logger, environment *model.Environment, data map[string]interface{}, baseDir string) ComponentManager {
	return &context{
		logger:      logger,
		environment: environment,
		directory:   filepath.Join(baseDir, "components"),
		components:  map[string]model.Component{},
		paths:       map[string]string{},
		data:        data,
	}
}

func (cm *context) RegisterComponent(c model.Component) {
	cm.logger.Println("Registering component " + c.Repository.String() + "@" + c.Version.String())
	cm.components[c.Id] = c
}

func (cm *context) ComponentPath(id string) string {
	return filepath.Join(cm.directory, id)
}

func (cm *context) ComponentsPaths() map[string]string {
	return cm.paths
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

func (cm *context) parseComponentDescriptor(cPath string) (*model.Environment, error) {
	cDescriptor := filepath.Join(cPath, util.DescriptorFileName)
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

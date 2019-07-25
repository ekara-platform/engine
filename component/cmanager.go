package component

import (
	"log"
	"os"

	"path/filepath"

	"github.com/ekara-platform/engine/component/scm"
	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
	"gopkg.in/yaml.v2"
)

var releaseNothing = func() {
	// Doing nothing and it's fine
}

type (
	// ComponentManager represents the common definition of all Component Manager
	ComponentManager struct {
		Logger      *log.Logger
		data        *model.TemplateContext
		environment *model.Environment
		Directory   string
		Paths       map[string]scm.FetchedComponent
	}

	localRef struct {
		component model.Component
	}

	// FileMap is used to Marshal the map of downloaded components
	fileMap struct {
		File map[string]string `yaml:"component_path"`
	}
)

//CreateComponentManager creates a new component manager
func CreateComponentManager(logger *log.Logger, data *model.TemplateContext, baseDir string) *ComponentManager {
	c := &ComponentManager{
		Logger:      logger,
		environment: nil,
		Directory:   filepath.Join(baseDir, "components"),
		Paths:       map[string]scm.FetchedComponent{},
		data:        data,
	}
	c.environment = model.InitEnvironment()
	return c
}

func (cm *ComponentManager) isComponentFetched(id string) (val scm.FetchedComponent, present bool) {
	val, present = cm.Paths[id]
	return
}

func (cm *ComponentManager) EnsureOneComponent(c model.Component) error {
	cm.Logger.Printf("ensuring component: %s", c.Id)
	path, fetched := cm.isComponentFetched(c.Id)
	if !fetched {
		fComp, err := fetch(cm, c)
		if err != nil {
			cm.Logger.Printf("error fetching the component %s", err.Error())
			return err
		}
		path = fComp
	}
	if path.HasDescriptor() {
		cm.Logger.Printf("creating partial environment based on component %s", c.Id)
		descriptorYaml, err := model.ParseYamlDescriptor(path.DescriptorUrl, cm.data)
		if err != nil {
			cm.Logger.Printf("error parsing the descriptor %s", err.Error())
			return err
		}

		cEnv, err := model.CreateEnvironment(path.DescriptorUrl.String(), descriptorYaml, c.Id)
		if err != nil {
			return err
		}

		// Merge or keep the resulting environment into the global one
		if cm.environment == nil {
			cm.environment = cEnv
		} else {
			// We don't want to merge the templates defined into the environment
			// But instead we want to keep them into the component
			cm.environment.Platform().KeepTemplates(c, cEnv.Templates)
			cEnv.Templates = model.Patterns{}
			err = cm.environment.Merge(cEnv)

			if err != nil {
				return err
			}
		}
	}
	cm.data.Model = model.CreateTEnvironmentForEnvironment(*cm.environment)

	return nil
}

func (cm *ComponentManager) Environment() *model.Environment {
	return cm.environment
}

func (cm *ComponentManager) ComponentsPaths() map[string]string {
	res := make(map[string]string)
	for k, v := range cm.Paths {
		res[k] = v.LocalPath
	}
	return res
}

func (cm *ComponentManager) SaveComponentsPaths(log *log.Logger, dest util.FolderPath) error {
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

func (cm *ComponentManager) ContainsFile(name string, in ...model.ComponentReferencer) MatchingPaths {
	return cm.contains(false, name, in...)
}

func (cm *ComponentManager) ContainsDirectory(name string, in ...model.ComponentReferencer) MatchingPaths {
	return cm.contains(true, name, in...)
}

func (cm *ComponentManager) contains(isFolder bool, name string, in ...model.ComponentReferencer) MatchingPaths {
	res := MatchingPaths{
		Paths: make([]MatchingPath, 0, 0),
	}
	if len(in) > 0 {
		for _, v := range in {
			uv, err := cm.Use(v)
			if err != nil {
				cm.Logger.Printf("An error occured using the component %s : %s", v.ComponentName(), err.Error())
			}
			if isFolder {
				if ok, match := uv.ContainsDirectory(name); ok {
					res.Paths = append(res.Paths, match)
				} else {
					uv.Release()
				}
			} else {
				if ok, match := uv.ContainsFile(name); ok {
					res.Paths = append(res.Paths, match)
				} else {
					uv.Release()
				}
			}
		}
	} else {
		for _, comp := range cm.environment.Platform().Components {
			lRef := localRef{
				component: comp,
			}
			uv, err := cm.Use(lRef)
			if err != nil {
				cm.Logger.Printf("An error occured using the component %s : %s", lRef.ComponentName(), err.Error())
			}
			if isFolder {
				if ok, match := uv.ContainsDirectory(name); ok {
					res.Paths = append(res.Paths, match)
				} else {
					uv.Release()
				}
			} else {
				if ok, match := uv.ContainsFile(name); ok {
					res.Paths = append(res.Paths, match)
				} else {
					uv.Release()
				}
			}
		}
	}
	return res
}

//Use returns a UsableComponent matching the given reference.
//If the component corresponding to the reference contains a template
//definition then the component will be duplicated and templated before
// being returned as a UsableComponent.
// Don't forget to Release the UsableComponent once is processing is over...
func (cm *ComponentManager) Use(cr model.ComponentReferencer) (UsableComponent, error) {
	c := cm.environment.Platform().Components[cr.ComponentName()]
	if ok, patterns := c.Templatable(); ok {
		path, err := runTemplate(*cm.data, cm.Paths[cr.ComponentName()].LocalPath, patterns, cr)
		if err != nil {
			return usable{}, err
		}
		// No error no path then it has not been templated
		if err == nil && path == "" {
			goto TemplateFalse
		}
		return usable{
			cm:        cm,
			path:      path,
			release:   cleanup(path),
			component: cm.environment.Platform().Components[cr.ComponentName()],
			templated: true,
		}, nil
	}
TemplateFalse:
	return usable{
		cm:        cm,
		release:   releaseNothing,
		path:      filepath.Join(cm.Directory, cr.ComponentName()),
		component: cm.environment.Platform().Components[cr.ComponentName()],
		templated: false,
	}, nil
}

func cleanup(path string) func() {
	return func() {
		os.RemoveAll(path)
	}
}

//Component returns the referenced component
func (r localRef) Component() (model.Component, error) {
	return r.component, nil
}

//ComponentName returns the referenced component name
func (r localRef) ComponentName() string {
	return r.component.Id
}

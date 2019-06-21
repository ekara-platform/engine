package component

import (
	"fmt"
	"log"
	"os"
	"sort"

	"path/filepath"

	"github.com/ekara-platform/engine/component/scm"
	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
	"gopkg.in/yaml.v2"
)

const maxFetchIterations = 9

var releaseNothing = func() {
	// Doing nothing and it's fine
}

//ComponentManager represents the common definition of all Component Manager
type (
	ComponentManager interface {
		//RegisterComponent register a new compoment.
		//
		//The registration key of a component is its id.
		//
		//If the component has already been registered it will remain unaffected
		// by a potential registration of a new version of the component.
		RegisterComponent(parent string, c model.Component)
		ComponentsPaths() map[string]string
		SaveComponentsPaths(log *log.Logger, dest util.FolderPath) error
		Ensure() error
		EnsureOneComponent(cID string, c model.Component) (bool, error)
		Use(cr model.ComponentReferencer) UsableComponent
		Environment() model.Environment
		ContainsFile(name string, in ...model.ComponentReferencer) MatchingPaths
		ContainsDirectory(name string, in ...model.ComponentReferencer) MatchingPaths
	}

	context struct {
		// Common to all environments
		logger      *log.Logger
		data        *model.TemplateContext
		environment *model.Environment
		// Local to one environment (in the case multiple environments will be supported)
		directory string
		paths     map[string]string
	}

	localRef struct {
		component model.Component
	}

	// FileMap is used to Marshal the map of downloaded components
	fileMap struct {
		File map[string]string `yaml:"component_path"`
	}
)

func CreateComponentManager(logger *log.Logger, data *model.TemplateContext, baseDir string) ComponentManager {
	c := &context{
		logger:      logger,
		environment: nil,
		directory:   filepath.Join(baseDir, "components"),
		paths:       map[string]string{},
		data:        data,
	}
	c.environment = model.InitEnvironment()
	return c
}

func (cm *context) isFetchNeeded() bool {
	for id := range cm.environment.Ekara.Components {
		if cm.isComponentFetchNeeded(id) {
			return true
		}
	}
	return false
}

func (cm *context) isComponentFetchNeeded(id string) bool {
	_, present := cm.paths[id]
	if id == model.MainComponentId || id == model.EkaraComponentId {
		return !present
	}
	return !present && cm.environment.Ekara.Used(id)
}

func (cm *context) RegisterComponent(parent string, c model.Component) {
	cm.logger.Printf("registering component %s@%s with parent %s", c.Repository.Url.String(), c.Repository.Ref, parent)
	cm.environment.Ekara.RegisterComponent(parent, c)
}

func (cm *context) Ensure() error {
	for i := 0; i < maxFetchIterations && cm.isFetchNeeded(); i++ {
		val, cpt := cm.environment.Ekara.ToFetch()
		for i := 0; i <= cpt-1; i++ {
			c := <-val
			b, err := cm.EnsureOneComponent(c.Id, c)
			if err != nil {
				return err
			}
			if b {
				goto EnsureAgain
			}
		}
	}
	if cm.isFetchNeeded() {
		return fmt.Errorf("not all components have been fetched after %d iterations, check for import loops in descriptors", maxFetchIterations)
	} else {
		return nil
	}
EnsureAgain:
	cm.Ensure()
	return nil
}

func (cm *context) EnsureOneComponent(cID string, c model.Component) (bool, error) {
	var toRegister []model.Component
	if cm.isComponentFetchNeeded(cID) {
		var err error
		toRegister, err = fetchComponent(cm, c)
		if err != nil {
			return len(toRegister) > 0, err
		}

		if cID == model.MainComponentId {
			// Registering the distribution
			if cm.environment != nil && cm.environment.Ekara != nil && cm.environment.Ekara.Distribution.Repository.Url != nil {
				d := model.Component(cm.environment.Ekara.Distribution)
				if _, ok := cm.environment.Ekara.Components[d.Id]; !ok {
					cm.logger.Printf("registering a distribution")
					toRegisterDist, err := fetchComponent(cm, d)
					cm.RegisterComponent(cID, d)
					if err != nil {
						return len(toRegister) > 0, err
					}
					for _, v := range toRegisterDist {
						cm.RegisterComponent(d.Id, v)
					}
				}
			}
		}
		for _, v := range toRegister {
			cm.RegisterComponent(cID, v)
		}
	}
	return len(toRegister) > 0, nil
}

func fetchComponent(cm *context, c model.Component) ([]model.Component, error) {
	toRegister := make([]model.Component, 0, 0)
	h, err := scm.GetHandler(cm.logger, cm.directory, c)
	if err != nil {
		return toRegister, err
	}
	cm.logger.Printf("fetching component %s ", c.Id)
	fComp, err := h()
	if err != nil {
		return toRegister, err
	}
	// Parse default descriptor
	toRegister, err = cm.parseComponentDescriptor(fComp)
	if err != nil {
		return toRegister, err
	}
	cm.logger.Printf("component %s is available in %s", c.Id, fComp.LocalPath)

	cm.paths[c.Id] = fComp.LocalPath
	cm.environment.Ekara.SortedFetchedComponents = append(cm.environment.Ekara.SortedFetchedComponents, c.Id)
	return toRegister, nil
}

func (cm *context) parseComponentDescriptor(fComp scm.FetchedComponent) ([]model.Component, error) {
	toRegister := make([]model.Component, 0, 0)
	if fComp.HasDescriptor() {
		// Parsing the descriptor
		cm.logger.Printf("creating partial environment based on component %s", fComp.ID)
		cEnv, err := model.CreateEnvironment(fComp.DescriptorUrl, cm.data)
		if err != nil {
			return toRegister, err
		}

		// If the parsed environment has components we prepare them in order
		// to be register them in alphabetical order..
		if len(cEnv.Ekara.Components) > 0 {
			var keys []string
			for k := range cEnv.Ekara.Components {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				if val, ok := cEnv.Ekara.Components[k]; ok {
					toRegister = append(toRegister, val)
				}
			}
		}

		// Merge or keep the resulting environment into the global one
		if cm.environment == nil {
			cm.environment = cEnv
		} else {
			// We don't want to merge the templates defined into the environment
			// But instead we want to keep them into the component
			if len(cEnv.Templates.Content) > 0 {
				cm.logger.Printf("env has template %s", fComp.ID)
				comp := cm.environment.Ekara.Components[fComp.ID]
				comp.Templates = cEnv.Templates
				cm.environment.Ekara.Components[fComp.ID] = comp

				comp = cm.environment.Ekara.Components[fComp.ID]
				comp.Templates = cEnv.Templates
				cm.environment.Ekara.Components[fComp.ID] = comp

				cEnv.Templates = model.Patterns{}
			}

			err = cm.environment.Merge(cEnv)
			if err != nil {
				return toRegister, err
			}
		}
	}
	return toRegister, nil
}

func (cm *context) Environment() model.Environment {
	return *cm.environment
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

func (cm *context) ContainsFile(name string, in ...model.ComponentReferencer) MatchingPaths {
	return cm.contains(false, name, in...)
}

func (cm *context) ContainsDirectory(name string, in ...model.ComponentReferencer) MatchingPaths {
	return cm.contains(true, name, in...)
}

func (cm *context) contains(isFolder bool, name string, in ...model.ComponentReferencer) MatchingPaths {
	res := MatchingPaths{
		Paths: make([]MatchingPath, 0, 0),
	}
	if len(in) > 0 {
		for _, v := range in {
			uv := cm.Use(v)
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
		for _, comp := range cm.environment.Ekara.Components {
			lRef := localRef{
				component: comp,
			}
			uv := cm.Use(lRef)
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

func (cm *context) Use(cr model.ComponentReferencer) UsableComponent {
	c := cm.environment.Ekara.Components[cr.ComponentName()]
	if ok, patterns := c.Templatable(); ok {
		path, err := runTemplate(*cm.data, cm.paths[cr.ComponentName()], patterns, cr)
		if err != nil {
			//TODO Return the error here !!!
		}
		// No error no path then it has not been templated
		if err == nil && path == "" {
			goto TemplateFalse
		}
		return usable{
			cm:        cm,
			path:      path,
			release:   cleanup(path),
			component: cm.environment.Ekara.Components[cr.ComponentName()],
			templated: true,
		}
	}
TemplateFalse:
	return usable{
		cm:        cm,
		release:   releaseNothing,
		path:      filepath.Join(cm.directory, cr.ComponentName()),
		component: cm.environment.Ekara.Components[cr.ComponentName()],
		templated: false,
	}
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

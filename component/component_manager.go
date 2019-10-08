package component

import (
	"log"
	"os"

	"path/filepath"

	"github.com/ekara-platform/engine/component/scm"
	"github.com/ekara-platform/model"
)

var releaseNothing = func() {
	// Doing nothing and it's fine
}

type (
	//Manager manages the fetch and the templating of components used into a descriptor
	Manager interface {
		Init(mainComponent model.Component) error
		Ensure() error
		Environment() *model.Environment
		TemplateContext() *model.TemplateContext
		ContainsFile(name string, in ...model.ComponentReferencer) MatchingPaths
		ContainsDirectory(name string, in ...model.ComponentReferencer) MatchingPaths
		Use(cr model.ComponentReferencer) (UsableComponent, error)
	}

	// ComponentManager downloads and keep track of ekara components on disk.
	manager struct {
		l                *log.Logger
		directory        string
		paths            map[string]scm.FetchedComponent
		referenceManager *ReferenceManager

		environment *model.Environment
		tplC        *model.TemplateContext
	}

	localRef struct {
		component model.Component
	}
)

//CreateComponentManager creates a new component manager
func CreateComponentManager(l *log.Logger, p model.Parameters, baseDir string) Manager {
	c := &manager{
		l:           l,
		environment: nil,
		directory:   filepath.Join(baseDir, "components"),
		paths:       map[string]scm.FetchedComponent{},
		tplC:        model.CreateTemplateContext(p),
	}
	c.environment = model.InitEnvironment()
	c.referenceManager = CreateReferenceManager(l, c)
	return c
}

func (cm *manager) isComponentFetched(id string) (val scm.FetchedComponent, present bool) {
	val, present = cm.paths[id]
	return
}

func (cm *manager) ensureOneComponent(c model.Component, data *model.TemplateContext) error {
	cm.l.Printf("ensuring component: %s", c.Id)
	path, fetched := cm.isComponentFetched(c.Id)
	if !fetched {
		fComp, err := fetch(cm.l, cm.directory, c)
		if err != nil {
			cm.l.Printf("error fetching the component: %s", err.Error())
			return err
		}
		cm.paths[c.Id] = fComp
		path = fComp
	}
	if path.HasDescriptor() {
		cm.l.Printf("creating partial environment based on component %s", c.Id)
		descriptorYaml, err := model.ParseYamlDescriptor(path.DescriptorUrl, data)
		if err != nil {
			cm.l.Printf("error parsing the descriptor: %s", err.Error())
			return err
		}

		cEnv, err := model.CreateEnvironment(path.DescriptorUrl.String(), descriptorYaml, c.Id)
		if err != nil {
			return err
		}

		// Customize or keep the resulting environment into the global one
		cm.l.Println("prepare partial environment customization")
		if cm.environment == nil {
			cm.environment = cEnv
			cm.l.Println("no customization required, it's the first built environment ")
		} else {
			// We don't want to customize the templates defined into the environment
			// But instead we want to keep them into the component
			cm.environment.Platform().KeepTemplates(c, cEnv.Templates)
			cEnv.Templates = model.Patterns{}
			cm.l.Println("partial environment should be used for customization")
			err = cm.environment.Customize(cEnv)

			if err != nil {
				cm.l.Printf("error customizing the environment %s", err.Error())
				return err
			}
		}
	}
	data.Model = model.CreateTEnvironmentForEnvironment(*cm.environment)

	return nil
}

func (cm *manager) Init(mainComponent model.Component) error {
	return cm.referenceManager.init(mainComponent)
}

func (cm *manager) Ensure() error {
	return cm.referenceManager.ensure()
}

func (cm *manager) Environment() *model.Environment {
	return cm.environment
}

func (cm *manager) TemplateContext() *model.TemplateContext {
	return cm.tplC
}

func (cm *manager) ContainsFile(name string, in ...model.ComponentReferencer) MatchingPaths {
	return cm.contains(false, name, in...)
}

func (cm *manager) ContainsDirectory(name string, in ...model.ComponentReferencer) MatchingPaths {
	return cm.contains(true, name, in...)
}

func (cm *manager) contains(isFolder bool, name string, in ...model.ComponentReferencer) MatchingPaths {
	res := MatchingPaths{
		Paths: make([]MatchingPath, 0, 0),
	}
	if len(in) > 0 {
		for _, v := range in {
			uv, err := cm.Use(v)
			if err != nil {
				cm.l.Printf("An error occurred using the component %s : %s", v.ComponentName(), err.Error())
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
				cm.l.Printf("An error occurred using the component %s : %s", lRef.ComponentName(), err.Error())
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
func (cm *manager) Use(cr model.ComponentReferencer) (UsableComponent, error) {
	c := cm.environment.Platform().Components[cr.ComponentName()]
	if ok, patterns := c.Templatable(); ok {
		path, err := runTemplate(cm.tplC, cm.paths[cr.ComponentName()].LocalPath, patterns, cr)
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
		path:      filepath.Join(cm.directory, cr.ComponentName()),
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

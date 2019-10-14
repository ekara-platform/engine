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
	/*
		Manager interface {
			TemplateContext() *model.TemplateContext
			ContainsFile(name string, in ...model.ComponentReferencer) MatchingPaths
			ContainsDirectory(name string, in ...model.ComponentReferencer) MatchingPaths
			Use(cr model.ComponentReferencer) (UsableComponent, error)
		}
	*/

	//Manager manages the fetch and the templating of components used into a descriptor
	Manager struct {
		l         *log.Logger
		platform  *model.Platform
		directory string
		paths     map[string]scm.FetchedComponent
		tplC      *model.TemplateContext
	}

	localRef struct {
		component model.Component
	}
)

//CreateComponentManager creates a new component manager
func CreateComponentManager(l *log.Logger, p model.Parameters, pl *model.Platform, baseDir string) *Manager {
	c := &Manager{
		l:         l,
		platform:  pl,
		directory: filepath.Join(baseDir, "components"),
		paths:     map[string]scm.FetchedComponent{},
		tplC:      model.CreateTemplateContext(p),
	}
	return c
}

func (cm *Manager) isComponentFetched(id string) (val scm.FetchedComponent, present bool) {
	val, present = cm.paths[id]
	return
}

func (cm *Manager) ensureOneComponent(c model.Component) (model.EkURL, bool, error) {
	cm.l.Printf("ensuring component: %s", c.Id)
	path, fetched := cm.isComponentFetched(c.Id)
	if !fetched {
		fComp, err := fetch(cm.l, cm.directory, c)
		if err != nil {
			cm.l.Printf("error fetching the component: %s", err.Error())
			return nil, false, err
		}
		cm.paths[c.Id] = fComp
		path = fComp
	}
	return path.DescriptorUrl, path.HasDescriptor(), nil
}

func (cm *Manager) TemplateContext() *model.TemplateContext {
	return cm.tplC
}

func (cm *Manager) ContainsFile(name string, in ...model.ComponentReferencer) MatchingPaths {
	return cm.contains(false, name, in...)
}

func (cm *Manager) ContainsDirectory(name string, in ...model.ComponentReferencer) MatchingPaths {
	return cm.contains(true, name, in...)
}

func (cm *Manager) contains(isFolder bool, name string, in ...model.ComponentReferencer) MatchingPaths {
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
		for _, comp := range cm.platform.Components {
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
func (cm *Manager) Use(cr model.ComponentReferencer) (UsableComponent, error) {
	c := cm.platform.Components[cr.ComponentName()]
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
			path:      path,
			release:   cleanup(path),
			component: cm.platform.Components[cr.ComponentName()],
			templated: true,
		}, nil
	}
TemplateFalse:
	return usable{
		release:   releaseNothing,
		path:      filepath.Join(cm.directory, cr.ComponentName()),
		component: cm.platform.Components[cr.ComponentName()],
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

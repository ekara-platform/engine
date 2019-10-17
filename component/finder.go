package component

import (
	"log"
	"path/filepath"

	"github.com/ekara-platform/model"
)

type (
	//Finder allows to locate and use component.
	//If required the located component will be templated.
	Finder interface {
		//ContainsFile returns paths pointing on templated components containing the given file
		//	Parameters
		//		name: the name of the file to search
		//		ctx: the context used to eventually template the matching components
		//		in: the component referencer where to look for the file, if not provided
		//          the Finder will look into all the components available into the platform.
		ContainsFile(name string, ctx model.TemplateContext, in ...model.ComponentReferencer) MatchingPaths

		//ContainsDirectory returns paths pointing on templated components containing the given directory
		//	Parameters
		//		name: the name of the directory to search
		//		ctx: the context used to eventually template the matching components
		//		in: the component referencer where to look for the directory, if not provided
		//          the Finder will look into all the components available into the platform.
		ContainsDirectory(name string, ctx model.TemplateContext, in ...model.ComponentReferencer) MatchingPaths

		//Use returns a UsableComponent matching the given reference.
		//If the component corresponding to the reference contains a template
		//definition then the component will be duplicated and templated before
		// being returned as a UsableComponent.
		// Don't forget to Release the UsableComponent once is processing is over...
		Use(cr model.ComponentReferencer, ctx model.TemplateContext) (UsableComponent, error)
	}

	finder struct {
		l    *log.Logger
		base string
		p    model.Platform
	}

	localRef struct {
		component model.Component
	}
)

func createFinder(l *log.Logger, baseDir string, p model.Platform) Finder {
	return finder{
		l:    l,
		base: baseDir,
		p:    p,
	}
}

func (f finder) ContainsFile(name string, ctx model.TemplateContext, in ...model.ComponentReferencer) MatchingPaths {
	return f.contains(false, name, ctx, in...)
}

func (f finder) ContainsDirectory(name string, ctx model.TemplateContext, in ...model.ComponentReferencer) MatchingPaths {
	return f.contains(true, name, ctx, in...)
}

func (f finder) contains(isFolder bool, name string, ctx model.TemplateContext, in ...model.ComponentReferencer) MatchingPaths {
	res := MatchingPaths{
		Paths: make([]MatchingPath, 0, 0),
	}
	if len(in) > 0 {
		for _, v := range in {
			if match, b := f.checkMatch(v, ctx, name, isFolder); b {
				res.Paths = append(res.Paths, match)
			}
		}
	} else {
		for _, comp := range f.p.Components {
			lRef := localRef{
				component: comp,
			}
			if match, b := f.checkMatch(lRef, ctx, name, isFolder); b {
				res.Paths = append(res.Paths, match)
			}
		}
	}
	return res
}

func (f finder) checkMatch(r model.ComponentReferencer, ctx model.TemplateContext, name string, isFolder bool) (MatchingPath, bool) {
	uv, err := f.Use(r, ctx)
	if err != nil {
		f.l.Printf("An error occurred using the component %s : %s", r.ComponentName(), err.Error())
		return mPath{}, false
	}
	if isFolder {
		if ok, match := uv.ContainsDirectory(name); ok {
			return match, true
		} else {
			uv.Release()
		}
	} else {
		if ok, match := uv.ContainsFile(name); ok {
			return match, true
		} else {
			uv.Release()
		}
	}
	return mPath{}, false
}

func (f finder) Use(cr model.ComponentReferencer, ctx model.TemplateContext) (UsableComponent, error) {

	lPath := filepath.Join(f.base, cr.ComponentName())
	c := f.p.Components[cr.ComponentName()]
	if ok, patterns := c.Templatable(); ok {
		path, err := runTemplate(ctx, lPath, patterns, cr)
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
			component: c,
			templated: true,
		}, nil
	}
TemplateFalse:
	return usable{
		release:   releaseNothing,
		path:      lPath,
		component: c,
		templated: false,
	}, nil
}

//Component returns the referenced component
func (r localRef) Component() (model.Component, error) {
	return r.component, nil
}

//ComponentName returns the referenced component name
func (r localRef) ComponentName() string {
	return r.component.Id
}

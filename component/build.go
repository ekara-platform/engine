package component

import (
	"log"
	"path/filepath"

	"github.com/ekara-platform/model"
)

//Build creates the base environment base on the given main component
func Build(mainComponent model.Component, ctx *model.TemplateContext, workDir string, l *log.Logger) (Finder, model.Environment, error) {

	rM := createReferenceManager(l)
	cM := createComponentManager(l, workDir)
	env := model.InitEnvironment()

	// Discover components starting from the main one
	err := rM.init(mainComponent, cM, ctx)
	if err != nil {
		return finder{}, *env, err
	}

	// Then ensure all effectively used components are fetched
	err = rM.ensure(env, cM, ctx)
	if err != nil {
		return finder{}, *env, err
	}

	finder := createFinder(l, filepath.Join(workDir, "components"), *env.Platform())
	return finder, *env, nil
}
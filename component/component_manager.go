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
	canager struct {
		l         *log.Logger
		directory string
		paths     map[string]scm.FetchedComponent
	}
)

//CreateComponentManager creates a new component manager
func createComponentManager(l *log.Logger, baseDir string) *canager {
	c := &Manager{
		l:         l,
		directory: filepath.Join(baseDir, "components"),
		paths:     map[string]scm.FetchedComponent{},
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

func cleanup(path string) func() {
	return func() {
		os.RemoveAll(path)
	}
}

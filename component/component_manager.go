package component

import (
	"log"
	"path/filepath"

	"github.com/ekara-platform/engine/component/scm"
	"github.com/ekara-platform/model"
)

type (
	//Manager manages the fetch and the templating of components used into a descriptor
	manager struct {
		l         *log.Logger
		directory string
		paths     map[string]scm.FetchedComponent
	}
)

//createComponentManager creates a new component manager
func createComponentManager(l *log.Logger, baseDir string) *manager {
	c := &manager{
		l:         l,
		directory: filepath.Join(baseDir, "components"),
		paths:     map[string]scm.FetchedComponent{},
	}
	return c
}

func (cm *manager) isComponentFetched(id string) (val scm.FetchedComponent, present bool) {
	val, present = cm.paths[id]
	return
}

func (cm *manager) ensureOneComponent(c model.Component) (model.EkURL, bool, error) {
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

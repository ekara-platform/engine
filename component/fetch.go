package component

import (
	"github.com/ekara-platform/engine/component/scm"
	"github.com/ekara-platform/model"
)

func fetch(cm *ComponentManager, c model.Component) (scm.FetchedComponent, error) {
	h, err := scm.GetHandler(cm.Logger, cm.Directory, c)
	if err != nil {
		return scm.FetchedComponent{}, err
	}
	cm.Logger.Printf("fetching component %s ", c.Id)
	fComp, err := h()
	if err != nil {
		return scm.FetchedComponent{}, err
	}
	cm.Logger.Printf("component %s is available in %s", c.Id, fComp.LocalPath)
	cm.Paths[c.Id] = fComp
	return fComp, nil
}

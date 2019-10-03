package component

import (
	"github.com/ekara-platform/engine/component/scm"
	"github.com/ekara-platform/model"
)

func fetch(cm *manager, c model.Component) (scm.FetchedComponent, error) {
	h, err := scm.GetHandler(cm.lC.Log(), cm.Directory, c)
	if err != nil {
		return scm.FetchedComponent{}, err
	}
	cm.lC.Log().Printf("fetching component %s ", c.Id)
	fComp, err := h()
	if err != nil {
		return scm.FetchedComponent{}, err
	}
	cm.lC.Log().Printf("component %s is available in %s", c.Id, fComp.LocalPath)
	cm.Paths[c.Id] = fComp
	return fComp, nil
}

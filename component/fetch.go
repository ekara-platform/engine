package component

import (
	"log"

	"github.com/ekara-platform/engine/component/scm"
	"github.com/ekara-platform/model"
)

func fetch(l *log.Logger, destination string, c model.Component) (scm.FetchedComponent, error) {
	h, err := scm.GetHandler(l, destination, c)
	if err != nil {
		return scm.FetchedComponent{}, err
	}
	l.Printf("fetching component %s ", c.Id)
	fComp, err := h()
	if err != nil {
		return scm.FetchedComponent{}, err
	}
	l.Printf("component %s is available in %s", c.Id, fComp.LocalPath)
	return fComp, nil
}

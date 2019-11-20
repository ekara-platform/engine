package action

import (
	"fmt"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
)

//Result represents the result of an action
type Result interface {
	//IsSuccess returns true id the action execution was successful
	IsSuccess() bool
	//AsJson returns the action returned content as JSON
	FromJson(s string) error
	//AsJson returns the action returned content as JSON
	AsJson() (string, error)
}

type (
	//Manager is the manager of all action available into the engine
	Manager interface {
		Run(id ActionID) (Result, error)
	}

	manager struct {
		actions map[ActionID]Action

		// Interfaces to other components
		lC util.LaunchContext
		cF component.Finder
		aM ansible.Manager

		// Model in use
		tplC *model.TemplateContext
		env  *model.Environment
	}
)

//CreateActionManager initializes the action manager and its content
func CreateActionManager(lC util.LaunchContext, tplC model.TemplateContext, env model.Environment, cF component.Finder, aM ansible.Manager) Manager {
	am := &manager{
		actions: make(map[ActionID]Action),
		lC:      lC,
		cF:      cF,
		aM:      aM,
		tplC:    &tplC,
		env:     &env,
	}

	for _, a := range allActions() {
		am.actions[a.id] = a
	}

	return am
}

func (am *manager) empty() bool {
	return len(am.actions) == 0
}

//get returns the action corresponding to the given id.
func (am *manager) get(id ActionID) (Action, error) {
	if val, ok := am.actions[id]; ok {
		return val, nil
	}
	return Action{}, fmt.Errorf("unsupported action")
}

//Run launches the action corresponding to the given id.
func (am *manager) Run(id ActionID) (Result, error) {
	rC := createRuntimeContext(am.lC, am.cF, am.aM, am.env, am.tplC)
	a, e := am.get(id)
	if e != nil {
		return nil, e
	}

	report, res, e := a.run(am, rC)
	if e != nil {
		return nil, e
	}

	loc, e := writeReport(*report, rC.lC.Ef().Output)
	if e != nil {
		return nil, e
	}
	rC.lC.Log().Printf(LogReportWritten, loc)

	return res, nil
}

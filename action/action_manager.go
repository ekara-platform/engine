package action

import (
	"fmt"
	"github.com/GroupePSA/componentizer"
	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/model"
	"github.com/ekara-platform/engine/util"
)

//Result represents the result of an action
type Result interface {
	//IsSuccess returns true id the action execution was successful
	IsSuccess() bool
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
		cM componentizer.ComponentManager
		aM ansible.Manager

		// Model in use
		tplC componentizer.TemplateContext
		env  model.Environment
	}
)

//CreateActionManager initializes the action manager and its content
func CreateActionManager(lC util.LaunchContext, cM componentizer.ComponentManager, aM ansible.Manager) Manager {
	am := &manager{
		actions: make(map[ActionID]Action),
		lC:      lC,
		cM:      cM,
		aM:      aM,
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
	rC := createRuntimeContext(am.lC, am.cM, am.aM, am.env, am.tplC)
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

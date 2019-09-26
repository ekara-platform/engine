package action

import (
	"fmt"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/engine/util"
)

type (
	//ActionManager is the manager of all action available into the engine
	ActionManager interface {
		// Run executes an engine action
		Run(id ActionID) error
	}

	actionManager struct {
		// launchContext from the engine holding this action manager
		lC util.LaunchContext
		// the component manager
		cM component.ComponentManager
		// the ansible manager
		aM ansible.AnsibleManager
		// available actions
		actions map[ActionID]Action
	}
)

//CreateActionManager initializes the action manager and its content
func CreateActionManager(lC util.LaunchContext, cM component.ComponentManager, aM ansible.AnsibleManager) ActionManager {
	am := actionManager{
		lC:      lC,
		cM:      cM,
		aM:      aM,
		actions: make(map[ActionID]Action),
	}
	for _, a := range allActions() {
		am.actions[a.id] = a
	}
	return &am
}

func (am *actionManager) empty() bool {
	return len(am.actions) == 0
}

//get returns the action corresponding to the given id.
func (am *actionManager) get(id ActionID) (Action, error) {
	if val, ok := am.actions[id]; ok {
		return val, nil
	}
	return Action{}, fmt.Errorf("Unsupported action")
}

//Run launches the action corresponding to the given id.
func (am *actionManager) Run(id ActionID) error {
	a, e := am.get(id)
	if e != nil {
		return e
	}

	am.lC.Log().Printf(LogLaunchingAction, a.name)
	report, e := a.run(am)
	if e != nil {
		return e
	}

	loc, e := writeReport(*report, am.lC.Ef().Output)
	if e != nil {
		return e
	}
	am.lC.Log().Printf(LogReportWritten, loc)

	return nil
}

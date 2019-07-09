package engine

import (
	"fmt"

	"github.com/ekara-platform/engine/ansible"
)

type (
	//ActionManager is the manager of all action available into the engine
	ActionManager struct {
		// available actions
		actions map[ActionID]Action
	}
)

//CreateActionManager initializes the action manager and its content
func CreateActionManager() ActionManager {
	am := ActionManager{
		actions: make(map[ActionID]Action),
	}
	for _, a := range InitActions() {
		am.actions[a.id] = a
	}
	return am
}

func (am ActionManager) empty() bool {
	return len(am.actions) == 0
}

//get returns the action corresponding to the given id.
func (am ActionManager) get(id ActionID) (Action, error) {
	if val, ok := am.actions[id]; ok {
		return val, nil
	}
	return Action{}, fmt.Errorf("Unsupported action")
}

//Run launches the action corresponding to the given id.
//The method will panic if the required action is missing.
func (am ActionManager) Run(id ActionID, lC LaunchContext) {
	lC.Ekara().ComponentManager().Ensure()
	a, e := am.get(id)
	if e != nil {
		panic(e)
	}
	// Initialization of the runtime context
	rC := &runtimeContext{}
	rC.buffer = make(map[string]ansible.Buffer)

	lC.Log().Printf(LogLaunchingAction, a.name)
	report, e := a.run(am, lC, rC)
	if e != nil {
		// Do something with the error here
		panic(e)
	}
	e = writeReport(*report)
	if e != nil {
		// DO something with the error here
	}
}

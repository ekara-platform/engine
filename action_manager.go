package engine

import (
	"fmt"

	"github.com/ekara-platform/engine/ansible"
)

type (
	// The manager of all action available into the engine
	actionManager struct {
		// available actions
		actions map[ActionId]action
	}
)

//CreateActionManager initializes the action manager and its content
func CreateActionManager() actionManager {
	am := actionManager{
		actions: make(map[ActionId]action),
	}
	for _, a := range InitActions() {
		am.actions[a.id] = a
	}
	return am
}

func (am actionManager) empty() bool {
	return len(am.actions) == 0
}

//get returns the action corresponding to the given id.
func (am actionManager) get(id ActionId) (action, error) {
	if val, ok := am.actions[id]; ok {
		return val, nil
	}
	return action{}, fmt.Errorf("Unsupported action")
}

//Run runs the action corresponding to the given id.
func (am actionManager) Run(id ActionId, lC LaunchContext) {
	a, e := am.get(id)
	if e != nil {
		panic(e)
	}
	// Initialization of the runtime context
	rC := &runtimeContext{}
	rC.buffer = make(map[string]ansible.Buffer)

	lC.Log().Printf(LOG_LAUNCHING_ACTION, a.name)
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

package engine

import (
	"fmt"
	"path/filepath"

	"github.com/GroupePSA/componentizer"
	"github.com/ekara-platform/engine/action"
	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/model"
	"github.com/ekara-platform/engine/util"
)

// Ekara is the facade used to manage environments.
type Ekara interface {
	Init(repo componentizer.Repository) error
	Environment() model.Environment
	Execute(id action.ActionID) (action.Result, error)
}

type engine struct {
	// Base context
	lC        util.LaunchContext
	directory string

	// Environment and its data
	environment model.Environment
	tplC        *model.TemplateContext

	// Available actions
	actions map[action.ActionID]action.Action

	// Subsystems
	componentManager componentizer.ComponentManager
	ansibleManager   ansible.Manager
}

// Create creates an engine environment descriptor based on the provided location.
//
//	Parameters:
//		lC: the launch context
//		workDir: the directory where the engine will do its work
func Create(lC util.LaunchContext, workDir string) Ekara {
	eng := engine{
		lC:        lC,
		directory: filepath.Clean(workDir),
		tplC:      model.CreateTemplateContext(lC.ExternalVars()),
		actions:   make(map[action.ActionID]action.Action),
	}

	// Register actions
	for _, a := range action.All() {
		eng.actions[a.Id] = a
	}

	// Initialize subsystems
	eng.componentManager = componentizer.CreateComponentManager(eng.lC.Log(), filepath.Join(eng.directory, "components"))
	eng.ansibleManager = ansible.CreateAnsibleManager(eng.lC, eng.componentManager)

	return &eng
}

func (eng *engine) Init(repo componentizer.Repository) error {
	m, err := eng.componentManager.Init(model.CreateComponent(model.MainComponentId, repo), eng.tplC)
	if err != nil {
		return err
	}
	eng.environment = m.(model.Environment)
	eng.tplC.Model = m.(model.Environment)
	return nil
}

func (eng *engine) Environment() model.Environment {
	return eng.environment
}

func (eng *engine) Execute(id action.ActionID) (action.Result, error) {
	rC := action.CreateRuntimeContext(eng.lC, eng.componentManager, eng.ansibleManager, eng.environment, eng.tplC)
	r := &action.ExecutionReport{}

	// Execute the action chain
	res, e := eng.execute(id, rC, r)
	if e != nil {
		return nil, e
	}

	// Write the report
	loc, e := r.Write(eng.lC.Ef().Output)
	if e != nil {
		return nil, e
	}
	eng.lC.Log().Printf("The execution report file has been written in %s\n", loc)

	return res, nil
}

func (eng *engine) execute(id action.ActionID, rC *action.RuntimeContext, report *action.ExecutionReport) (action.Result, error) {
	a, ok := eng.actions[id]
	if !ok {
		return nil, fmt.Errorf("unsupported action")
	}

	if a.DependsOn != action.NilActionID {
		// Execute dependent actions but drop the result
		_, e := eng.execute(a.DependsOn, rC, report)
		if e != nil {
			return nil, e
		}
	}

	// Execute the action
	rep, res := a.Execute(rC)
	if rep.Error != nil {
		return nil, rep.Error
	} else {
		report.Aggregate(rep)
	}
	return res, nil
}

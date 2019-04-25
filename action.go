package engine

import (
	"strconv"
)

type (
	// ActionID represents the id of an action available on an environment
	ActionID int

	//Action represents an action available on an environment
	Action struct {
		// The action id
		id ActionID
		// The action id  on which this action depends
		dependsOn ActionID
		// The name of the action
		name string
		// The action steps
		steps []step
	}
)

const (
	// ActionFailID identifies the action of marking a StepResults
	//with an error in case of validation error into the descriptor
	ActionFailID ActionID = iota
	//ActionReportID identifies the action of reading an existing execution report
	ActionReportID
	//ActionCreateID identifies the action of creation environment's machines
	ActionCreateID
	//ActionInstallID identifies the action of installing the environment orchestrator
	ActionInstallID
	//ActionDeployID identifies the action of deveploying the environment stacks
	ActionDeployID
	//ActionCheckID identifies the action of validating the environment descriptor
	ActionCheckID
	//ActionDumpID identifies the acion of dumping the environment descriptor content
	ActionDumpID
	//ActionUpdateID identifies the action of updating of an environment
	ActionUpdateID
	//ActionDeleteID identifies the action of deleting an environment
	ActionDeleteID
	//ActionNilID identifies no action, used to indicate that an action depends on nothing
	ActionNilID
	//ActionRegisterID identifies the action of registering an environment
	// through the api once it has been create
	ActionRegisterID
)

// String returns the string representation of the action id
func (a ActionID) String() string {
	return strconv.Itoa(int(a))
}

//InitActions initializes all actions available into the engien
func InitActions() []Action {
	r := make([]Action, 0)
	r = append(r, Action{ActionFailID, ActionNilID, "FailOnError", failOnEkaraErrorSteps})
	r = append(r, Action{ActionReportID, ActionFailID, "Report", reportSteps})
	r = append(r, Action{ActionCreateID, ActionReportID, "Create", createSteps})
	r = append(r, Action{ActionInstallID, ActionCreateID, "Install", installSteps})
	r = append(r, Action{ActionDeployID, ActionInstallID, "Deploy", deploySteps})
	r = append(r, Action{ActionCheckID, ActionNilID, "Check", checkSteps})
	r = append(r, Action{ActionDumpID, ActionCheckID, "Dump", dumpSteps})
	return r
}

// run runs an action for the given action manager and contexts
func (a Action) run(m ActionManager, lC LaunchContext, rC *runtimeContext) (*ExecutionReport, error) {
	r := &ExecutionReport{
		Context: lC,
	}

	if a.dependsOn != ActionNilID {
		d, e := m.get(a.dependsOn)
		if e != nil {
			return r, e
		}
		// Run the dependent action
		rep, e := d.run(m, lC, rC)
		if e != nil {
			return r, e
		}
		r.aggregate(*rep)
		if rep.Error != nil {
			return r, nil
		}

	}

	lC.Log().Printf(LogRunningAction, a.name)

	// Run the actions steps
	rep := launch(a.steps, lC, rC)
	r.aggregate(rep)
	return r, nil
}

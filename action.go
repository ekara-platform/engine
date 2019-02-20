package engine

import (
	"strconv"
)

type (
	// ActionId represents the id of an action available on an environment
	ActionId int

	// action represents an action available on an environment
	action struct {
		// The action id
		id ActionId
		// The action id  on which this action depends
		dependsOn ActionId
		// The name of the action
		name string
		// The action steps
		steps []step
	}
)

const (
	// Mark a StepResults with an error in case of validation error into the descriptor
	ActionFailId ActionId = iota
	// Read an existing execution report
	ActionReportId
	// Creation of an environment's machines
	ActionCreateId
	// Installation of the environment orchestrator
	ActionInstallId
	// Deployment of the environment stacks
	ActionDeployId
	// Validation of the environment descriptor
	ActionCheckId
	// Dump of the environment descriptor
	ActionDumpId
	// Update of an environment
	ActionUpdateId
	// Deletion of an environment
	ActionDeleteId
	// Nil action Id, used to indicate that an action depends  on nothing
	ActionNilId
)

// String returns the string representation of the action id
func (a ActionId) String() string {
	return strconv.Itoa(int(a))
}

//InitActions initializes all actions available into the engien
func InitActions() []action {
	r := make([]action, 0)
	r = append(r, action{ActionFailId, ActionNilId, "FailOnError", failOnEkaraErrorSteps})
	r = append(r, action{ActionReportId, ActionFailId, "Report", reportSteps})
	r = append(r, action{ActionCreateId, ActionReportId, "Create", createSteps})
	r = append(r, action{ActionInstallId, ActionCreateId, "Install", installSteps})
	r = append(r, action{ActionDeployId, ActionInstallId, "Deploy", deploySteps})
	r = append(r, action{ActionCheckId, ActionNilId, "Check", checkSteps})
	r = append(r, action{ActionDumpId, ActionCheckId, "Dump", dumpSteps})
	return r
}

// run runs an action for the given action manager and contexts
func (a action) run(m actionManager, lC LaunchContext, rC *runtimeContext) (error, *ExecutionReport) {
	r := &ExecutionReport{
		Context: lC,
	}

	if a.dependsOn != ActionNilId {
		d, e := m.get(a.dependsOn)
		if e != nil {
			return e, r
		}
		// Run the dependent action
		e, rep := d.run(m, lC, rC)
		if e != nil {
			return e, r
		}
		r.aggregate(*rep)
		if rep.Error != nil {
			return nil, r
		}

	}

	lC.Log().Printf(LOG_RUNNING_ACTION, a.name)

	// Run the actions steps
	rep := launch(a.steps, lC, rC)
	r.aggregate(rep)
	return nil, r
}

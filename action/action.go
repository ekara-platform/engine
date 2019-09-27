package action

import (
	"strconv"
	"time"
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
	// NilActionID identifies no action, used to indicate that an action depends on nothing
	NilActionID ActionID = iota
	// ValidateActionID identifies the action of validating an environment model.
	ValidateActionID
	// CheckActionID identifies the action of returning the validation results of an environment model.
	CheckActionID
	// FailOnErrorActionID identifies the action of failing if the environment model has validation errors.
	FailOnErrorActionID
	// DumpActionID identifies the acion of dumping the effective environment model.
	DumpActionID
	// ApplyActionID identifies the action of applying a descriptor to reality
	ApplyActionID
	// DestroyActionID identifies the action of destroying an environment.
	DestroyActionID
)

// String returns the string representation of the action id
func (a ActionID) String() string {
	return strconv.Itoa(int(a))
}

func allActions() []Action {
	r := make([]Action, 0)
	r = append(r, applyAction)
	r = append(r, checkAction)
	r = append(r, dumpAction)
	r = append(r, validateAction)
	return r
}

// run runs an action for the given action manager and contexts
func (a Action) run(am *actionManager) (*ExecutionReport, error) {
	r := &ExecutionReport{}

	if a.dependsOn != NilActionID {
		d, e := am.get(a.dependsOn)
		if e != nil {
			return r, e
		}
		// Run the dependent action
		rep, e := d.run(am)
		if e != nil {
			return r, e
		}
		r.aggregate(*rep)
		if rep.Error != nil {
			return r, nil
		}

	}

	am.lC.Log().Printf(LogRunningAction, a.name)

	// Run the actions steps
	rep := a.launch(CreateRuntimeContext(am.lC, am.cM, am.aM))
	r.aggregate(rep)
	return r, nil
}

// launch runs a slice of step functions
//
// If one step in the slice returns an error then the launch process will stop and
// the cleanup will be invoked on all previously launched steps
func (a Action) launch(rC *runtimeContext) ExecutionReport {
	r := ExecutionReport{}

	cleanups := []Cleanup{}
	for _, f := range a.steps {
		ctx := f(rC)
		for _, sr := range ctx.Results {
			i := int64(sr.ExecutionTime / time.Millisecond)
			if i == 0 {
				sr.ExecutionTime, _ = time.ParseDuration("1ms")
			}

			r.Steps.Results = append(r.Steps.Results, sr)
			r.Steps.TotalExecutionTime = r.Steps.TotalExecutionTime + sr.ExecutionTime

			if sr.cleanUp != nil {
				cleanups = append(cleanups, sr.cleanUp)
			}

			e := sr.error
			if e != nil {
				cleanLaunched(cleanups, rC.lC)
				r.Error = e
				return r
			}
		}
	}

	return r
}

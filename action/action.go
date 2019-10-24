package action

import (
	"time"
)

type (
	// ActionID represents the id of an action available on an environment
	ActionID string

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
	NilActionID ActionID = "NIL"
	// ValidateActionID identifies the action of validating an environment model.
	ValidateActionID = "VALIDATE"
	// CheckActionID identifies the action of returning the validation results of an environment model.
	CheckActionID = "CHECK"
	// DumpActionID identifies the action of dumping the effective environment model.
	DumpActionID = "DUMP"
	// ApplyActionID identifies the action of applying a descriptor to reality
	ApplyActionID = "APPLY"
)

// String returns the string representation of the action id
func (a ActionID) String() string {
	return string(a)
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
func (a Action) run(am manager) (*ExecutionReport, Result, error) {
	r := &ExecutionReport{}

	if a.dependsOn != NilActionID {
		// Obtain the dependent action
		d, e := am.get(a.dependsOn)
		if e != nil {
			return r, nil, e
		}

		// Run the dependent action and aggregate its report (dropping the intermediate result)
		rep, _, e := d.run(am)
		if e != nil {
			return r, nil, e
		}
		r.aggregate(*rep)

		// If the report contains an error return it
		if rep.Error != nil {
			return r, nil, rep.Error
		}
	}

	am.rC.lC.Log().Printf(LogRunningAction, a.name)

	// Run the final action and return its result
	rep, res := a.launch(am.rC)
	r.aggregate(rep)
	return r, res, nil
}

// launch runs the action on the given contenxt
func (a Action) launch(rC *runtimeContext) (ExecutionReport, Result) {
	r := ExecutionReport{}

	cleanups := []Cleanup{}
	var finalRes Result
	for _, f := range a.steps {
		sCs, res := f(rC)
		for _, sr := range sCs.Status {
			i := int64(sr.ExecutionTime / time.Millisecond)
			if i == 0 {
				sr.ExecutionTime, _ = time.ParseDuration("1ms")
			}

			r.Steps.Status = append(r.Steps.Status, sr)
			r.Steps.TotalExecutionTime = r.Steps.TotalExecutionTime + sr.ExecutionTime

			if sr.cleanUp != nil {
				cleanups = append(cleanups, sr.cleanUp)
			}

			e := sr.error
			if e != nil {
				cleanLaunched(cleanups, rC.lC)
				r.Error = e
				return r, nil
			}
		}
		if res != nil {
			finalRes = res
		}
	}
	return r, finalRes
}

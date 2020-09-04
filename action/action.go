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
		Id ActionID
		// The action id  on which this action depends
		DependsOn ActionID
		// The name of the action
		Name string
		// The action steps
		steps []step
	}

	Result interface {
		//IsSuccess returns true id the action execution was successful
		IsSuccess() bool
		//AsJson returns the action returned content as JSON
		AsJson() (string, error)
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
	// DestroyActionID identifies the action of destroying an existing environment
	DestroyActionID = "DESTROY"
)

// String returns the string representation of the action id
func (a ActionID) String() string {
	return string(a)
}

func All() []Action {
	r := make([]Action, 0)
	r = append(r, applyAction)
	r = append(r, destroyAction)
	r = append(r, checkAction)
	r = append(r, dumpAction)
	r = append(r, validateAction)
	return r
}

// launch runs the action on the given context
func (a Action) Execute(rC *RuntimeContext) (ExecutionReport, Result) {
	r := ExecutionReport{}

	cleanups := []Cleanup{}
	var finalRes Result
	for _, f := range a.steps {
		sCs := f(rC)
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
		if rC.result != nil {
			finalRes = rC.result
		}
	}
	return r, finalRes
}

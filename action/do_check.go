package action

import (
	"fmt"
)

var (
	checkAction = Action{
		CheckActionID,
		NilActionID,
		"Check",
		[]step{doCheck},
	}
)

func doCheck(rC *runtimeContext) (StepResults, Result) {
	sc := InitCodeStepResult("Checking if the environment has any validation error", nil, NoCleanUpRequired)
	vErrs := rC.cM.Environment().Validate()
	if vErrs.HasErrors() {
		// in case of validation error we stop
		FailsOnDescriptor(&sc, fmt.Errorf("Descriptor error"), "The descriptor is not valid", nil)
	}
	return sc.Build(), nil
}

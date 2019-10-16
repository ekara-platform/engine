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

	rC.pN.Notify("check", "Checking for environment model problems")

	vErrs := rC.environment.Validate()
	if vErrs.HasErrors() {
		// in case of validation error we stop
		FailsOnDescriptor(&sc, fmt.Errorf("descriptor error"), "The descriptor is not valid", nil)
	}

	rC.pN.Notify("check", "Environment model checked")

	return sc.Build(), nil
}

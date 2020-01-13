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

func doCheck(rC *runtimeContext) StepResults {
	sc := InitCodeStepResult("Checking if the environment has any validation error", nil, NoCleanUpRequired)

	rC.lC.Feedback().Progress("check", "Checking for environment model problems")

	vErrs := rC.environment.Validate()
	if vErrs.HasErrors() {
		// in case of validation error we stop
		FailsOnDescriptor(&sc, fmt.Errorf("descriptor error"), "The descriptor is not valid", nil)
	}

	rC.lC.Feedback().Progress("check", "Environment model checked")

	return sc.Build()
}

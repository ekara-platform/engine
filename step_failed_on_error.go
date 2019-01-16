package engine

import (
	"fmt"

	"github.com/ekara-platform/model"
)

var failOnEkaraErrorSteps = []step{ffailOnEkaraError}

func ffailOnEkaraError(lC LaunchContext, rC *runtimeContext) StepResults {
	sc := InitCodeStepResult("Stopping the process in case of validation errors", nil, NoCleanUpRequired)
	if rC.ekaraError != nil {
		vErrs, ok := rC.ekaraError.(model.ValidationErrors)
		if ok {
			if vErrs.HasErrors() {
				// in case of validation error we stop
				lC.Log().Println(rC.ekaraError)
				FailsOnDescriptor(&sc, rC.ekaraError, fmt.Sprintf(ERROR_PARSING_DESCRIPTOR, rC.ekaraError.Error()), nil)
				goto MoveOut
			}
		} else {
			FailsOnDescriptor(&sc, rC.ekaraError, fmt.Sprintf(ERROR_PARSING_DESCRIPTOR, rC.ekaraError.Error()), nil)
			goto MoveOut
		}
	}
MoveOut:
	return sc.Array()
}

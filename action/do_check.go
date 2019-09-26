package action

import (
	"fmt"

	"github.com/ekara-platform/model"
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
	if rC.ekaraError != nil {
		vErrs, ok := rC.ekaraError.(model.ValidationErrors)
		if ok {
			if vErrs.HasErrors() {
				// in case of validation error we stop
				rC.lC.Log().Println(rC.ekaraError)
				FailsOnDescriptor(&sc, rC.ekaraError, fmt.Sprintf(ErrorParsingDescriptor, rC.ekaraError.Error()), nil)
			}
		} else {
			FailsOnDescriptor(&sc, rC.ekaraError, fmt.Sprintf(ErrorParsingDescriptor, rC.ekaraError.Error()), nil)
		}
	}
	return sc.Array()
}

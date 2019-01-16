package engine

import (
	"fmt"

	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
)

var checkSteps = []step{flogCheck}

func flogCheck(lC LaunchContext, rC *runtimeContext) StepResults {
	sc := InitCodeStepResult("Validating the environment content", nil, NoCleanUpRequired)
	ve := rC.ekaraError
	if ve != nil {
		vErrs, ok := ve.(model.ValidationErrors)
		// if the error is not a "validation error" then we return it
		if !ok {
			FailsOnDescriptor(&sc, ve, fmt.Sprintf(ERROR_PARSING_ENVIRONMENT, ve.Error()), nil)
		} else {
			lC.Log().Printf(ve.Error())
			b, e := vErrs.JSonContent()
			if e != nil {
				FailsOnDescriptor(&sc, e, fmt.Sprintf(ERROR_GENERIC, e), nil)
			}
			// print both errors and warnings into the report file
			path, err := util.SaveFile(lC.Log(), *lC.Ef().Output, VALIDATION_OUTPUT_FILE, b)
			if err != nil {
				// in case of error writing the report file
				FailsOnDescriptor(&sc, err, fmt.Sprintf(ERROR_CREATING_REPORT_FILE, path), nil)
			}

			if vErrs.HasErrors() {
				// in case of validation error we stop
				FailsOnDescriptor(&sc, ve, fmt.Sprintf(ERROR_PARSING_DESCRIPTOR, ve.Error()), nil)
			}
		}
	} else {
		lC.Log().Printf(LOG_VALIDATION_SUCCESSFUL)
	}
	return sc.Array()
}

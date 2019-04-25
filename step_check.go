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
			FailsOnDescriptor(&sc, ve, fmt.Sprintf(ErrorParsingEnvironment, ve.Error()), nil)
		} else {
			lC.Log().Printf("%s\n", ve.Error())
			b, e := vErrs.JSonContent()
			if e != nil {
				FailsOnCode(&sc, e, fmt.Sprintf(ErrorGeneric, e), nil)
			}
			// print both errors and warnings into the report file
			path, err := util.SaveFile(lC.Log(), *lC.Ef().Output, ValidationOutputFile, b)
			if err != nil {
				// in case of error writing the report file
				FailsOnCode(&sc, err, fmt.Sprintf(ErrorCreatingReportFile, path), nil)
			}

			if vErrs.HasErrors() {
				// in case of validation error we stop
				FailsOnDescriptor(&sc, ve, fmt.Sprintf(ErrorParsingDescriptor, ve.Error()), nil)
			}
		}
	} else {
		lC.Log().Printf(LogValidationSuccessful)
	}
	return sc.Array()
}

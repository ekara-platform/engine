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

	rC.lC.Feedback().Progress("check", "Checking model and components")

	// Validate the descriptor
	vErrs := rC.environment.Validate()
	if vErrs.HasErrors() {
		rC.lC.Feedback().Error("Environment model is not valid")
		FailsOnModel(&sc, fmt.Errorf("model error"), "Environment model is not valid", nil)
		return sc.Build()
	} else {
		// Validate all components
		hasErrors := false
		for cName, comp := range rC.environment.Platform().Components {
			uc, err := rC.cF.Use(comp, rC.tplC)
			if err != nil {
				rC.lC.Feedback().Error("Component %s is not valid: %s", cName, err.Error())
				hasErrors = true
			} else {
				uc.Release()
			}
		}
		if hasErrors {
			FailsOnComponent(&sc, fmt.Errorf("at least one component is not valid"), "", nil)
			return sc.Build()
		}
	}

	rC.lC.Feedback().Progress("check", "Model and components checked")
	return sc.Build()
}

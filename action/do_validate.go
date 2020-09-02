package action

import (
	"errors"

	"github.com/ekara-platform/engine/model"
)

var (
	validateAction = Action{
		ValidateActionID,
		NilActionID,
		"Validate",
		[]step{doValidate},
	}
)

//ValidateResult contains validation errors ready to be serialized
type ValidateResult struct {
	model.ValidationErrors
}

//IsSuccess returns true id the validate execution was successful
func (v ValidateResult) IsSuccess() bool {
	return !v.HasErrors()
}

//FromJson fills an action returned content from a JSON content
func (v ValidateResult) FromJson(s string) error {
	return errors.New("not implemented")
}

//AsJson returns the validation content as JSON
func (v ValidateResult) AsJson() (string, error) {
	return "", errors.New("not implemented")
}

func doValidate(rC *runtimeContext) StepResults {
	sc := InitCodeStepResult("Validating the environment content", nil, NoCleanUpRequired)
	rC.result = ValidateResult{ValidationErrors: rC.environment.Validate()}
	return sc.Build()
}

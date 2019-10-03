package action

import (
	"fmt"

	"github.com/ekara-platform/model"
)

const (
	validationOutputFile = "validation.json"
)

var (
	validateAction = Action{
		ValidateActionID,
		NilActionID,
		"Validate",
		[]step{doValidate},
	}
)

type ValidateResult struct {
	vErrs model.ValidationErrors
}

func (v ValidateResult) IsSuccess() bool {
	return !v.vErrs.HasErrors()
}

func (v ValidateResult) AsJson() (string, error) {
	b, err := v.vErrs.JSonContent()
	if err != nil {
		return "", err
	}
	return fmt.Sprint(b), nil
}

func (v ValidateResult) AsYaml() (string, error) {
	return v.AsJson()
}

func (v ValidateResult) AsPlainText() ([]string, error) {
	errors := make([]string, 0)
	warnings := make([]string, 0)
	for _, vErr := range v.vErrs.Errors {
		if vErr.ErrorType == model.Error {
			errors = append(errors, "ERROR "+vErr.Message)
		} else {
			warnings = append(warnings, "WARN  "+vErr.Message)
		}
	}
	return append(errors, warnings...), nil
}

func doValidate(rC *runtimeContext) (StepResults, Result) {
	sc := InitCodeStepResult("Validating the environment content", nil, NoCleanUpRequired)
	err := rC.ekaraError
	if err != nil {
		vErrs, ok := err.(model.ValidationErrors)
		// if the error is not a "validation error" then we return it
		if !ok {
			FailsOnDescriptor(&sc, err, fmt.Sprintf(ErrorParsingEnvironment, err.Error()), nil)
			return sc.Build(), nil
		} else {
			return sc.Build(), ValidateResult{vErrs: vErrs}
		}
	} else {
		return sc.Build(), ValidateResult{}
	}
}

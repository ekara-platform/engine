package action

import (
	"errors"
	"github.com/ekara-platform/engine/model"
	"gopkg.in/yaml.v2"
)

var (
	dumpAction = Action{
		DumpActionID,
		NilActionID,
		"Dump",
		[]step{doDump},
	}
)

//DumpResult contains the built environment ready to be serialized
type DumpResult struct {
	Env model.Environment
}

//IsSuccess returns true id the dump execution was successful
func (r DumpResult) IsSuccess() bool {
	return true
}

//FromJson fills an action returned content from a JSON content
func (r DumpResult) FromJson(s string) error {
	return errors.New("not implemented")
}

//AsJson returns the dump content as JSON
func (r DumpResult) AsJson() (string, error) {
	return "", errors.New("not implemented")
}

//AsJson returns the dump content as JSON
func (r DumpResult) AsYaml() (string, error) {
	envYaml, err := yaml.Marshal(r.Env)
	if err != nil {
		return "", err
	}
	return string(envYaml), nil
}

func doDump(rC *runtimeContext) StepResults {
	sc := InitCodeStepResult("Retrieving aggregated environment model", nil, NoCleanUpRequired)
	rC.result = DumpResult{Env: rC.environment}
	return sc.Build()
}

package action

import (
	"encoding/json"
	"strings"

	"github.com/ekara-platform/model"

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
	env *model.Environment
}

//IsSuccess returns true id the dump execution was successful
func (r DumpResult) IsSuccess() bool {
	return r.env != nil
}

//AsJson returns the dump content as JSON
func (r DumpResult) AsJson() (string, error) {
	envJson, err := json.MarshalIndent(r.env, "", "    ")
	if err != nil {
		return "", err
	}
	return string(envJson), nil
}

//AsYaml returns the dump content as YAML
func (r DumpResult) AsYaml() (string, error) {
	envYaml, err := yaml.Marshal(r.env)
	if err != nil {
		return "", err
	}
	return string(envYaml), nil
}

//AsPlainText returns the dump content as plain text
func (r DumpResult) AsPlainText() ([]string, error) {
	yamlEnv, err := r.AsYaml()
	if err != nil {
		return []string{}, err
	}
	return strings.Split(yamlEnv, "\n"), nil
}

func doDump(rC *runtimeContext) (StepResults, Result) {
	sc := InitCodeStepResult("Retrieving aggregated environment model", nil, NoCleanUpRequired)
	return sc.Build(), DumpResult{env: rC.environment}
}

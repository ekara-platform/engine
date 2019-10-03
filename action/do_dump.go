package action

import (
	"encoding/json"
	"github.com/ekara-platform/model"
	"strings"

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

type DumpResult struct {
	env *model.Environment
}

func (r DumpResult) IsSuccess() bool {
	return r.env != nil
}

func (r DumpResult) AsJson() (string, error) {
	envJson, err := json.MarshalIndent(r.env, "", "    ")
	if err != nil {
		return "", err
	}
	return string(envJson), nil
}

func (r DumpResult) AsYaml() (string, error) {
	envYaml, err := yaml.Marshal(r.env)
	if err != nil {
		return "", err
	}
	return string(envYaml), nil
}

func (r DumpResult) AsPlainText() ([]string, error) {
	yamlEnv, err := r.AsYaml()
	if err != nil {
		return []string{}, err
	}
	return strings.Split(yamlEnv, "\n"), nil
}

func doDump(rC *runtimeContext) (StepResults, Result) {
	sc := InitCodeStepResult("Retrieving aggregated environment model", nil, NoCleanUpRequired)
	return sc.Build(), DumpResult{env: rC.cM.Environment()}
}

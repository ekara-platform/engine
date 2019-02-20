package engine

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/ekara-platform/engine/util"
)

var dumpSteps = []step{fdump}

func fdump(lC LaunchContext, rC *runtimeContext) StepResults {
	sc := InitCodeStepResult("Dumping the environment content", nil, NoCleanUpRequired)

	environmentYaml, err := yaml.Marshal(lC.Ekara().ComponentManager().Environment())
	if err != nil {
		FailsOnDescriptor(&sc, err, fmt.Sprintf("Error marshalling the environment Yaml content:", err.Error()), nil)
	}

	path, err := util.SaveFile(lC.Log(), *lC.Ef().Output, DUMP_YAML_OUTPUT_FILE, environmentYaml)
	if err != nil {
		// in case of error writing the yaml dump
		FailsOnDescriptor(&sc, err, fmt.Sprintf(ERROR_CREATING_REPORT_FILE, path), nil)
	}

	environmentYaml, err = yaml.Marshal(lC.Ekara().ComponentManager().Environment().OriginalEnv)
	if err != nil {
		FailsOnDescriptor(&sc, err, fmt.Sprintf("Error marshalling the environment source Yaml content:", err.Error()), nil)
	}

	path, err = util.SaveFile(lC.Log(), *lC.Ef().Output, DUMP_SOURCE_YAML_OUTPUT_FILE, environmentYaml)
	if err != nil {
		// in case of error writing the yaml dump
		FailsOnDescriptor(&sc, err, fmt.Sprintf(ERROR_CREATING_REPORT_FILE, path), nil)
	}

	return sc.Array()
}

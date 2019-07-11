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
		FailsOnDescriptor(&sc, err, fmt.Sprintf("Error marshalling the environment Yaml content: %s", err.Error()), nil)
	}

	path, err := util.SaveFile(lC.Log(), *lC.Ef().Output, DumpYamlOutputFile, environmentYaml)
	if err != nil {
		// in case of error writing the yaml dump
		FailsOnDescriptor(&sc, err, fmt.Sprintf(ErrorCreatingDumpFile, path), nil)
	}

	return sc.Array()
}

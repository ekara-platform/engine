package action

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/ekara-platform/engine/util"
)

const (
	dumpYamlFile = "dump.yaml"
)

var (
	dumpAction = Action{
		DumpActionID,
		NilActionID,
		"Dump",
		[]step{doDump},
	}
)

func doDump(rC *runtimeContext) StepResults {
	sc := InitCodeStepResult("Dumping the effective environment model in "+dumpYamlFile, nil, NoCleanUpRequired)

	environmentYaml, err := yaml.Marshal(rC.cM.Environment())
	if err != nil {
		FailsOnDescriptor(&sc, err, fmt.Sprintf("Error marshalling the environment YAML: %s", err.Error()), nil)
	}

	path, err := util.SaveFile(rC.lC.Ef().Output, dumpYamlFile, environmentYaml)
	if err != nil {
		// in case of error writing the yaml dump
		FailsOnDescriptor(&sc, err, fmt.Sprintf(ErrorCreatingDumpFile, path), nil)
	}

	return sc.Array()
}

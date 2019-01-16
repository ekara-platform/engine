package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ekara-platform/engine/util"
)

var reportSteps = []step{freport}

// freport reads the content of the eventually existing report file
func freport(lC LaunchContext, rC *runtimeContext) StepResults {
	sc := InitCodeStepResult("Validating the environment content", nil, NoCleanUpRequired)
	ok := lC.Ef().Output.Contains(REPORT_OUTPUT_FILE)
	if ok {
		lC.Log().Println("A report file from a previous execution has been located")
		b, err := ioutil.ReadFile(util.JoinPaths(lC.Ef().Output.Path(), REPORT_OUTPUT_FILE))
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf(ERROR_READING_REPORT, REPORT_OUTPUT_FILE, err.Error()), nil)
			goto MoveOut
		}

		report := ReportFileContent{}

		err = json.Unmarshal(b, &report)
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf(ERROR_UNMARSHALLING_REPORT, REPORT_OUTPUT_FILE, err.Error()), nil)
			goto MoveOut
		}
		rC.report = report
	} else {
		lC.Log().Println("Unable to locate a report file from a previous execution")
	}
MoveOut:
	return sc.Array()
}

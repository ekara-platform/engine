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
	sc := InitCodeStepResult("Reporting the execution details", nil, NoCleanUpRequired)
	ok := lC.Ef().Output.Contains(REPORT_OUTPUT_FILE)
	if ok {
		lC.Log().Println("A report file from a previous execution has been located")
		b, err := ioutil.ReadFile(util.JoinPaths(lC.Ef().Output.Path(), REPORT_OUTPUT_FILE))
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf(ERROR_READING_REPORT, REPORT_OUTPUT_FILE, err.Error()), nil)
			goto MoveOut
		}
		reportNoTime := ReportFileContentNoTime{}

		err = json.Unmarshal(b, &reportNoTime)
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf(ERROR_UNMARSHALLING_REPORT, REPORT_OUTPUT_FILE, err.Error()), nil)
			goto MoveOut
		}
		report := ReportFileContent{}
		report.Results = make([]StepResult, 0)
		for _, v := range reportNoTime.Results {
			report.Results = append(report.Results, v.ToStepResult())
		}
		rC.report = report

	} else {
		lC.Log().Println("Unable to locate a report file from a previous execution")
	}
MoveOut:
	return sc.Array()
}

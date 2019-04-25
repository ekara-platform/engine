package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ekara-platform/engine/util"
)

var reportSteps = []step{freport}

//freport reads the content of the eventually existing report file
func freport(lC LaunchContext, rC *runtimeContext) StepResults {
	sc := InitCodeStepResult("Reporting the execution details", nil, NoCleanUpRequired)
	ok := lC.Ef().Output.Contains(ReportOutputFile)
	if ok {
		lC.Log().Println("A report file from a previous execution has been located")
		b, err := ioutil.ReadFile(util.JoinPaths(lC.Ef().Output.Path(), ReportOutputFile))
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf(ErrorReadingReport, ReportOutputFile, err.Error()), nil)
			goto MoveOut
		}
		reportNoTime := ReportFileContentNoTime{}

		err = json.Unmarshal(b, &reportNoTime)
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf(ErrorUnmarshallingReport, ReportOutputFile, err.Error()), nil)
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

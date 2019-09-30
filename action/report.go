package action

import (
	"encoding/json"
	"time"

	"github.com/ekara-platform/engine/util"
)

const (
	reportOutputFile = "report.json"
)

type (
	//ExecutionReport contains all the execution details of what has been executed by the engine.
	ExecutionReport struct {
		//Error contains the error which stopped the engine execution or nil if everything was fine
		Error error
		//Steps represent the details of each step executed by the engine
		Steps StepResults
	}

	//ReportFileContent is used to marchal the execution report
	ReportFileContent struct {
		Results []StepResult
	}

	//ReportFailures contains the steps who failed during the engine execution
	ReportFailures struct {
		playBookFailures []StepResult
		otherFailures    []StepResult
	}
)

func writeReport(rep ExecutionReport, path util.FolderPath) (string, error) {
	loc, e := rep.generate(path)
	if e != nil {
		return "", e
	}
	if rep.Error != nil {
		return "", rep.Error
	}
	return loc, nil
}

// Content returns the json representation of the report steps
func (er ExecutionReport) Content() (b []byte, e error) {
	b, e = json.MarshalIndent(&er.Steps, "", "    ")
	return
}

func (er *ExecutionReport) aggregate(r ExecutionReport) {
	if r.Error != nil {
		er.Error = r.Error
	}
	res := make([]StepResult, 0)
	var totalDuration time.Duration

	for _, v := range er.Steps.Status {
		i := int64(v.ExecutionTime / time.Millisecond)
		if i == 0 {
			v.ExecutionTime, _ = time.ParseDuration("1ms")
		}
		res = append(res, v)
		totalDuration = totalDuration + v.ExecutionTime
	}

	for _, v := range r.Steps.Status {
		i := int64(v.ExecutionTime / time.Millisecond)
		if i == 0 {
			v.ExecutionTime, _ = time.ParseDuration("1ms")
		}
		res = append(res, v)
		totalDuration = totalDuration + v.ExecutionTime
	}

	i := int64(totalDuration / time.Millisecond)
	if i == 0 {
		totalDuration, _ = time.ParseDuration("1ms")
	}

	er.Steps = StepResults{
		Status:             res,
		TotalExecutionTime: totalDuration,
	}
}

func (er ExecutionReport) generate(path util.FolderPath) (string, error) {
	b, err := er.Content()
	if err != nil {
		return "", err
	}
	return util.SaveFile(path, reportOutputFile, b)
}

func (rfc ReportFileContent) hasFailure() (bool, ReportFailures) {
	r := ReportFailures{}

	for _, v := range rfc.Results {
		if v.Status == STEP_STATUS_FAILURE {

			switch v.FailureCause {
			case codeFailure, descriptorFailure:
				r.otherFailures = append(r.otherFailures, v)
				break

			case playBookFailure:
				r.playBookFailures = append(r.playBookFailures, v)
				break

			}
		}
	}
	return len(r.playBookFailures) > 0 || len(r.otherFailures) > 0, r
}

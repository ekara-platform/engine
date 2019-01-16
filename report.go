package engine

import (
	"encoding/json"

	"github.com/ekara-platform/engine/util"
)

type (
	ExecutionReport struct {
		Error   error
		Steps   StepResults
		Context LaunchContext
	}

	ReportFileContent struct {
		Results []StepResult
	}

	ReportFailures struct {
		playBookFailures []StepResult
		otherFailures    []StepResult
	}
)

func writeReport(rep ExecutionReport) error {
	loc, e := rep.generate()
	if e != nil {
		return e
	}
	rep.Context.Log().Printf(LOG_REPORT_WRITTEN, loc)
	if rep.Error != nil {
		return rep.Error
	}
	return nil
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
	for _, v := range er.Steps.Results {
		res = append(res, v)
	}

	for _, v := range r.Steps.Results {
		res = append(res, v)
	}
	er.Steps = StepResults{
		Results: res,
	}
}

func (er ExecutionReport) generate() (string, error) {
	b, err := er.Content()
	if err != nil {
		return "", err
	}
	return util.SaveFile(er.Context.Log(), *er.Context.Ef().Output, REPORT_OUTPUT_FILE, b)
}

func (rfc ReportFileContent) hasFailure() (bool, ReportFailures) {
	r := ReportFailures{}

	for _, v := range rfc.Results {
		if v.Status == STEP_STATUS_FAILURE {

			switch v.FailureCause {
			case CODE_FAILURE, DESCRIPTOR_FAILURE:
				r.otherFailures = append(r.otherFailures, v)
				break

			case PLAYBOOK_FAILURE:
				r.playBookFailures = append(r.playBookFailures, v)
				break

			}
		}
	}
	return len(r.playBookFailures) > 0 || len(r.otherFailures) > 0, r
}

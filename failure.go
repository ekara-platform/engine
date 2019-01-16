package engine

import (
	"encoding/json"
)

type (
	// Type used to identify the cause of the failure
	failureCause string

	// Detail of a playbook execution failure
	playBookFailureDetail struct {
		// The name of the playbook
		Playbook string
		// The compoment holding the playbook
		Compoment string
		// The ansible return code
		Code int
	}
)

const (
	// Failure due to an Ekara code execution error
	NOT_IMPLEMENTED_FAILURE failureCause = "NotImplemented"

	// Failure due to an Ekara code execution error
	CODE_FAILURE failureCause = "Code"
	// Failure due to a invalid environment descriptor content
	DESCRIPTOR_FAILURE failureCause = "Descriptor"
	// Failure due to a playbook execution error
	PLAYBOOK_FAILURE failureCause = "PlayBook"
)

var FailsOnNotImplemented = failOn(NOT_IMPLEMENTED_FAILURE)
var FailsOnCode = failOn(CODE_FAILURE)
var FailsOnDescriptor = failOn(DESCRIPTOR_FAILURE)
var FailsOnPlaybook = failOn(PLAYBOOK_FAILURE)

func failOn(fc failureCause) func(sr *StepResult, err error, detail string, content interface{}) {
	return func(sr *StepResult, err error, detail string, content interface{}) {
		sr.FailureCause = fc
		sr.error = err
		if err != nil {
			sr.ErrorMessage = err.Error()
		}
		sr.ReadableMessage = detail
		sr.Status = STEP_STATUS_FAILURE
		if content != nil {
			bs, _ := json.Marshal(content)
			sr.RawContent = string(bs)
		}
	}
}

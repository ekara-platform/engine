package action

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
		// The component holding the playbook
		Component string
		// The ansible return code
		Code int
	}
)

const (
	notImplementedFailure failureCause = "NotImplemented"
	codeFailure           failureCause = "Code"
	descriptorFailure     failureCause = "Descriptor"
	playBookFailure       failureCause = "PlayBook"
)

// FailsOnNotImplemented allows to create a failure on an execution step
// because of a missing implementation
var FailsOnNotImplemented = failOn(notImplementedFailure)

// FailsOnCode allows to create a failure on an execution step
// because of an error returned by ekara, other than errors generated
// by the processing of the descriptor of by a playbook execution
var FailsOnCode = failOn(codeFailure)

// FailsOnDescriptor allows to create a failure on an execution step
// because of an invalid descriptor
var FailsOnDescriptor = failOn(descriptorFailure)

// FailsOnPlaybook allows to create a failure on an execution step
// because of an error return by a playbook execution
var FailsOnPlaybook = failOn(playBookFailure)

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

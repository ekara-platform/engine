package engine

import (
	"encoding/json"
	"time"

	"github.com/ekara-platform/model"
)

type (
	//stepResultStatus is the type used to identify final status of the step
	stepResultStatus string
	//stepResultContext is the type used to identify the step context,
	// or if you prefer what the step was doing
	stepResultContext string

	//StepResultsNotime is a lighter representation of a  chain of steps execution results
	//It is used to unmarchal step results
	StepResultsNotime struct {
		Results []StepResultNoTime
	}

	//StepResults represents a chain of steps execution results
	//It is used to marchal step results
	StepResults struct {
		Results            []StepResult
		TotalExecutionTime time.Duration
	}

	StepResultNoTime struct {
		StepName        string
		AppliedToType   string `json:",omitempty"`
		AppliedToName   string `json:",omitempty"`
		Status          stepResultStatus
		Context         stepResultContext
		FailureCause    failureCause `json:",omitempty"`
		ErrorMessage    string       `json:",omitempty"`
		ReadableMessage string       `json:",omitempty"`
		RawContent      interface{}  `json:",omitempty"`
	}

	StepResult struct {
		StepName        string
		AppliedToType   string `json:",omitempty"`
		AppliedToName   string `json:",omitempty"`
		Status          stepResultStatus
		Context         stepResultContext
		FailureCause    failureCause `json:",omitempty"`
		error           error
		ErrorMessage    string      `json:",omitempty"`
		ReadableMessage string      `json:",omitempty"`
		RawContent      interface{} `json:",omitempty"`
		cleanUp         Cleanup
		startedAt       time.Time
		ExecutionTime   time.Duration
	}

	// step represents a sinlge step used to compose a process executed by the installer
	step func(lC LaunchContext, rC *runtimeContext) StepResults
)

const (

	// Step execution successful
	STEP_STATUS_FAILURE stepResultStatus = "Failure"
	// Step execution failed
	STEP_STATUS_SUCCESS stepResultStatus = "Success"
	// The step belongs internally to Ekara
	STEP_CONTEXT_CODE stepResultContext = "Ekara execution"
	// The step deals with the environment descriptor
	STEP_CONTEXT_DESCRIPTOR stepResultContext = "Environment descriptor content"
	//The step deals with the parameter file used to fill descriptor variables
	STEP_CONTEXT_PARAMETER_FILE stepResultContext = "Environment parameter file"
	//The step purpose is to launch an ansible playbook
	STEP_CONTEXT_PLABOOK stepResultContext = "Playbook execution"
	// The step purpose is to launch a hook through an ansible playbook
	STEP_CONTEXT_HOOK_PLABOOK stepResultContext = "Hook execution"
)

// InitCodeStepResult initializes a step result to run Ekara code
var InitCodeStepResult = initResult(STEP_CONTEXT_CODE)

// InitDescriptorStepResult initializes a step result to process the environment descriptor
var InitDescriptorStepResult = initResult(STEP_CONTEXT_DESCRIPTOR)

// InitParameterStepResult initializes a step results to process the parameter file
var InitParameterStepResult = initResult(STEP_CONTEXT_PARAMETER_FILE)

// InitPlaybookStepResult initializes a step results to execute an ansible playbook
var InitPlaybookStepResult = initResult(STEP_CONTEXT_PLABOOK)

// InitHookStepResult initializes a step results to execute a hook through an ansible playbook
var InitHookStepResult = initResult(STEP_CONTEXT_HOOK_PLABOOK)

// initResult initializes a step results for the given context
func initResult(o stepResultContext) func(stepName string, appliedTo model.Describable, c Cleanup) StepResult {

	return func(stepName string, appliedTo model.Describable, c Cleanup) StepResult {
		result := StepResult{Context: o}
		result.startedAt = time.Now()
		result.StepName = stepName
		if appliedTo != nil {
			result.AppliedToType = appliedTo.DescType()
			result.AppliedToName = appliedTo.DescName()
		}
		result.cleanUp = c
		result.Status = STEP_STATUS_SUCCESS
		return result
	}
}

// Array() initializes an array with the step result
func (sr StepResult) Array() StepResults {
	sr.ExecutionTime = time.Since(sr.startedAt)
	i := int64(sr.ExecutionTime / time.Millisecond)
	if i == 0 {
		sr.ExecutionTime, _ = time.ParseDuration("1ms")
	}
	r := StepResults{
		Results:            []StepResult{sr},
		TotalExecutionTime: sr.ExecutionTime,
	}

	return r

}

// Add adds the given stepResult to the results
func (sr *StepResults) Add(c StepResult) {
	c.ExecutionTime = time.Since(c.startedAt)
	i := int64(c.ExecutionTime / time.Millisecond)
	if i == 0 {
		c.ExecutionTime, _ = time.ParseDuration("1ms")
	}
	sr.Results = append(sr.Results, c)
	sr.TotalExecutionTime = sr.TotalExecutionTime + c.ExecutionTime
}

// InitStepResults initializes an empty stepResults
func InitStepResults() *StepResults {
	sRs := &StepResults{}
	sRs.Results = make([]StepResult, 0)
	return sRs
}

// Content returns the json representation of the report steps
func (sr *StepResults) MarshalJSON() (b []byte, e error) {
	temp := struct {
		Results            []StepResult
		TotalExecutionTime string
	}{
		Results:            sr.Results,
		TotalExecutionTime: fmtDuration(sr.TotalExecutionTime),
	}
	b, e = json.MarshalIndent(&temp, "", "    ")
	return
}

func (sr *StepResult) MarshalJSON() (b []byte, e error) {
	temp := struct {
		StepName        string
		AppliedToType   string `json:",omitempty"`
		AppliedToName   string `json:",omitempty"`
		Status          stepResultStatus
		Context         stepResultContext
		FailureCause    failureCause `json:",omitempty"`
		ErrorMessage    string       `json:",omitempty"`
		ReadableMessage string       `json:",omitempty"`
		RawContent      interface{}  `json:",omitempty"`
		ExecutionTime   string
	}{
		StepName:        sr.StepName,
		AppliedToType:   sr.AppliedToType,
		AppliedToName:   sr.AppliedToName,
		Status:          sr.Status,
		Context:         sr.Context,
		FailureCause:    sr.FailureCause,
		ErrorMessage:    sr.ErrorMessage,
		ReadableMessage: sr.ReadableMessage,
		RawContent:      sr.RawContent,
		ExecutionTime:   fmtDuration(sr.ExecutionTime),
	}
	b, e = json.MarshalIndent(&temp, "", "    ")
	return
}

func (srnt StepResultNoTime) ToStepResult() StepResult {
	return StepResult{
		StepName:        srnt.StepName,
		AppliedToType:   srnt.AppliedToType,
		AppliedToName:   srnt.AppliedToName,
		Status:          srnt.Status,
		Context:         srnt.Context,
		FailureCause:    srnt.FailureCause,
		ErrorMessage:    srnt.ErrorMessage,
		ReadableMessage: srnt.ReadableMessage,
		RawContent:      srnt.RawContent,
	}
}

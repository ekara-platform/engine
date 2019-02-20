package engine

import (
	"github.com/ekara-platform/model"
)

type (
	// Type used to identify final status of the step
	stepResultStatus string
	// Type used to identify the step context, or if you prefer what the step was doing
	stepResultContext string

	// stepResults represents a chain of steps execution results
	StepResults struct {
		Results []StepResult
	}

	// stepResult represents the execution result of a single step with its context
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
	// The step deals with the parameter file used to fill descriptor variables
	STEP_CONTEXT_PARAMETER_FILE stepResultContext = "Environment parameter file"
	// The step purpose is to launch an ansible playbook
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
func (sc StepResult) Array() StepResults {
	return StepResults{
		Results: []StepResult{sc},
	}
}

// Add adds the given stepResult to the results
func (sc *StepResults) Add(c StepResult) {
	sc.Results = append(sc.Results, c)
}

// InitStepResults initializes an empty stepResults
func InitStepResults() *StepResults {
	sRs := &StepResults{}
	sRs.Results = make([]StepResult, 0)
	return sRs
}

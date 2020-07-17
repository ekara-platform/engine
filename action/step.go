package action

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ekara-platform/engine/model"
)

type (
	//stepResultStatus is the type used to identify final status of the step
	stepStatus string
	//stepResultContext is the type used to identify the step context,
	// or if you prefer what the step was doing
	stepInfo string

	//StepResults represents a chain of steps execution results
	//It is used to marchal step results
	StepResults struct {
		Status             []StepResult
		TotalExecutionTime time.Duration
	}

	//StepResult represents an step execution results
	StepResult struct {
		StepName        string
		AppliedToType   string `json:",omitempty"`
		AppliedToName   string `json:",omitempty"`
		Status          stepStatus
		Context         stepInfo
		FailureCause    failureCause `json:",omitempty"`
		ErrorMessage    string       `json:",omitempty"`
		ReadableMessage string       `json:",omitempty"`
		RawContent      interface{}  `json:",omitempty"`
		ExecutionTime   time.Duration
		error           error
		cleanUp         Cleanup
		startedAt       time.Time
	}

	// step represents a sinlge step used to compose a process executed by the installer
	step func(rC *runtimeContext) StepResults
)

const (

	//stepStatusFailure Step execution successful
	stepStatusFailure stepStatus = "Failure"

	//stepStatusSuccess Step execution failed
	stepStatusSuccess stepStatus = "Success"

	//stepContextCode The step belongs internally to Ekara
	stepContextCode stepInfo = "Ekara execution"

	//stepContextDescriptor The step deals with the environment descriptor
	stepContextDescriptor stepInfo = "Environment descriptor content"

	//stepContextParameterFIle The step deals with the parameter file used to fill descriptor variables
	stepContextParameterFIle stepInfo = "Environment parameter file"

	//stepContextPlaybook The step purpose is to launch an ansible playbook
	stepContextPlaybook stepInfo = "Playbook execution"

	//stepContextHookPlaybook The step purpose is to launch a hook through an ansible playbook
	stepContextHookPlaybook stepInfo = "Hook execution"
)

// InitCodeStepResult initializes a step result to run Ekara code
var InitCodeStepResult = initResult(stepContextCode)

// InitDescriptorStepResult initializes a step result to process the environment descriptor
var InitDescriptorStepResult = initResult(stepContextDescriptor)

// InitParameterStepResult initializes a step results to process the parameter file
var InitParameterStepResult = initResult(stepContextParameterFIle)

// InitPlaybookStepResult initializes a step results to execute an ansible playbook
var InitPlaybookStepResult = initResult(stepContextPlaybook)

// InitHookStepResult initializes a step results to execute a hook through an ansible playbook
var InitHookStepResult = initResult(stepContextHookPlaybook)

// initResult initializes a step results for the given context
func initResult(o stepInfo) func(stepName string, appliedTo model.Describable, c Cleanup) StepResult {
	return func(stepName string, appliedTo model.Describable, c Cleanup) StepResult {
		result := StepResult{Context: o}
		result.startedAt = time.Now()
		result.StepName = stepName
		if appliedTo != nil {
			result.AppliedToType = appliedTo.DescType()
			result.AppliedToName = appliedTo.DescName()
		}
		result.cleanUp = c
		result.Status = stepStatusSuccess
		return result
	}
}

// Build initializes an array with the step result
func (sr StepResult) Build() StepResults {
	sr.ExecutionTime = time.Since(sr.startedAt)
	i := int64(sr.ExecutionTime / time.Millisecond)
	if i == 0 {
		sr.ExecutionTime, _ = time.ParseDuration("1ms")
	}
	r := StepResults{
		Status:             []StepResult{sr},
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
	sr.Status = append(sr.Status, c)
	sr.TotalExecutionTime = sr.TotalExecutionTime + c.ExecutionTime
}

// InitStepResults initializes an empty stepResults
func InitStepResults() *StepResults {
	sRs := &StepResults{}
	sRs.Status = make([]StepResult, 0)
	return sRs
}

// MarshalJSON returns the json representation of the report steps
func (sr *StepResults) MarshalJSON() (b []byte, e error) {
	temp := struct {
		Results            []StepResult
		TotalExecutionTime string
	}{
		Results:            sr.Status,
		TotalExecutionTime: fmtDuration(sr.TotalExecutionTime),
	}
	b, e = json.MarshalIndent(&temp, "", "    ")
	return
}

// MarshalJSON returns the json representation of a report step
func (sr *StepResult) MarshalJSON() (b []byte, e error) {
	temp := struct {
		StepName        string
		AppliedToType   string `json:",omitempty"`
		AppliedToName   string `json:",omitempty"`
		Status          stepStatus
		Context         stepInfo
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

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Millisecond)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	ms := d / time.Millisecond
	return fmt.Sprintf("%02dh%02dm%02ds%03d", h, m, s, ms)
}

package engine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockDescribable struct {
}

func (d MockDescribable) DescType() string {
	return "type"
}

func (d MockDescribable) DescName() string {
	return "name"
}

func TestStepAddedExecutionTime(t *testing.T) {

	results := InitStepResults()
	sc := InitCodeStepResult("FirstStep", MockDescribable{}, NoCleanUpRequired)

	assert.Equal(t, sc.ExecutionTime, time.Duration(0))
	assert.Equal(t, results.TotalExecutionTime, time.Duration(0))

	time.Sleep(10 * time.Millisecond)
	results.Add(sc)

	execTime := results.Results[0].ExecutionTime

	//Adding the step into the result should define its execution time
	assert.NotEqual(t, results.Results[0].ExecutionTime, time.Duration(0))

	// It also should define the total of the execution time for the whole result
	assert.NotEqual(t, results.TotalExecutionTime, time.Duration(0))
	assert.Equal(t, results.TotalExecutionTime, execTime)

	report1 := ExecutionReport{
		Steps: *results,
	}

	report2 := ExecutionReport{}

	report2.aggregate(report1)

	// Merging to report should update the total of the execution time for the whole report
	assert.NotEqual(t, report2.Steps.TotalExecutionTime, time.Duration(0))
	assert.Equal(t, report2.Steps.TotalExecutionTime, execTime)

}

func TestStepArrayExecutionTime(t *testing.T) {

	sc := InitCodeStepResult("FirstStep", MockDescribable{}, NoCleanUpRequired)
	assert.Equal(t, sc.ExecutionTime, time.Duration(0))

	time.Sleep(10 * time.Millisecond)
	//Calling Array() on the step should define its execution time
	results := sc.Array()
	execTime := results.Results[0].ExecutionTime
	assert.NotEqual(t, results.Results[0].ExecutionTime, time.Duration(0))

	// It also should define the total of the execution time for the whole result
	assert.NotEqual(t, results.TotalExecutionTime, time.Duration(0))
	assert.Equal(t, results.TotalExecutionTime, execTime)

	report1 := ExecutionReport{
		Steps: results,
	}

	report2 := ExecutionReport{}

	report2.aggregate(report1)

	// Merging to report should update the total of the execution time for the whole report
	assert.NotEqual(t, report2.Steps.TotalExecutionTime, time.Duration(0))
	assert.Equal(t, report2.Steps.TotalExecutionTime, execTime)

}

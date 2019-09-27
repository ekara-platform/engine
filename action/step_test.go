package action

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

func TestTimeFormat(t *testing.T) {
	assert.Equal(t, fmtDuration(1*time.Millisecond), "00h00m00s001")
	assert.Equal(t, fmtDuration(100*time.Millisecond), "00h00m00s100")
	assert.Equal(t, fmtDuration(1000*time.Millisecond), "00h00m01s000")
	assert.Equal(t, fmtDuration(60000*time.Millisecond), "00h01m00s000")
	assert.Equal(t, fmtDuration(3600000*time.Millisecond), "01h00m00s000")
	assert.Equal(t, fmtDuration(3600001*time.Millisecond), "01h00m00s001")
	assert.Equal(t, fmtDuration(3600100*time.Millisecond), "01h00m00s100")
	assert.Equal(t, fmtDuration(3601001*time.Millisecond), "01h00m01s001")
	assert.Equal(t, fmtDuration(3661001*time.Millisecond), "01h01m01s001")

}

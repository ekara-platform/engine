package engine

import (
	"fmt"
	"testing"

	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"
)

func TestApplyTo(t *testing.T) {

	p := model.Provider{Name: "toName"}

	sc := InitCodeStepResult("stepName", p, NoCleanUpRequired)
	assert.NotNil(t, sc)
	assert.Equal(t, sc.AppliedToType, "Provider")
	assert.Equal(t, sc.AppliedToName, "toName")

}

func TestInitCodeStep(t *testing.T) {

	sc := InitCodeStepResult("stepName", nil, NoCleanUpRequired)
	assert.NotNil(t, sc)
	assert.Equal(t, sc.StepName, "stepName")
	assert.Equal(t, sc.AppliedToName, "")
	assert.Equal(t, sc.AppliedToType, "")
	assert.Equal(t, sc.Context, STEP_CONTEXT_CODE)
	assert.NotNil(t, sc.cleanUp)
	checkNoError(t, sc)

	a := sc.Array()
	assert.Equal(t, len(a.Results), 1)

}

func TestInitDescriptorStep(t *testing.T) {

	sc := InitDescriptorStepResult("stepName", nil, NoCleanUpRequired)
	assert.NotNil(t, sc)
	assert.Equal(t, sc.StepName, "stepName")
	assert.Equal(t, sc.AppliedToType, "")
	assert.Equal(t, sc.AppliedToName, "")
	assert.Equal(t, sc.Context, STEP_CONTEXT_DESCRIPTOR)
	assert.NotNil(t, sc.cleanUp)
	checkNoError(t, sc)

	a := sc.Array()
	assert.Equal(t, len(a.Results), 1)

}

func TestInitParamStep(t *testing.T) {

	sc := InitParameterStepResult("stepName", nil, NoCleanUpRequired)
	assert.NotNil(t, sc)
	assert.Equal(t, sc.StepName, "stepName")
	assert.Equal(t, sc.AppliedToType, "")
	assert.Equal(t, sc.AppliedToName, "")
	assert.Equal(t, sc.Context, STEP_CONTEXT_PARAMETER_FILE)
	assert.NotNil(t, sc.cleanUp)
	checkNoError(t, sc)

	a := sc.Array()
	assert.Equal(t, len(a.Results), 1)

}

func TestInitPlaybookStep(t *testing.T) {

	sc := InitPlaybookStepResult("stepName", nil, NoCleanUpRequired)
	assert.NotNil(t, sc)
	assert.Equal(t, sc.StepName, "stepName")
	assert.Equal(t, sc.AppliedToType, "")
	assert.Equal(t, sc.AppliedToName, "")
	assert.Equal(t, sc.Context, STEP_CONTEXT_PLABOOK)
	assert.NotNil(t, sc.cleanUp)
	checkNoError(t, sc)

	a := sc.Array()
	assert.Equal(t, len(a.Results), 1)

}

func checkNoError(t *testing.T, sc StepResult) {
	assert.Equal(t, sc.Status, STEP_STATUS_SUCCESS)
	assert.Equal(t, string(sc.FailureCause), "")
	assert.Nil(t, sc.error)
	assert.Equal(t, sc.ErrorMessage, "")
	assert.Equal(t, sc.ReadableMessage, "")
	assert.Nil(t, sc.RawContent)

}

func TestInitResults(t *testing.T) {

	s := InitStepResults()
	assert.Equal(t, len(s.Results), 0)
	s.Add(InitCodeStepResult("stepName1", nil, NoCleanUpRequired))
	assert.Equal(t, len(s.Results), 1)
	s.Add(InitCodeStepResult("stepName2", nil, NoCleanUpRequired))
	assert.Equal(t, len(s.Results), 2)

}

func TestLaunchSteps(t *testing.T) {
	calls := []step{
		fStepMock1,
		fStepMock2,
		fStepMock3,
	}
	rep := launch(calls, &MockLaunchContext{}, &runtimeContext{})
	assert.NotNil(t, rep)
	assert.Nil(t, rep.Error)
	srs := rep.Steps.Results
	assert.Equal(t, len(srs), 3)

	// Check the order of the executed steps
	assert.Equal(t, srs[0].StepName, "Dummy step 1")
	assert.Equal(t, srs[1].StepName, "Dummy step 2")
	assert.Equal(t, srs[2].StepName, "Dummy step 3")
}

func TestLaunchStepsError(t *testing.T) {
	calls := []step{
		fStepMock1,
		fStepMock2,
		fStepMock3,
		fStepMockError,
	}
	rep := launch(calls, &MockLaunchContext{}, &runtimeContext{})
	assert.NotNil(t, rep)
	assert.NotNil(t, rep.Error)
	srs := rep.Steps.Results
	assert.Equal(t, len(srs), 4)

	// Check the order of the executed steps
	assert.Equal(t, srs[0].StepName, "Dummy step 1")
	assert.Equal(t, srs[1].StepName, "Dummy step 2")
	assert.Equal(t, srs[2].StepName, "Dummy step 3")
	assert.Equal(t, srs[3].StepName, "Dummy step on error")

}

func TestLaunchStepsError2(t *testing.T) {
	calls := []step{
		fStepMock1,
		fStepMock2,
		fStepMockError,
		fStepMock3,
	}
	rep := launch(calls, &MockLaunchContext{}, &runtimeContext{})
	assert.NotNil(t, rep)
	assert.NotNil(t, rep.Error)
	// Because fStepMockError throws an error fStepMock3 is not invoked and
	// then it is never returned into the report
	srs := rep.Steps.Results
	assert.Equal(t, len(srs), 3)

	// Check the order of the executed steps
	assert.Equal(t, srs[0].StepName, "Dummy step 1")
	assert.Equal(t, srs[1].StepName, "Dummy step 2")
	assert.Equal(t, srs[2].StepName, "Dummy step on error")

}

func TestLaunchStepsMultiples(t *testing.T) {
	calls := []step{
		fStepMock1,
		fStepMock2,
		fStepMock3,
		fStepMockMultipleContext,
	}
	rep := launch(calls, &MockLaunchContext{}, &runtimeContext{})
	assert.NotNil(t, rep)
	assert.Nil(t, rep.Error)
	scs := rep.Steps.Results
	assert.Equal(t, len(scs), 6)
	// Check the order of the executed steps
	assert.Equal(t, scs[0].StepName, "Dummy step 1")
	assert.Equal(t, scs[1].StepName, "Dummy step 2")
	assert.Equal(t, scs[2].StepName, "Dummy step 3")
	assert.Equal(t, scs[3].StepName, "Dummy step, multiple 1")
	assert.Equal(t, scs[4].StepName, "Dummy step, multiple 2")
	assert.Equal(t, scs[5].StepName, "Dummy step, multiple 3")
}

func fStepMock1(lC LaunchContext, rC *runtimeContext) StepResults {
	sc := InitCodeStepResult("Dummy step 1", nil, NoCleanUpRequired)
	return sc.Array()
}

func fStepMock2(lC LaunchContext, rC *runtimeContext) StepResults {
	sc := InitCodeStepResult("Dummy step 2", nil, NoCleanUpRequired)
	return sc.Array()
}

func fStepMock3(lC LaunchContext, rC *runtimeContext) StepResults {
	sc := InitCodeStepResult("Dummy step 3", nil, NoCleanUpRequired)
	return sc.Array()
}

func fStepMockError(lC LaunchContext, rC *runtimeContext) StepResults {
	sc := InitCodeStepResult("Dummy step on error", nil, NoCleanUpRequired)
	sc.error = fmt.Errorf("Dummy error")
	return sc.Array()
}

func fStepMockMultipleContext(lC LaunchContext, rC *runtimeContext) StepResults {
	srs := InitStepResults()
	srs.Add(InitCodeStepResult("Dummy step, multiple 1", nil, NoCleanUpRequired))
	srs.Add(InitCodeStepResult("Dummy step, multiple 2", nil, NoCleanUpRequired))
	srs.Add(InitCodeStepResult("Dummy step, multiple 3", nil, NoCleanUpRequired))
	return *srs
}

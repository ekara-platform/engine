package action

import (
	"fmt"
	"testing"

	"github.com/ekara-platform/engine/model"
	"github.com/ekara-platform/engine/util"
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
	assert.Equal(t, sc.Context, stepContextCode)
	assert.NotNil(t, sc.cleanUp)
	checkNoError(t, sc)

	a := sc.Build()
	assert.Equal(t, len(a.Status), 1)

}

func TestInitDescriptorStep(t *testing.T) {

	sc := InitDescriptorStepResult("stepName", nil, NoCleanUpRequired)
	assert.NotNil(t, sc)
	assert.Equal(t, sc.StepName, "stepName")
	assert.Equal(t, sc.AppliedToType, "")
	assert.Equal(t, sc.AppliedToName, "")
	assert.Equal(t, sc.Context, stepContextDescriptor)
	assert.NotNil(t, sc.cleanUp)
	checkNoError(t, sc)

	a := sc.Build()
	assert.Equal(t, len(a.Status), 1)

}

func TestInitParamStep(t *testing.T) {

	sc := InitParameterStepResult("stepName", nil, NoCleanUpRequired)
	assert.NotNil(t, sc)
	assert.Equal(t, sc.StepName, "stepName")
	assert.Equal(t, sc.AppliedToType, "")
	assert.Equal(t, sc.AppliedToName, "")
	assert.Equal(t, sc.Context, stepContextParameterFIle)
	assert.NotNil(t, sc.cleanUp)
	checkNoError(t, sc)

	a := sc.Build()
	assert.Equal(t, len(a.Status), 1)

}

func TestInitPlaybookStep(t *testing.T) {

	sc := InitPlaybookStepResult("stepName", nil, NoCleanUpRequired)
	assert.NotNil(t, sc)
	assert.Equal(t, sc.StepName, "stepName")
	assert.Equal(t, sc.AppliedToType, "")
	assert.Equal(t, sc.AppliedToName, "")
	assert.Equal(t, sc.Context, stepContextPlaybook)
	assert.NotNil(t, sc.cleanUp)
	checkNoError(t, sc)

	a := sc.Build()
	assert.Equal(t, len(a.Status), 1)

}

func checkNoError(t *testing.T, sc StepResult) {
	assert.Equal(t, sc.Status, stepStatusSuccess)
	assert.Equal(t, string(sc.FailureCause), "")
	assert.Nil(t, sc.error)
	assert.Equal(t, sc.ErrorMessage, "")
	assert.Equal(t, sc.ReadableMessage, "")
	assert.Nil(t, sc.RawContent)

}

func TestInitResults(t *testing.T) {

	s := InitStepResults()
	assert.Equal(t, len(s.Status), 0)
	s.Add(InitCodeStepResult("stepName1", nil, NoCleanUpRequired))
	assert.Equal(t, len(s.Status), 1)
	s.Add(InitCodeStepResult("stepName2", nil, NoCleanUpRequired))
	assert.Equal(t, len(s.Status), 2)

}

func TestLaunchSteps(t *testing.T) {
	action := Action{
		steps: []step{
			fStepMock1,
			fStepMock2,
			fStepMock3,
		}}
	rep, _ := action.Execute(mockRuntimeContext())
	assert.NotNil(t, rep)
	assert.Nil(t, rep.Error)
	srs := rep.Steps.Status
	assert.Equal(t, 3, len(srs))

	// Check the order of the executed steps
	assert.Equal(t, srs[0].StepName, "Dummy step 1")
	assert.Equal(t, srs[1].StepName, "Dummy step 2")
	assert.Equal(t, srs[2].StepName, "Dummy step 3")
}

func TestLaunchStepsError(t *testing.T) {
	action := Action{
		steps: []step{
			fStepMock1,
			fStepMock2,
			fStepMock3,
			fStepMockError,
		}}
	rep, _ := action.Execute(mockRuntimeContext())
	assert.NotNil(t, rep)
	assert.NotNil(t, rep.Error)
	srs := rep.Steps.Status
	assert.Equal(t, len(srs), 4)

	// Check the order of the executed steps
	assert.Equal(t, srs[0].StepName, "Dummy step 1")
	assert.Equal(t, srs[1].StepName, "Dummy step 2")
	assert.Equal(t, srs[2].StepName, "Dummy step 3")
	assert.Equal(t, srs[3].StepName, "Dummy step on error")

}

func TestLaunchStepsError2(t *testing.T) {
	action := Action{
		steps: []step{
			fStepMock1,
			fStepMock2,
			fStepMockError,
			fStepMock3,
		}}
	rep, _ := action.Execute(mockRuntimeContext())
	assert.NotNil(t, rep)
	assert.NotNil(t, rep.Error)
	// Because fStepMockError throws an error fStepMock3 is not invoked and
	// then it is never returned into the report
	srs := rep.Steps.Status
	assert.Equal(t, len(srs), 3)

	// Check the order of the executed steps
	assert.Equal(t, srs[0].StepName, "Dummy step 1")
	assert.Equal(t, srs[1].StepName, "Dummy step 2")
	assert.Equal(t, srs[2].StepName, "Dummy step on error")

}

func TestLaunchStepsMultiples(t *testing.T) {
	action := Action{
		steps: []step{
			fStepMock1,
			fStepMock2,
			fStepMock3,
			fStepMockMultipleContext,
		}}
	rep, _ := action.Execute(mockRuntimeContext())
	assert.NotNil(t, rep)
	assert.Nil(t, rep.Error)
	scs := rep.Steps.Status
	assert.Equal(t, len(scs), 6)
	// Check the order of the executed steps
	assert.Equal(t, scs[0].StepName, "Dummy step 1")
	assert.Equal(t, scs[1].StepName, "Dummy step 2")
	assert.Equal(t, scs[2].StepName, "Dummy step 3")
	assert.Equal(t, scs[3].StepName, "Dummy step, multiple 1")
	assert.Equal(t, scs[4].StepName, "Dummy step, multiple 2")
	assert.Equal(t, scs[5].StepName, "Dummy step, multiple 3")
}

func fStepMock1(rC *RuntimeContext) StepResults {
	sc := InitCodeStepResult("Dummy step 1", nil, NoCleanUpRequired)
	return sc.Build()
}

func fStepMock2(rC *RuntimeContext) StepResults {
	sc := InitCodeStepResult("Dummy step 2", nil, NoCleanUpRequired)
	return sc.Build()
}

func fStepMock3(rC *RuntimeContext) StepResults {
	sc := InitCodeStepResult("Dummy step 3", nil, NoCleanUpRequired)
	return sc.Build()
}

func fStepMockError(rC *RuntimeContext) StepResults {
	sc := InitCodeStepResult("Dummy step on error", nil, NoCleanUpRequired)
	sc.error = fmt.Errorf("Dummy error")
	return sc.Build()
}

func fStepMockMultipleContext(rC *RuntimeContext) StepResults {
	srs := InitStepResults()
	srs.Add(InitCodeStepResult("Dummy step, multiple 1", nil, NoCleanUpRequired))
	srs.Add(InitCodeStepResult("Dummy step, multiple 2", nil, NoCleanUpRequired))
	srs.Add(InitCodeStepResult("Dummy step, multiple 3", nil, NoCleanUpRequired))
	return *srs
}

func mockRuntimeContext() *RuntimeContext {
	return CreateRuntimeContext(util.CreateMockLaunchContext(false), nil, nil, model.Environment{}, &model.TemplateContext{})
}

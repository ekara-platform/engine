package action

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeFail(t *testing.T) {

	sc := InitCodeStepResult("DUMMY_STEP", nil, NoCleanUpRequired)
	FailsOnCode(&sc, fmt.Errorf("DUMMY_ERROR"), "DUMMY_DETAILS", nil)

	assert.Equal(t, sc.StepName, "DUMMY_STEP")
	assert.Equal(t, sc.AppliedToType, "")
	assert.Equal(t, sc.AppliedToName, "")
	assert.Equal(t, sc.Status, stepStatusFailure)
	assert.Equal(t, sc.Context, stepContextCode)
	assert.Equal(t, sc.FailureCause, codeFailure)
	assert.NotNil(t, sc.error)
	assert.Equal(t, sc.ErrorMessage, "DUMMY_ERROR")
	assert.Equal(t, sc.ReadableMessage, "DUMMY_DETAILS")
	assert.Nil(t, sc.RawContent)

}

func TestDescriptorFail(t *testing.T) {

	sc := InitCodeStepResult("DUMMY_STEP", nil, NoCleanUpRequired)
	FailsOnModel(&sc, fmt.Errorf("DUMMY_ERROR"), "DUMMY_DETAILS", nil)

	assert.Equal(t, sc.StepName, "DUMMY_STEP")
	assert.Equal(t, sc.AppliedToType, "")
	assert.Equal(t, sc.AppliedToName, "")
	assert.Equal(t, sc.Status, stepStatusFailure)
	assert.Equal(t, sc.Context, stepContextCode)
	assert.Equal(t, sc.FailureCause, modelFailure)
	assert.NotNil(t, sc.error)
	assert.Equal(t, sc.ErrorMessage, "DUMMY_ERROR")
	assert.Equal(t, sc.ReadableMessage, "DUMMY_DETAILS")
	assert.Nil(t, sc.RawContent)

}

func TestPlaybookFail(t *testing.T) {

	sc := InitCodeStepResult("DUMMY_STEP", nil, NoCleanUpRequired)
	FailsOnPlaybook(&sc, fmt.Errorf("DUMMY_ERROR"), "DUMMY_DETAILS", nil)

	assert.Equal(t, sc.StepName, "DUMMY_STEP")
	assert.Equal(t, sc.AppliedToType, "")
	assert.Equal(t, sc.AppliedToName, "")
	assert.Equal(t, sc.Status, stepStatusFailure)
	assert.Equal(t, sc.Context, stepContextCode)
	assert.Equal(t, sc.FailureCause, playBookFailure)
	assert.NotNil(t, sc.error)
	assert.Equal(t, sc.ErrorMessage, "DUMMY_ERROR")
	assert.Equal(t, sc.ReadableMessage, "DUMMY_DETAILS")
	assert.Nil(t, sc.RawContent)

}

func TestNotImplentedFail(t *testing.T) {

	sc := InitCodeStepResult("DUMMY_STEP", nil, NoCleanUpRequired)
	FailsOnNotImplemented(&sc, fmt.Errorf("DUMMY_ERROR"), "DUMMY_DETAILS", nil)

	assert.Equal(t, sc.StepName, "DUMMY_STEP")
	assert.Equal(t, sc.AppliedToType, "")
	assert.Equal(t, sc.AppliedToName, "")
	assert.Equal(t, sc.Status, stepStatusFailure)
	assert.Equal(t, sc.Context, stepContextCode)
	assert.Equal(t, sc.FailureCause, notImplementedFailure)
	assert.NotNil(t, sc.error)
	assert.Equal(t, sc.ErrorMessage, "DUMMY_ERROR")
	assert.Equal(t, sc.ReadableMessage, "DUMMY_DETAILS")
	assert.Nil(t, sc.RawContent)

}

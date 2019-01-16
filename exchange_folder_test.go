package engine

import (
	"log"
	"os"
	"testing"

	"github.com/ekara-platform/engine/util"
	"github.com/stretchr/testify/assert"
)

func TestChildExchangeFolderOk(t *testing.T) {

	ef, e := util.CreateExchangeFolder("./", "testFolder")
	assert.Nil(t, e)
	assert.NotNil(t, ef)
	defer ef.Delete()

	e = ef.Create()
	assert.Nil(t, e)

	log := log.New(os.Stdout, "Test", log.Ldate|log.Ltime|log.Lmicroseconds)
	sc := InitCodeStepResult("DummyStep", nil, NoCleanUpRequired)

	subEf, ko := createChildExchangeFolder(ef.Input, "subTestFolder", &sc, log)
	assert.False(t, ko)
	assert.NotNil(t, subEf)

	assert.Equal(t, sc.StepName, "DummyStep")
	assert.Equal(t, sc.AppliedToType, "")
	assert.Equal(t, sc.AppliedToName, "")
	assert.Equal(t, sc.Status, STEP_STATUS_SUCCESS)
	assert.Equal(t, sc.Context, STEP_CONTEXT_CODE)
	assert.Equal(t, string(sc.FailureCause), "")
	assert.Nil(t, sc.error)
	assert.Equal(t, sc.ErrorMessage, "")
	assert.Equal(t, sc.ReadableMessage, "")
	assert.Nil(t, sc.RawContent)

	_, err := os.Stat(subEf.Location.Path())
	assert.Nil(t, err)
}

func TestChildExchangeFolderKo(t *testing.T) {

	ef, e := util.CreateExchangeFolder("./", "testFolder")
	assert.Nil(t, e)
	assert.NotNil(t, ef)
	defer ef.Delete()

	// We are not calling ef.Create() in order to get an error creating the child

	log := log.New(os.Stdout, "Test", log.Ldate|log.Ltime|log.Lmicroseconds)
	sc := InitCodeStepResult("DummyStep", nil, NoCleanUpRequired)

	subEf, ko := createChildExchangeFolder(ef.Input, "subTestFolfer", &sc, log)
	assert.True(t, ko)
	assert.NotNil(t, subEf)

	assert.Equal(t, sc.StepName, "DummyStep")
	assert.Equal(t, sc.AppliedToType, "")
	assert.Equal(t, sc.AppliedToName, "")
	assert.Equal(t, sc.Status, STEP_STATUS_FAILURE)
	assert.Equal(t, sc.Context, STEP_CONTEXT_CODE)
	assert.Equal(t, sc.FailureCause, CODE_FAILURE)
	assert.NotNil(t, sc.error)
	assert.Nil(t, sc.RawContent)

}

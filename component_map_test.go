package engine

import (
	"log"
	"os"
	"testing"

	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"
)

func TestSaveComponentMapOk(t *testing.T) {

	ef, e := util.CreateExchangeFolder("./", "testFolder")
	assert.Nil(t, e)
	assert.NotNil(t, ef)
	defer ef.Delete()

	assert.Nil(t, e)

	sc := InitCodeStepResult("DummyStep", nil, NoCleanUpRequired)

	c := MockLaunchContext{
		efolder:              ef,
		logger:               log.New(os.Stdout, "Test", log.Ldate|log.Ltime|log.Lmicroseconds),
		sshPublicKeyContent:  "sshPublicKey_content",
		sshPrivateKeyContent: "sshPrivateKey_content",
		engine: EkaraMock{
			Env: model.Environment{
				Name:      "NameContent",
				Qualifier: "QualifierContent",
			},
		},
	}

	ko := saveComponentMap(&c, ef.Input, &sc)
	assert.False(t, ko)

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

	ok := ef.Input.Contains(util.ComponentPathsFileName)
	assert.True(t, ok)

}

package engine

import (
	"log"
	"os"
	"testing"

	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"
)

func TestSaveBaseParamOk(t *testing.T) {

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
	bp := BuildBaseParam(c, "nodeId", "providerName")

	ko := saveBaseParams(bp, &c, ef.Input, &sc)
	assert.False(t, ko)

	assert.Nil(t, sc.error)
	assert.Equal(t, sc.ErrorMessage, "")
	assert.Equal(t, string(sc.FailureCause), "")
	assert.Equal(t, string(sc.Status), string(STEP_STATUS_SUCCESS))

	ok, _, err := ef.Input.ContainsParamYaml()
	assert.True(t, ok)
	assert.Nil(t, err)
}

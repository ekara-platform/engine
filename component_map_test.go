package engine

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"
)

func TestSaveComponentMapOk(t *testing.T) {

	p, _ := model.CreateParameters(map[string]interface{}{
		"ek": map[interface{}]interface{}{
			"aws": map[interface{}]interface{}{
				"region": "dummy",
				"accessKey": map[interface{}]interface{}{
					"id":     "dummy",
					"secret": "dummy",
				},
			},
		},
	})
	mainPath := "./testdata/gittest/descriptor"
	tc := model.CreateContext(p)

	ef, e := util.CreateExchangeFolder("./", "testFolder")
	assert.Nil(t, e)
	assert.NotNil(t, ef)
	defer ef.Delete()

	c := &MockLaunchContext{
		locationContent:      mainPath,
		templateContext:      tc,
		efolder:              ef,
		logger:               log.New(ioutil.Discard, "Test", log.Ldate|log.Ltime|log.Lmicroseconds),
		sshPublicKeyContent:  "sshPublicKey_content",
		sshPrivateKeyContent: "sshPrivateKey_content",
	}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDesc := tester.createRep(mainPath)

	descContent := `
name: NameContent
qualifier: QualifierContent

`
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	sc := InitCodeStepResult("DummyStep", nil, NoCleanUpRequired)

	ko := saveComponentMap(c, ef.Input, &sc)
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
	//TODO check the file content here
	assert.True(t, ok)

}

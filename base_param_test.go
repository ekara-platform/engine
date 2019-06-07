package engine

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"
)

func TestSaveBaseParamOk(t *testing.T) {
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
	tester := gitTester(t, c)
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
	bp := BuildBaseParam(c, "nodeId", "providerName")
	ko := saveBaseParams(bp, c, ef.Input, &sc)
	assert.False(t, ko)

	assert.Nil(t, sc.error)
	assert.Equal(t, sc.ErrorMessage, "")
	assert.Equal(t, string(sc.FailureCause), "")
	assert.Equal(t, string(sc.Status), string(STEP_STATUS_SUCCESS))

	ok, _, err := ef.Input.ContainsParamYaml()
	assert.True(t, ok)
	assert.Nil(t, err)

}

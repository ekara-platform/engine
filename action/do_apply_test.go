package action

import (
	"os"
	"testing"

	"github.com/ekara-platform/engine/component"

	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"
)

func TestChildExchangeFolderOk(t *testing.T) {
	ef, e := util.CreateExchangeFolder("./", "testFolder")
	assert.Nil(t, e)
	assert.NotNil(t, ef)
	defer ef.Delete()

	e = ef.Create()
	assert.Nil(t, e)

	sc := InitCodeStepResult("DummyStep", nil, NoCleanUpRequired)

	subEf, ko := createChildExchangeFolder(ef.Input, "subTestFolder", &sc)
	assert.False(t, ko)
	assert.NotNil(t, subEf)

	assert.Equal(t, sc.StepName, "DummyStep")
	assert.Equal(t, sc.AppliedToType, "")
	assert.Equal(t, sc.AppliedToName, "")
	assert.Equal(t, sc.Status, stepStatusSuccess)
	assert.Equal(t, sc.Context, stepContextCode)
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
	sc := InitCodeStepResult("DummyStep", nil, NoCleanUpRequired)

	subEf, ko := createChildExchangeFolder(ef.Input, "subTestFolfer", &sc)
	assert.True(t, ko)
	assert.NotNil(t, subEf)

	assert.Equal(t, sc.StepName, "DummyStep")
	assert.Equal(t, sc.AppliedToType, "")
	assert.Equal(t, sc.AppliedToName, "")
	assert.Equal(t, sc.Status, stepStatusFailure)
	assert.Equal(t, sc.Context, stepContextCode)
	assert.Equal(t, sc.FailureCause, codeFailure)
	assert.NotNil(t, sc.error)
	assert.Nil(t, sc.RawContent)

}

func TestSaveBaseParamOk(t *testing.T) {
	p := model.CreateParameters(map[string]interface{}{
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

	ef, e := util.CreateExchangeFolder("./", "testFolder")
	assert.Nil(t, e)
	assert.NotNil(t, ef)
	defer ef.Delete()

	c := util.CreateMockLaunchContextWithDataAndFolder(mainPath, p, ef, false)
	tester := component.CreateComponentTester(t, c)
	defer tester.Clean()

	repDesc := tester.CreateRep(mainPath)

	descContent := `
name: NameContent
qualifier: QualifierContent

`
	repDesc.WriteCommit("ekara.yaml", descContent)

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)

	sc := InitCodeStepResult("DummyStep", nil, NoCleanUpRequired)
	bp := buildBaseParam(mockRuntimeContextWithParameters(c), "nodeId")
	ko := saveBaseParams(bp, ef.Input, &sc)
	assert.False(t, ko)

	assert.Nil(t, sc.error)
	assert.Equal(t, sc.ErrorMessage, "")
	assert.Equal(t, string(sc.FailureCause), "")
	assert.Equal(t, string(sc.Status), string(stepStatusSuccess))

	ok, _, err := ef.Input.ContainsParamYaml()
	assert.True(t, ok)
	assert.Nil(t, err)

}

func mockRuntimeContextWithParameters(lC util.LaunchContext) *runtimeContext {
	env := model.InitEnvironment()
	return createRuntimeContext(lC, nil, nil, env, &model.TemplateContext{})
}

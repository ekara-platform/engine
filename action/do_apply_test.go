package action

import (
	"os"
	"testing"

	"github.com/ekara-platform/engine/model"
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

	tester := util.CreateComponentTester(t, p)
	defer tester.Clean()

	mainPath := "./testdata/gittest/descriptor"

	ef, e := util.CreateExchangeFolder("./", "testFolder")
	assert.Nil(t, e)
	assert.NotNil(t, ef)
	defer ef.Delete()

	repDesc := tester.CreateDir(mainPath)

	descContent := `
name: NameContent
qualifier: QualifierContent

`
	repDesc.WriteCommit("ekara.yaml", descContent)

	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)

	sc := InitCodeStepResult("DummyStep", nil, NoCleanUpRequired)
	c := util.CreateMockLaunchContextWithDataAndFolder(p, ef, false)
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

func TestStacksDependencies(t *testing.T) {
	sts := model.Stacks{}

	// This test will check the dependencies on the following tree
	//
	//        0                    1
	//                           / | \
	//                          2  7  8
	//                        / |     | \
	//                       3  6     9  12
	//                      / |       | \
	//                     4  5       10  11

	sts["6"] = model.Stack{Name: "6", Dependencies: []string{"2"}}
	sts["1"] = model.Stack{Name: "1"}
	sts["3"] = model.Stack{Name: "3", Dependencies: []string{"2"}}
	sts["4"] = model.Stack{Name: "4", Dependencies: []string{"3"}}
	sts["9"] = model.Stack{Name: "9", Dependencies: []string{"8"}}
	sts["5"] = model.Stack{Name: "5", Dependencies: []string{"3"}}
	sts["7"] = model.Stack{Name: "7", Dependencies: []string{"1"}}
	sts["8"] = model.Stack{Name: "8", Dependencies: []string{"1"}}
	sts["10"] = model.Stack{Name: "10", Dependencies: []string{"9"}}
	sts["2"] = model.Stack{Name: "2", Dependencies: []string{"1"}}
	sts["12"] = model.Stack{Name: "12", Dependencies: []string{"8"}}
	sts["11"] = model.Stack{Name: "11", Dependencies: []string{"9"}}
	sts["0"] = model.Stack{Name: "0"}

	assert.Len(t, sts, 13)
	resolved, err := sortStacks(sts)
	assert.Equal(t, len(resolved), len(sts))
	assert.Nil(t, err)

	processed := make(map[string]model.Stack)
	for _, val := range resolved {
		if len(processed) == 0 {
			processed[val.Name] = val
			continue
		}
		//Check than all the dependencies has been already processd
		for _, d := range val.Dependencies {
			if _, ok := processed[d]; !ok {
				assert.Fail(t, "Dependency %s has not been yet processed", d)
			}

		}
		processed[val.Name] = val
	}
	//Check that the original Stacks has been untouched
	assert.Len(t, sts, 13)

}

func TestStacksMultiplesDependencies(t *testing.T) {
	sts := model.Stacks{}

	// This test will check the dependencies on the following tree
	//
	//        0                    1
	//                           / | \
	//                          2  7  8
	//                        / | /    | \
	//                       3  6     9  12
	//                      / |   \   | \
	//                     4  5    \-10  11

	sts["6"] = model.Stack{Name: "6", Dependencies: []string{"2", "7"}}
	sts["1"] = model.Stack{Name: "1"}
	sts["3"] = model.Stack{Name: "3", Dependencies: []string{"2"}}
	sts["4"] = model.Stack{Name: "4", Dependencies: []string{"3"}}
	sts["9"] = model.Stack{Name: "9", Dependencies: []string{"8"}}
	sts["5"] = model.Stack{Name: "5", Dependencies: []string{"3"}}
	sts["7"] = model.Stack{Name: "7", Dependencies: []string{"1"}}
	sts["8"] = model.Stack{Name: "8", Dependencies: []string{"1"}}
	sts["10"] = model.Stack{Name: "10", Dependencies: []string{"9", "6"}}
	sts["2"] = model.Stack{Name: "2", Dependencies: []string{"1"}}
	sts["12"] = model.Stack{Name: "12", Dependencies: []string{"8"}}
	sts["11"] = model.Stack{Name: "11", Dependencies: []string{"9"}}
	sts["0"] = model.Stack{Name: "0"}

	assert.Len(t, sts, 13)
	resolved, err := sortStacks(sts)
	assert.Equal(t, len(resolved), len(sts))
	assert.Nil(t, err)

	processed := make(map[string]model.Stack)
	for _, val := range resolved {
		if len(processed) == 0 {
			processed[val.Name] = val
			continue
		}
		//Check than all the dependencies has been already processd
		for _, d := range val.Dependencies {
			if _, ok := processed[d]; !ok {
				assert.Fail(t, "Dependency %s has not been yet processed", d)
			}

		}
		processed[val.Name] = val
	}
	//Check that the original Stacks has been untouched
	assert.Len(t, sts, 13)

}

func TestStacksCyclicDependencies(t *testing.T) {
	sts := model.Stacks{}

	sts["1"] = model.Stack{Name: "1", Dependencies: []string{"3"}}
	sts["2"] = model.Stack{Name: "2", Dependencies: []string{"1"}}
	sts["3"] = model.Stack{Name: "3", Dependencies: []string{"2"}}

	assert.Len(t, sts, 3)
	_, err := sortStacks(sts)
	assert.NotNil(t, err)

	//Check that the original Stacks has been untouched
	assert.Len(t, sts, 3)

}

func mockRuntimeContextWithParameters(lC util.LaunchContext) *RuntimeContext {
	return CreateRuntimeContext(lC, nil, nil, model.Environment{}, &model.TemplateContext{})
}

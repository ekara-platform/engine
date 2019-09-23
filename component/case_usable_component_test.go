package component

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

var (
	usableComp1Content = `
templates:
  - "{{ .Vars.templateDef }}"
`

	usableParentContent = `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
`
	usableDescContent = `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent

# Following content just to force the download of comp1 and comp2
orchestrator:
  component: comp1
providers:
  p1:
    component: comp1

nodes:
  node1:
    instances: 1
    provider:
      name: p1
`
)

func TestUsableTemplateOneMatch(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	p, _ := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget1.yaml",
	})

	c := util.CreateMockLaunchContextWithData(mainPath, p, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	writecheckUsableCommon(t, tester, mainPath)

	cptComp := 3
	env, cm := checkUsableCommon(t, tester, cptComp)
	oComp, err := env.Orchestrator.Component()
	assert.Nil(t, err)
	ok, _ := oComp.Templatable()
	if assert.True(t, ok) {
		usableComp, err := cm.Use(env.Orchestrator, tester.cM.TemplateContext())
		assert.Nil(t, err)
		assert.True(t, usableComp.Templated())
		// Check the existence of the templated folder
		assert.Equal(t, cptComp+1, tester.ComponentCount())
		assert.True(t, tester.RootContainsComponent(usableComp.RootPath()))

		// Check  that the file content has been templated
		tester.CheckFile(usableComp, "templateTarget1.yaml", "templateContentFromCli")
		// Check  that the file content has not been templated
		tester.CheckFile(usableComp, "templateTarget2.yaml", "{{ .Vars.templateContent }}")

		// Check the release
		usableComp.Release()
		assert.Equal(t, cptComp, tester.ComponentCount())
		assert.False(t, tester.RootContainsComponent(usableComp.RootPath()))
	}
}

func TestUsableTemplateMatch2Usable(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	p, _ := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget1.yaml",
	})

	c := util.CreateMockLaunchContextWithData(mainPath, p, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	writecheckUsableCommon(t, tester, mainPath)

	cptComp := 3
	env, cm := checkUsableCommon(t, tester, cptComp)

	oComp, err := env.Orchestrator.Component()
	assert.Nil(t, err)
	ok, _ := oComp.Templatable()
	if assert.True(t, ok) {
		usableComp1, err := cm.Use(env.Orchestrator, tester.cM.TemplateContext())
		assert.Nil(t, err)
		assert.True(t, usableComp1.Templated())
		// Check the existence of the templated folder
		assert.Equal(t, cptComp+1, tester.ComponentCount())
		assert.True(t, tester.RootContainsComponent(usableComp1.RootPath()))

		usableComp2, err := cm.Use(env.Orchestrator, tester.cM.TemplateContext())
		assert.Nil(t, err)
		assert.True(t, usableComp2.Templated())
		// Check the existence of a new templated folder
		assert.Equal(t, cptComp+2, tester.ComponentCount())
		assert.True(t, tester.RootContainsComponent(usableComp2.RootPath()))

		// Check the release of usableComp1
		usableComp1.Release()
		assert.Equal(t, cptComp+1, tester.ComponentCount())
		assert.False(t, tester.RootContainsComponent(usableComp1.RootPath()))

		// Check the release of usableComp2
		usableComp2.Release()
		assert.Equal(t, cptComp, tester.ComponentCount())
		assert.False(t, tester.RootContainsComponent(usableComp2.RootPath()))

	}
}

func TestUsableTemplateDoubleMatch(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	p, _ := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget[12].yaml",
	})

	c := util.CreateMockLaunchContextWithData(mainPath, p, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	writecheckUsableCommon(t, tester, mainPath)

	cptComp := 3
	env, cm := checkUsableCommon(t, tester, cptComp)

	oComp, err := env.Orchestrator.Component()
	assert.Nil(t, err)
	ok, _ := oComp.Templatable()
	if assert.True(t, ok) {
		usableComp, err := cm.Use(env.Orchestrator, tester.cM.TemplateContext())
		assert.Nil(t, err)
		assert.True(t, usableComp.Templated())
		// Check the existence of the templated folder
		assert.Equal(t, cptComp+1, tester.ComponentCount())
		assert.True(t, tester.RootContainsComponent(usableComp.RootPath()))

		// Check  that the two file content has been templated
		tester.CheckFile(usableComp, "templateTarget1.yaml", "templateContentFromCli")
		tester.CheckFile(usableComp, "templateTarget2.yaml", "templateContentFromCli")

		// Check the release
		usableComp.Release()
		assert.Equal(t, cptComp, tester.ComponentCount())
		assert.False(t, tester.RootContainsComponent(usableComp.RootPath()))
	}
}

func TestUsableTemplateNoMatch(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	p, _ := model.CreateParameters(map[string]interface{}{
		"templateDef": "/noMatchinTarget.yaml",
	})

	c := util.CreateMockLaunchContextWithData(mainPath, p, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	writecheckUsableCommon(t, tester, mainPath)

	cptComp := 3
	env, cm := checkUsableCommon(t, tester, cptComp)

	oComp, err := env.Orchestrator.Component()
	assert.Nil(t, err)
	ok, _ := oComp.Templatable()
	if assert.True(t, ok) {
		usableComp, err := cm.Use(env.Orchestrator, tester.cM.TemplateContext())
		assert.Nil(t, err)
		assert.False(t, usableComp.Templated())
		// Check that no templated folder has been created
		assert.Equal(t, cptComp, tester.ComponentCount())
		assert.True(t, tester.RootContainsComponent(usableComp.RootPath()))
		// Check  that the files content has not been templated
		tester.CheckFile(usableComp, "templateTarget1.yaml", "{{ .Vars.templateContent }}")
		tester.CheckFile(usableComp, "templateTarget2.yaml", "{{ .Vars.templateContent }}")
		// Check that the release has no effect
		usableComp.Release()
		assert.Equal(t, cptComp, tester.ComponentCount())
		assert.True(t, tester.RootContainsComponent(usableComp.RootPath()))
	}
}

func writecheckUsableCommon(t *testing.T, tester *ComponentTester, d string) {
	repParent := tester.CreateRep("./testdata/gittest/parent")
	repComp1 := tester.CreateRep("./testdata/gittest/comp1")
	repDesc := tester.CreateRep(d)

	repComp1.WriteCommit("ekara.yaml", usableComp1Content)
	repComp1.WriteCommit("templateTarget1.yaml", `{{ .Vars.templateContent }}`)
	repComp1.WriteCommit("templateTarget2.yaml", `{{ .Vars.templateContent }}`)
	repParent.WriteCommit("ekara.yaml", usableParentContent)
	repDesc.WriteCommit("ekara.yaml", usableDescContent)

}

func checkUsableCommon(t *testing.T, tester *ComponentTester, initialComp int) (model.Environment, ComponentManager) {
	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)

	tester.AssertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1")

	cm := tester.cM
	assert.NotNil(t, cm)

	assert.Equal(t, initialComp, tester.ComponentCount())
	return env, cm
}

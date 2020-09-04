package engine

import (
	"github.com/GroupePSA/componentizer"
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/engine/model"

	"github.com/stretchr/testify/assert"
)

var (
	usableComp1Content = `
ekara:
  templates:
    - "{{ .Vars.templateDef }}"
`

	usableParentContent = `
ekara:
  components:
    comp1:
      repository: comp1
`
	usableDescContent = `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: parent

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

	mainPath := "descriptor"

	p := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget1.yaml",
	})

	tester := util.CreateComponentTester(t, p)
	defer tester.Clean()

	repo := writeCheckUsableCommon(tester, mainPath)

	cptComp := 3
	env, cm := checkUsableCommon(tester, repo, cptComp)
	oComp, err := env.Orchestrator.Component(tester.Model())
	assert.Nil(t, err)
	ok, _ := oComp.GetTemplates()
	if assert.True(t, ok) {
		usableComp, err := cm.Use(env.Orchestrator, tester.TemplateContext())
		assert.Nil(t, err)
		assert.True(t, usableComp.Templated())
		// Check the existence of the templated folder
		assert.Equal(t, cptComp, tester.ComponentCount())
		// TODO: // TODO: assert.True(t, tester.RootContainsComponent(usableComp.RootPath()))

		// Check  that the file content has been templated
		tester.AssertFileContent(usableComp, "templateTarget1.yaml", "templateContentFromCli")
		// Check  that the file content has not been templated
		tester.AssertFileContent(usableComp, "templateTarget2.yaml", "{{ .Vars.templateContent }}")

		// Check the release
		usableComp.Release()
		assert.Equal(t, cptComp, tester.ComponentCount())
		// TODO: assert.False(t, tester.RootContainsComponent(usableComp.RootPath()))
	}
}

func TestUsableTemplateMatch2Usable(t *testing.T) {

	mainPath := "descriptor"

	p := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget1.yaml",
	})

	tester := util.CreateComponentTester(t, p)
	defer tester.Clean()

	repo := writeCheckUsableCommon(tester, mainPath)

	cptComp := 3
	env, cm := checkUsableCommon(tester, repo, cptComp)

	oComp, err := env.Orchestrator.Component(tester.Model())
	assert.Nil(t, err)
	ok, _ := oComp.GetTemplates()
	if assert.True(t, ok) {
		usableComp1, err := cm.Use(env.Orchestrator, tester.TemplateContext())
		assert.Nil(t, err)
		assert.True(t, usableComp1.Templated())
		// Check the existence of the templated folder
		assert.Equal(t, cptComp, tester.ComponentCount())
		// TODO: assert.True(t, tester.RootContainsComponent(usableComp1.RootPath()))

		usableComp2, err := cm.Use(env.Orchestrator, tester.TemplateContext())
		assert.Nil(t, err)
		assert.True(t, usableComp2.Templated())
		// Check the existence of a new templated folder
		assert.Equal(t, cptComp, tester.ComponentCount())
		// TODO: assert.True(t, tester.RootContainsComponent(usableComp2.RootPath()))

		// Check the release of usableComp1
		usableComp1.Release()
		assert.Equal(t, cptComp, tester.ComponentCount())
		// TODO: assert.False(t, tester.RootContainsComponent(usableComp1.RootPath()))

		// Check the release of usableComp2
		usableComp2.Release()
		assert.Equal(t, cptComp, tester.ComponentCount())
		// TODO: assert.False(t, tester.RootContainsComponent(usableComp2.RootPath()))

	}
}

func TestUsableTemplateDoubleMatch(t *testing.T) {

	mainPath := "descriptor"

	p := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget[12].yaml",
	})

	tester := util.CreateComponentTester(t, p)
	defer tester.Clean()

	repo := writeCheckUsableCommon(tester, mainPath)

	cptComp := 3
	env, cm := checkUsableCommon(tester, repo, cptComp)

	oComp, err := env.Orchestrator.Component(tester.Model())
	assert.Nil(t, err)
	ok, _ := oComp.GetTemplates()
	if assert.True(t, ok) {
		usableComp, err := cm.Use(env.Orchestrator, tester.TemplateContext())
		assert.Nil(t, err)
		assert.True(t, usableComp.Templated())
		// Check the existence of the templated folder
		assert.Equal(t, cptComp, tester.ComponentCount())
		// TODO: assert.True(t, tester.RootContainsComponent(usableComp.RootPath()))

		// Check  that the two file content has been templated
		tester.AssertFileContent(usableComp, "templateTarget1.yaml", "templateContentFromCli")
		tester.AssertFileContent(usableComp, "templateTarget2.yaml", "templateContentFromCli")

		// Check the release
		usableComp.Release()
		assert.Equal(t, cptComp, tester.ComponentCount())
		// TODO: assert.False(t, tester.RootContainsComponent(usableComp.RootPath()))
	}
}

func TestUsableTemplateNoMatch(t *testing.T) {

	mainPath := "descriptor"

	p := model.CreateParameters(map[string]interface{}{
		"templateDef": "/noMatchinTarget.yaml",
	})

	tester := util.CreateComponentTester(t, p)
	defer tester.Clean()

	repo := writeCheckUsableCommon(tester, mainPath)

	cptComp := 3
	env, cm := checkUsableCommon(tester, repo, cptComp)

	oComp, err := env.Orchestrator.Component(tester.Model())
	assert.Nil(t, err)
	ok, _ := oComp.GetTemplates()
	if assert.True(t, ok) {
		usableComp, err := cm.Use(env.Orchestrator, tester.TemplateContext())
		assert.Nil(t, err)
		assert.False(t, usableComp.Templated())
		// Check that no templated folder has been created
		assert.Equal(t, cptComp, tester.ComponentCount())
		// TODO: assert.True(t, tester.RootContainsComponent(usableComp.RootPath()))
		// Check  that the files content has not been templated
		tester.AssertFileContent(usableComp, "templateTarget1.yaml", "{{ .Vars.templateContent }}")
		tester.AssertFileContent(usableComp, "templateTarget2.yaml", "{{ .Vars.templateContent }}")
		// Check that the release has no effect
		usableComp.Release()
		assert.Equal(t, cptComp, tester.ComponentCount())
		// TODO: assert.True(t, tester.RootContainsComponent(usableComp.RootPath()))
	}
}

func writeCheckUsableCommon(tester util.EkaraComponentTester, d string) componentizer.Repository {
	repParent := tester.CreateDir("parent")
	repComp1 := tester.CreateDir("comp1")
	repDesc := tester.CreateDir(d)

	repComp1.WriteCommit("ekara.yaml", usableComp1Content)
	repComp1.WriteCommit("templateTarget1.yaml", `{{ .Vars.templateContent }}`)
	repComp1.WriteCommit("templateTarget2.yaml", `{{ .Vars.templateContent }}`)
	repParent.WriteCommit("ekara.yaml", usableParentContent)
	repDesc.WriteCommit("ekara.yaml", usableDescContent)
	return repDesc.AsRepository("master")
}

func checkUsableCommon(tester util.EkaraComponentTester, repo componentizer.Repository, initialComp int) (model.Environment, componentizer.ComponentManager) {
	tester.Init(repo)
	env := tester.Env()
	assert.NotNil(tester.T(), env)

	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1")

	assert.Equal(tester.T(), initialComp, tester.ComponentCount())
	return env, tester.ComponentManager()
}

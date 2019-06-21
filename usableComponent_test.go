package engine

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

var (
	usableComp1Content = `
templates:
  - "{{ .Vars.templateDef }}"
`

	usableDistContent = `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
`
	usableDescContent = `
name: ekara-demo-var
qualifier: dev

ekara:
  distribution:
    repository: ./testdata/gittest/distribution

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
	tc := model.CreateContext(p)

	c := &MockLaunchContext{locationContent: mainPath, templateContext: tc}
	tester := gitTester(t, c)
	defer tester.clean()

	writecheckUsableCommon(t, tester, mainPath)

	cptComp := 3
	env, cm := checkUsableCommon(t, c, tester, cptComp)
	oComp, err := env.Orchestrator.Component()
	assert.Nil(t, err)
	ok, _ := oComp.Templatable()
	if assert.True(t, ok) {
		usableComp := cm.Use(env.Orchestrator)
		assert.True(t, usableComp.Templated())
		// Check the existence of the templated folder
		assert.Equal(t, cptComp+1, tester.countComponent())
		assert.True(t, tester.rootContainsComponent(usableComp.RootPath()))

		// Check  that the file content has been templated
		checkFile(t, usableComp, "templateTarget1.yaml", "templateContentFromCli")
		// Check  that the file content has not been templated
		checkFile(t, usableComp, "templateTarget2.yaml", "{{ .Vars.templateContent }}")

		// Check the release
		usableComp.Release()
		assert.Equal(t, cptComp, tester.countComponent())
		assert.False(t, tester.rootContainsComponent(usableComp.RootPath()))
	}
}

func TestUsableTemplateMatch2Usable(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	p, _ := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget1.yaml",
	})
	tc := model.CreateContext(p)

	c := &MockLaunchContext{locationContent: mainPath, templateContext: tc}
	tester := gitTester(t, c)
	defer tester.clean()

	writecheckUsableCommon(t, tester, mainPath)

	cptComp := 3
	env, cm := checkUsableCommon(t, c, tester, cptComp)

	oComp, err := env.Orchestrator.Component()
	assert.Nil(t, err)
	ok, _ := oComp.Templatable()
	if assert.True(t, ok) {
		usableComp1 := cm.Use(env.Orchestrator)
		assert.True(t, usableComp1.Templated())
		// Check the existence of the templated folder
		assert.Equal(t, cptComp+1, tester.countComponent())
		assert.True(t, tester.rootContainsComponent(usableComp1.RootPath()))

		usableComp2 := cm.Use(env.Orchestrator)
		assert.True(t, usableComp2.Templated())
		// Check the existence of a new templated folder
		assert.Equal(t, cptComp+2, tester.countComponent())
		assert.True(t, tester.rootContainsComponent(usableComp2.RootPath()))

		// Check the release of usableComp1
		usableComp1.Release()
		assert.Equal(t, cptComp+1, tester.countComponent())
		assert.False(t, tester.rootContainsComponent(usableComp1.RootPath()))

		// Check the release of usableComp2
		usableComp2.Release()
		assert.Equal(t, cptComp, tester.countComponent())
		assert.False(t, tester.rootContainsComponent(usableComp2.RootPath()))

	}
}

func TestUsableTemplateDoubleMatch(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	p, _ := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget[12].yaml",
	})
	tc := model.CreateContext(p)

	c := &MockLaunchContext{locationContent: mainPath, templateContext: tc}
	tester := gitTester(t, c)
	defer tester.clean()

	writecheckUsableCommon(t, tester, mainPath)

	cptComp := 3
	env, cm := checkUsableCommon(t, c, tester, cptComp)

	oComp, err := env.Orchestrator.Component()
	assert.Nil(t, err)
	ok, _ := oComp.Templatable()
	if assert.True(t, ok) {
		usableComp := cm.Use(env.Orchestrator)
		assert.True(t, usableComp.Templated())
		// Check the existence of the templated folder
		assert.Equal(t, cptComp+1, tester.countComponent())
		assert.True(t, tester.rootContainsComponent(usableComp.RootPath()))

		// Check  that the two file content has been templated
		checkFile(t, usableComp, "templateTarget1.yaml", "templateContentFromCli")
		checkFile(t, usableComp, "templateTarget2.yaml", "templateContentFromCli")

		// Check the release
		usableComp.Release()
		assert.Equal(t, cptComp, tester.countComponent())
		assert.False(t, tester.rootContainsComponent(usableComp.RootPath()))
	}
}

func TestUsableTemplateNoMatch(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	p, _ := model.CreateParameters(map[string]interface{}{
		"templateDef": "/noMatchinTarget.yaml",
	})
	tc := model.CreateContext(p)

	c := &MockLaunchContext{locationContent: mainPath, templateContext: tc}
	tester := gitTester(t, c)
	defer tester.clean()

	writecheckUsableCommon(t, tester, mainPath)

	cptComp := 3
	env, cm := checkUsableCommon(t, c, tester, cptComp)

	oComp, err := env.Orchestrator.Component()
	assert.Nil(t, err)
	ok, _ := oComp.Templatable()
	if assert.True(t, ok) {
		usableComp := cm.Use(env.Orchestrator)
		assert.False(t, usableComp.Templated())
		// Check that no templated folder has been created
		assert.Equal(t, cptComp, tester.countComponent())
		assert.True(t, tester.rootContainsComponent(usableComp.RootPath()))
		// Check  that the files content has not been templated
		checkFile(t, usableComp, "templateTarget1.yaml", "{{ .Vars.templateContent }}")
		checkFile(t, usableComp, "templateTarget2.yaml", "{{ .Vars.templateContent }}")
		// Check that the release has no effect
		usableComp.Release()
		assert.Equal(t, cptComp, tester.countComponent())
		assert.True(t, tester.rootContainsComponent(usableComp.RootPath()))
	}
}

func writecheckUsableCommon(t *testing.T, tester *tester, d string) {
	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repDesc := tester.createRep(d)

	repComp1.writeCommit(t, "ekara.yaml", usableComp1Content)
	repComp1.writeCommit(t, "templateTarget1.yaml", `{{ .Vars.templateContent }}`)
	repComp1.writeCommit(t, "templateTarget2.yaml", `{{ .Vars.templateContent }}`)
	repDist.writeCommit(t, "ekara.yaml", usableDistContent)
	repDesc.writeCommit(t, "ekara.yaml", usableDescContent)

}

func checkUsableCommon(t *testing.T, c *MockLaunchContext, tester *tester, initialComp int) (model.Environment, component.ComponentManager) {
	err := tester.initEngine()
	assert.Nil(t, err)
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	tester.assertComponentsContains("__main__", "__ekara__", "comp1")

	cm := c.Ekara().ComponentManager()
	assert.NotNil(t, cm)

	assert.Equal(t, initialComp, tester.countComponent())
	return env, cm
}

func checkFile(t *testing.T, u component.UsableComponent, file, wanted string) {
	b, err := ioutil.ReadFile(filepath.Join(u.RootPath(), file))
	assert.Nil(t, err)
	assert.Equal(t, wanted, string(b))
}

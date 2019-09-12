package engine

import (
	"testing"

	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

var (
	byMatchComp1Content = `
templates:
  - "{{ .Vars.templateDef }}"
`
	byMatchComp2Content = `
templates:
  - "{{ .Vars.templateDef }}"
`
	byMatchParentContent = `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
    comp2:
      repository: ./testdata/gittest/comp2
    comp3:
      repository: ./testdata/gittest/comp3
`
	byMatchDescContent = `
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
    component: comp2
  p2:
    component: comp3

nodes:
  node1:
    instances: 1
    provider:
      name: p1
  node2:
    instances: 1
    provider:
      name: p2  
`
)

func TestByMatchNoMatchOnSearch(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	p, _ := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget1.yaml",
	})

	c := &MockLaunchContext{locationContent: mainPath, data: p}
	tester := gitTester(t, c, false)
	defer tester.clean()

	writecheckByMatchCommon(t, tester, mainPath)

	cptComp := 5
	env, cm := checkByMatchCommon(t, c, tester, cptComp)

	comp1 := env.Orchestrator
	comp2 := env.Providers["p1"]

	matches := cm.ContainsFile("dummy.file", tester.context.engine.Context().data)
	assert.NotNil(t, matches)
	assert.Equal(t, 0, matches.Count())

	//Nothing found the no template folder have been created
	assert.Equal(t, cptComp, tester.countComponent())

	matches = cm.ContainsFile("dummy.file", tester.context.engine.Context().data, comp1, comp2)
	assert.NotNil(t, matches)
	assert.Equal(t, 0, matches.Count())

	//Nothing found the no template folder have been created
	assert.Equal(t, cptComp, tester.countComponent())

}

func TestByMatchOneMatchOnSearch(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	p, _ := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget1.yaml",
	})

	c := &MockLaunchContext{locationContent: mainPath, data: p}
	tester := gitTester(t, c, false)
	defer tester.clean()

	writecheckByMatchCommon(t, tester, mainPath)

	cptComp := 5
	env, cm := checkByMatchCommon(t, c, tester, cptComp)

	comp1 := env.Orchestrator
	comp2 := env.Providers["p1"]
	comp3 := env.Providers["p2"]

	checker := checkByMatch(t, tester, cptComp)

	// 2 matches but only 1 (comp2) templated
	matches := cm.ContainsFile("search2.file", tester.context.engine.Context().data)
	releaseCheck := checker(matches, 2, 1, comp2)
	checkByMatchContent(t, matches, comp2, "templateTarget1.yaml", "templateContentFromCli")
	releaseCheck()

	matches = cm.ContainsFile("search2.file", tester.context.engine.Context().data, comp1, comp2, comp3)
	releaseCheck = checker(matches, 2, 1, comp2)
	checkByMatchContent(t, matches, comp2, "templateTarget1.yaml", "templateContentFromCli")
	releaseCheck()

}

func TestByMatchTwoMatchesOnSearch(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	p, _ := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget1.yaml",
	})

	c := &MockLaunchContext{locationContent: mainPath, data: p}
	tester := gitTester(t, c, false)
	defer tester.clean()

	writecheckByMatchCommon(t, tester, mainPath)

	cptComp := 5
	env, cm := checkByMatchCommon(t, c, tester, cptComp)

	comp1 := env.Orchestrator
	comp2 := env.Providers["p1"]
	comp3 := env.Providers["p2"]

	checker := checkByMatch(t, tester, cptComp)

	// 3 matches but only 2 (comp1,comp2) templated
	matches := cm.ContainsFile("search1.file", tester.context.engine.Context().data)
	releaseCheck := checker(matches, 3, 2, comp1, comp2)
	checkByMatchContent(t, matches, comp1, "templateTarget1.yaml", "templateContentFromCli")
	checkByMatchContent(t, matches, comp2, "templateTarget1.yaml", "templateContentFromCli")
	releaseCheck()

	matches = cm.ContainsFile("search1.file", tester.context.engine.Context().data, comp1, comp2, comp3)
	releaseCheck = checker(matches, 3, 2, comp1, comp2)
	checkByMatchContent(t, matches, comp1, "templateTarget1.yaml", "templateContentFromCli")
	checkByMatchContent(t, matches, comp2, "templateTarget1.yaml", "templateContentFromCli")
	releaseCheck()

}

func checkByMatchContent(t *testing.T, ms component.MatchingPaths, r model.ComponentReferencer, file, wanted string) {
	for _, m := range ms.Paths {
		if m.Component().Name() == r.ComponentName() {
			checkFile(t, m.Component(), file, wanted)
		}
	}
}

func checkByMatch(t *testing.T, tester *tester, initialCpt int) func(ms component.MatchingPaths, length int, templatedCpt int, templated ...model.ComponentReferencer) func() {
	return func(ms component.MatchingPaths, length int, templatedCpt int, templated ...model.ComponentReferencer) func() {
		assert.NotNil(t, ms)
		// Check that the right number of templates has been created
		assert.Equal(t, initialCpt+templatedCpt, tester.countComponent())
		// Check that the right number of matching component has been located
		if assert.Equal(t, length, ms.Count()) {
			for _, r := range templated {
				assert.True(t, hasBeenTemplated(t, ms, r))
			}
		}

		return func() {
			ms.Release()
			//Check that the number of component is back to the initial once after the release
			assert.Equal(t, initialCpt, tester.countComponent())
		}
	}
}

func hasBeenTemplated(t *testing.T, ps component.MatchingPaths, r model.ComponentReferencer) bool {
	for _, v := range ps.Paths {
		if v.Component().Name() == r.ComponentName() {
			return assert.True(t, v.Component().Templated())
		}
	}
	return false
}

func writecheckByMatchCommon(t *testing.T, tester *tester, d string) {
	repParent := tester.createRep("./testdata/gittest/parent")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repComp3 := tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp3")
	repDesc := tester.createRep(d)

	repComp1.writeCommit(t, "ekara.yaml", byMatchComp1Content)
	repComp1.writeCommit(t, "search1.file", "")
	repComp1.writeCommit(t, "templateTarget1.yaml", `{{ .Vars.templateContent }}`)

	repComp2.writeCommit(t, "ekara.yaml", byMatchComp2Content)
	repComp2.writeCommit(t, "search1.file", "")
	repComp2.writeCommit(t, "search2.file", "")
	repComp2.writeCommit(t, "templateTarget1.yaml", `{{ .Vars.templateContent }}`)

	repComp3.writeCommit(t, "search1.file", "")
	repComp3.writeCommit(t, "search2.file", "")

	repParent.writeCommit(t, "ekara.yaml", byMatchParentContent)
	repDesc.writeCommit(t, "ekara.yaml", byMatchDescContent)
}

func checkByMatchCommon(t *testing.T, c *MockLaunchContext, tester *tester, initialComp int) (model.Environment, *component.ComponentManager) {
	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	tester.assertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1", "comp2", "comp3")

	cm := c.Ekara().ComponentManager()
	assert.NotNil(t, cm)

	assert.Equal(t, initialComp, tester.countComponent())
	return env, cm
}

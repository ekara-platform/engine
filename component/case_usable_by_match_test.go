package component

import (
	"testing"

	"github.com/ekara-platform/engine/util"

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
name: ekaraDemoVar
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

	c := util.CreateMockLaunchContextWithData(mainPath, p, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	writecheckByMatchCommon(tester, mainPath)

	cptComp := 5
	env, cm := checkByMatchCommon(tester, cptComp)

	comp1 := env.Orchestrator
	comp2 := env.Providers["p1"]

	matches := cm.ContainsFile("dummy.file", tester.cM.TemplateContext())
	assert.NotNil(t, matches)
	assert.Equal(t, 0, matches.Count())

	//Nothing found the no template folder have been created
	assert.Equal(t, cptComp, tester.ComponentCount())

	matches = cm.ContainsFile("dummy.file", tester.cM.TemplateContext(), comp1, comp2)
	assert.NotNil(t, matches)
	assert.Equal(t, 0, matches.Count())

	//Nothing found the no template folder have been created
	assert.Equal(t, cptComp, tester.ComponentCount())

}

func TestByMatchOneMatchOnSearch(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	p, _ := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget1.yaml",
	})

	c := util.CreateMockLaunchContextWithData(mainPath, p, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	writecheckByMatchCommon(tester, mainPath)

	cptComp := 5
	env, cm := checkByMatchCommon(tester, cptComp)

	comp1 := env.Orchestrator
	comp2 := env.Providers["p1"]
	comp3 := env.Providers["p2"]

	checker := checkByMatch(tester, cptComp)

	// 2 matches but only 1 (comp2) templated
	matches := cm.ContainsFile("search2.file", tester.cM.TemplateContext())
	releaseCheck := checker(matches, 2, 1, comp2)
	checkByMatchContent(tester, matches, comp2, "templateTarget1.yaml", "templateContentFromCli")
	releaseCheck()

	matches = cm.ContainsFile("search2.file", tester.cM.TemplateContext(), comp1, comp2, comp3)
	releaseCheck = checker(matches, 2, 1, comp2)
	checkByMatchContent(tester, matches, comp2, "templateTarget1.yaml", "templateContentFromCli")
	releaseCheck()

}

func TestByMatchTwoMatchesOnSearch(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	p, _ := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget1.yaml",
	})

	c := util.CreateMockLaunchContextWithData(mainPath, p, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	writecheckByMatchCommon(tester, mainPath)

	cptComp := 5
	env, cm := checkByMatchCommon(tester, cptComp)

	comp1 := env.Orchestrator
	comp2 := env.Providers["p1"]
	comp3 := env.Providers["p2"]

	checker := checkByMatch(tester, cptComp)

	// 3 matches but only 2 (comp1,comp2) templated
	matches := cm.ContainsFile("search1.file", tester.cM.TemplateContext())
	releaseCheck := checker(matches, 3, 2, comp1, comp2)
	checkByMatchContent(tester, matches, comp1, "templateTarget1.yaml", "templateContentFromCli")
	checkByMatchContent(tester, matches, comp2, "templateTarget1.yaml", "templateContentFromCli")
	releaseCheck()

	matches = cm.ContainsFile("search1.file", tester.cM.TemplateContext(), comp1, comp2, comp3)
	releaseCheck = checker(matches, 3, 2, comp1, comp2)
	checkByMatchContent(tester, matches, comp1, "templateTarget1.yaml", "templateContentFromCli")
	checkByMatchContent(tester, matches, comp2, "templateTarget1.yaml", "templateContentFromCli")
	releaseCheck()

}

func checkByMatchContent(tester *ComponentTester, ms MatchingPaths, r model.ComponentReferencer, file, wanted string) {
	for _, m := range ms.Paths {
		if m.Component().Name() == r.ComponentName() {
			tester.CheckFile(m.Component(), file, wanted)
		}
	}
}

func checkByMatch(tester *ComponentTester, initialCpt int) func(ms MatchingPaths, length int, templatedCpt int, templated ...model.ComponentReferencer) func() {
	return func(ms MatchingPaths, length int, templatedCpt int, templated ...model.ComponentReferencer) func() {
		assert.NotNil(tester.t, ms)
		// Check that the right number of templates has been created
		assert.Equal(tester.t, initialCpt+templatedCpt, tester.ComponentCount())
		// Check that the right number of matching component has been located
		if assert.Equal(tester.t, length, ms.Count()) {
			for _, r := range templated {
				assert.True(tester.t, hasBeenTemplated(tester.t, ms, r))
			}
		}

		return func() {
			ms.Release()
			//Check that the number of component is back to the initial once after the release
			assert.Equal(tester.t, initialCpt, tester.ComponentCount())
		}
	}
}

func hasBeenTemplated(t *testing.T, ps MatchingPaths, r model.ComponentReferencer) bool {
	for _, v := range ps.Paths {
		if v.Component().Name() == r.ComponentName() {
			return assert.True(t, v.Component().Templated())
		}
	}
	return false
}

func writecheckByMatchCommon(tester *ComponentTester, d string) {
	repParent := tester.CreateRep("./testdata/gittest/parent")
	repComp1 := tester.CreateRep("./testdata/gittest/comp1")
	repComp2 := tester.CreateRep("./testdata/gittest/comp2")
	repComp3 := tester.CreateRepDefaultDescriptor("./testdata/gittest/comp3")
	repDesc := tester.CreateRep(d)

	repComp1.WriteCommit("ekara.yaml", byMatchComp1Content)
	repComp1.WriteCommit("search1.file", "")
	repComp1.WriteCommit("templateTarget1.yaml", `{{ .Vars.templateContent }}`)

	repComp2.WriteCommit("ekara.yaml", byMatchComp2Content)
	repComp2.WriteCommit("search1.file", "")
	repComp2.WriteCommit("search2.file", "")
	repComp2.WriteCommit("templateTarget1.yaml", `{{ .Vars.templateContent }}`)

	repComp3.WriteCommit("search1.file", "")
	repComp3.WriteCommit("search2.file", "")

	repParent.WriteCommit("ekara.yaml", byMatchParentContent)
	repDesc.WriteCommit("ekara.yaml", byMatchDescContent)
}

func checkByMatchCommon(tester *ComponentTester, initialComp int) (model.Environment, ComponentManager) {
	err := tester.Init()
	assert.Nil(tester.t, err)
	env := tester.Env()
	assert.NotNil(tester.t, env)

	tester.AssertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1", "comp2", "comp3")

	cm := tester.cM
	assert.NotNil(tester.t, cm)

	assert.Equal(tester.t, initialComp, tester.ComponentCount())
	return env, cm
}

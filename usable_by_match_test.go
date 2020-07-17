package engine

import (
	"github.com/GroupePSA/componentizer"
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/engine/model"

	"github.com/stretchr/testify/assert"
)

var (
	byMatchComp1Content = `
ekara:
  templates:
    - "{{ .Vars.templateDef }}"
`
	byMatchComp2Content = `
ekara:
  templates:
    - "{{ .Vars.templateDef }}"
`
	byMatchParentContent = `
ekara:
  components:
    comp1:
      repository: comp1
    comp2:
      repository: comp2
    comp3:
      repository: comp3
`
	byMatchDescContent = `
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

	mainPath := "descriptor"

	p := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget1.yaml",
	})

	tester := util.CreateComponentTester(t, p)
	defer tester.Clean()

	repo := writeCheckByMatchCommon(tester, mainPath)

	cptComp := 5
	env, cm := checkByMatchCommon(tester, repo, cptComp)

	comp1 := env.Orchestrator
	comp2 := env.Providers["p1"]

	matches := cm.ContainsFile("dummy.file", tester.TemplateContext())
	assert.NotNil(t, matches)
	assert.Equal(t, 0, matches.Count())

	//Nothing found the no template folder have been created
	assert.Equal(t, cptComp, tester.ComponentCount())

	matches = cm.ContainsFile("dummy.file", tester.TemplateContext(), comp1, comp2)
	assert.NotNil(t, matches)
	assert.Equal(t, 0, matches.Count())

	//Nothing found the no template folder have been created
	assert.Equal(t, cptComp, tester.ComponentCount())

}

func TestByMatchOneMatchOnSearch(t *testing.T) {

	mainPath := "descriptor"

	p := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget1.yaml",
	})

	tester := util.CreateComponentTester(t, p)
	defer tester.Clean()

	repo := writeCheckByMatchCommon(tester, mainPath)

	cptComp := 5
	env, cm := checkByMatchCommon(tester, repo, cptComp)

	comp1 := env.Orchestrator
	comp2 := env.Providers["p1"]
	comp3 := env.Providers["p2"]

	checker := checkByMatch(tester, cptComp)

	// 2 matches but only 1 (comp2) templated
	matches := cm.ContainsFile("search2.file", tester.TemplateContext())
	releaseCheck := checker(matches, 2, comp2)
	checkByMatchContent(tester, matches, comp2, "templateTarget1.yaml", "templateContentFromCli")
	releaseCheck()

	matches = cm.ContainsFile("search2.file", tester.TemplateContext(), comp1, comp2, comp3)
	releaseCheck = checker(matches, 2, comp2)
	checkByMatchContent(tester, matches, comp2, "templateTarget1.yaml", "templateContentFromCli")
	releaseCheck()

}

func TestByMatchTwoMatchesOnSearch(t *testing.T) {

	mainPath := "descriptor"

	p := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget1.yaml",
	})

	tester := util.CreateComponentTester(t, p)
	defer tester.Clean()

	repo := writeCheckByMatchCommon(tester, mainPath)

	cptComp := 5
	env, cm := checkByMatchCommon(tester, repo, cptComp)

	comp1 := env.Orchestrator
	comp2 := env.Providers["p1"]
	comp3 := env.Providers["p2"]

	checker := checkByMatch(tester, cptComp)

	// 3 matches but only 2 (comp1,comp2) templated
	matches := cm.ContainsFile("search1.file", tester.TemplateContext())
	releaseCheck := checker(matches, 3, comp1, comp2)
	checkByMatchContent(tester, matches, comp1, "templateTarget1.yaml", "templateContentFromCli")
	checkByMatchContent(tester, matches, comp2, "templateTarget1.yaml", "templateContentFromCli")
	releaseCheck()

	matches = cm.ContainsFile("search1.file", tester.TemplateContext(), comp1, comp2, comp3)
	releaseCheck = checker(matches, 3, comp1, comp2)
	checkByMatchContent(tester, matches, comp1, "templateTarget1.yaml", "templateContentFromCli")
	checkByMatchContent(tester, matches, comp2, "templateTarget1.yaml", "templateContentFromCli")
	releaseCheck()

}

func checkByMatchContent(tester util.EkaraComponentTester, ms componentizer.MatchingPaths, r componentizer.ComponentRef, file, wanted string) {
	for _, m := range ms.Paths {
		if m.Owner().Id() == r.ComponentId() {
			tester.AssertFileContent(m.Owner(), file, wanted)
		}
	}
}

func checkByMatch(tester util.EkaraComponentTester, initialCpt int) func(ms componentizer.MatchingPaths, length int, templated ...componentizer.ComponentRef) func() {
	return func(ms componentizer.MatchingPaths, length int, templated ...componentizer.ComponentRef) func() {
		assert.NotNil(tester.T(), ms)
		// Check that the right number of templates has been created
		assert.Equal(tester.T(), initialCpt, tester.ComponentCount())
		// Check that the right number of matching component has been located
		if assert.Equal(tester.T(), length, ms.Count()) {
			for _, r := range templated {
				assert.True(tester.T(), hasBeenTemplated(tester.T(), ms, r))
			}
		}

		return func() {
			ms.Release()
			//Check that the number of component is back to the initial once after the release
			assert.Equal(tester.T(), initialCpt, tester.ComponentCount())
		}
	}
}

func hasBeenTemplated(t *testing.T, ps componentizer.MatchingPaths, r componentizer.ComponentRef) bool {
	for _, v := range ps.Paths {
		if v.Owner().Id() == r.ComponentId() {
			return assert.True(t, v.Owner().Templated())
		}
	}
	return false
}

func writeCheckByMatchCommon(tester util.EkaraComponentTester, d string) componentizer.Repository {
	repParent := tester.CreateDir("parent")
	repComp1 := tester.CreateDir("comp1")
	repComp2 := tester.CreateDir("comp2")
	repComp3 := tester.CreateDirEmptyDesc("comp3")
	repDesc := tester.CreateDir(d)

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

	return repDesc.AsRepository("master")
}

func checkByMatchCommon(tester util.EkaraComponentTester, repo componentizer.Repository, initialComp int) (model.Environment, componentizer.ComponentManager) {
	tester.Init(repo)
	env := tester.Env()
	assert.NotNil(tester.T(), env)

	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1", "comp2", "comp3")

	assert.Equal(tester.T(), initialComp, tester.ComponentCount())
	return env, tester.ComponentManager()
}

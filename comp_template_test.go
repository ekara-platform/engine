package engine

import (
	"github.com/ekara-platform/engine/util"
	"testing"

	"github.com/ekara-platform/engine/model"

	"github.com/stretchr/testify/assert"
)

func TestComponentTemplatable(t *testing.T) {

	comp1Content := `
ekara:
  templates:
    - "*.yml"
`

	parentContent := `
ekara:
  components:
    comp1:
      repository: comp1
    comp2:
      repository: comp2	  
`
	descContent := `
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

nodes:
  node1:
    instances: 1
    provider:
      name: p1
`

	p := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget1.yaml",
	})
	tester := util.CreateComponentTester(t, p)
	defer tester.Clean()

	repParent := tester.CreateDir("parent")
	repComp1 := tester.CreateDir("comp1")
	tester.CreateDirEmptyDesc("comp2")
	repDesc := tester.CreateDir("descriptor")
	repComp1.WriteCommit("ekara.yaml", comp1Content)
	repParent.WriteCommit("ekara.yaml", parentContent)
	repDesc.WriteCommit("ekara.yaml", descContent)
	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)

	tester.AssertComponentAvailable(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1", "comp2")
	oComp, err := env.Orchestrator.Component(env)
	assert.Nil(t, err)
	ok, patterns := oComp.GetTemplates()
	if assert.True(t, ok) {
		assert.Contains(t, patterns, "*.yml")
	}

	pComp, err := env.Providers["p1"].Component(env)
	assert.Nil(t, err)
	ok, patterns = pComp.GetTemplates()
	if assert.False(t, ok) {
		assert.Equal(t, 0, len(patterns))
	}
}

package component

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestComponentTemplatable(t *testing.T) {

	comp1Content := `
templates:
  - "*.yml"
`

	parentContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
    comp2:
      repository: ./testdata/gittest/comp2	  
`
	descContent := `
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

nodes:
  node1:
    instances: 1
    provider:
      name: p1
`

	mainPath := "./testdata/gittest/descriptor"

	p, _ := model.CreateParameters(map[string]interface{}{
		"templateContent": "templateContentFromCli",
		"templateDef":     "/templateTarget1.yaml",
	})

	c := util.CreateMockLaunchContextWithData(mainPath, p, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	repParent := tester.CreateRep("./testdata/gittest/parent")
	repComp1 := tester.CreateRep("./testdata/gittest/comp1")
	tester.CreateRepDefaultDescriptor("./testdata/gittest/comp2")
	repDesc := tester.CreateRep(mainPath)

	repComp1.WriteCommit("ekara.yaml", comp1Content)
	repParent.WriteCommit("ekara.yaml", parentContent)
	repDesc.WriteCommit("ekara.yaml", descContent)

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)

	tester.AssertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1", "comp2")

	cm := tester.cM
	assert.NotNil(t, cm)

	oComp, err := env.Orchestrator.Component()
	assert.Nil(t, err)
	ok, patterns := oComp.Templatable()
	if assert.True(t, ok) {
		assert.Contains(t, patterns.Content, "*.yml")
	}

	pComp, err := env.Providers["p1"].Component()
	assert.Nil(t, err)
	ok, patterns = pComp.Templatable()
	if assert.False(t, ok) {
		assert.Equal(t, 0, len(patterns.Content))
	}
}

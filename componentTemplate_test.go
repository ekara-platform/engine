package engine

import (
	"testing"

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

	c := &MockLaunchContext{locationContent: mainPath, data: p}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repParent := tester.createRep("./testdata/gittest/parent")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

	repComp1.writeCommit(t, "ekara.yaml", comp1Content)
	repParent.writeCommit(t, "ekara.yaml", parentContent)
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	tester.assertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1", "comp2")

	cm := c.Ekara().ComponentManager()
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

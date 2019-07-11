package engine

import (
	"testing"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestOverwriteInMain(t *testing.T) {

	distContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
`
	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/parent")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp1Overwritten := tester.createRep("./testdata/gittest/comp1Overwritten")
	repDesc := tester.createRep(mainPath)

	repDist.writeCommit(t, "ekara.yaml", distContent)

	repComp1.writeCommit(t, "content.txt", "comp content from parent")
	repComp1Overwritten.writeCommit(t, "content.txt", "comp content overwriten in descriptor")

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent	
  components:
    comp1:
      repository: ./testdata/gittest/comp1Overwritten
# Following content just to force the download of comp1
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
	repDesc.writeCommit(t, "ekara.yaml", descContent)
	
	err := tester.initEngine()
	assert.Nil(t, err)
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator and provider
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId, "comp1")

	cm := c.Ekara().ComponentManager()
	assert.NotNil(t, cm)

	usableComp, err := cm.Use(env.Orchestrator)
	defer usableComp.Release()
	// Check that the comp1 used is the one defined into the main descriptor
	checkFile(t, usableComp, "content.txt", "comp content overwriten in descriptor")
}

func TestOverwriteInDescriptorComponent(t *testing.T) {

	comp2Content := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
`

	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/parent")
	repDist.writeCommit(t, "ekara.yaml", "")

	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repComp2.writeCommit(t, "ekara.yaml", comp2Content)

	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp1.writeCommit(t, "content.txt", "comp content from parent")

	repComp1Overwritten := tester.createRep("./testdata/gittest/comp1Overwritten")
	repComp1Overwritten.writeCommit(t, "content.txt", "comp content overwriten in descriptor")

	repDesc := tester.createRep(mainPath)

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent
  components:
    comp1:
      repository: ./testdata/gittest/comp1Overwritten
    comp2:
      repository: ./testdata/gittest/comp2

# Following content just to force the download of comp1
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
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator
	// comp2 should be downloaded because it's used as  provider
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId, "comp1", "comp2")
	cm := c.Ekara().ComponentManager()
	assert.NotNil(t, cm)

	usableComp, err := cm.Use(env.Orchestrator)
	defer usableComp.Release()
	// Check that the comp1 used is the one defined into the main descriptor
	checkFile(t, usableComp, "content.txt", "comp content overwriten in descriptor")

}

package component

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

// Test that a component defined into the parent can be overwritten
// into the main descriptor.
func TestOverwritenFromParentByMain(t *testing.T) {
	parentContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
`
	mainPath := "./testdata/gittest/descriptor"
	c := util.CreateMockLaunchContext(mainPath, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	repParent := tester.CreateRep("./testdata/gittest/parent")
	repComp1 := tester.CreateRep("./testdata/gittest/comp1")
	repComp1Overwritten := tester.CreateRep("./testdata/gittest/comp1Overwritten")
	repDesc := tester.CreateRep(mainPath)

	repParent.WriteCommit("ekara.yaml", parentContent)

	repComp1.WriteCommit("content.txt", "comp content from parent")
	repComp1Overwritten.WriteCommit("content.txt", "comp content overwriten in descriptor")

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
	repDesc.WriteCommit("ekara.yaml", descContent)

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator and provider
	tester.AssertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", "comp1")

	cm := tester.cM
	assert.NotNil(t, cm)

	usableComp, err := cm.Use(env.Orchestrator, tester.cM.tplC)
	defer usableComp.Release()
	// Check that the comp1 used is the one defined into the main descriptor
	tester.CheckFile(usableComp, "content.txt", "comp content overwriten in descriptor")
}

// Test that a component defined into a parent can be overwritten
// into a child.
func TestOverwritenFromParentByChild(t *testing.T) {

	parent2Content := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
`

	parent1Content := `
ekara:
  parent:
    repository: ./testdata/gittest/parent2	
  components:
    comp1:
      repository: ./testdata/gittest/comp1Overwritten
`

	mainPath := "./testdata/gittest/descriptor"
	c := util.CreateMockLaunchContext(mainPath, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	repParent1 := tester.CreateRep("./testdata/gittest/parent1")
	repParent2 := tester.CreateRep("./testdata/gittest/parent2")
	repComp1 := tester.CreateRep("./testdata/gittest/comp1")
	repComp1Overwritten := tester.CreateRep("./testdata/gittest/comp1Overwritten")
	repDesc := tester.CreateRep(mainPath)

	repParent1.WriteCommit("ekara.yaml", parent1Content)
	repParent2.WriteCommit("ekara.yaml", parent2Content)

	repComp1.WriteCommit("content.txt", "comp content from parent")
	repComp1Overwritten.WriteCommit("content.txt", "comp content overwriten in descriptor")

	descContent := `
name: ekara-demo-var
qualifier: dev
ekara:
  parent:
    repository: ./testdata/gittest/parent1	

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
	repDesc.WriteCommit("ekara.yaml", descContent)

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator and provider
	tester.AssertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", model.EkaraComponentId+"2", "comp1")

	cm := tester.cM
	assert.NotNil(t, cm)

	usableComp, err := cm.Use(env.Orchestrator, tester.cM.tplC)
	defer usableComp.Release()
	// Check that the comp1 used is the one defined into the main descriptor
	tester.CheckFile(usableComp, "content.txt", "comp content overwriten in descriptor")
}

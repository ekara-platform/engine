package engine

import (
	"github.com/ekara-platform/engine/util"
	"testing"

	"github.com/ekara-platform/engine/model"
	"github.com/stretchr/testify/assert"
)

// Test that a component defined into the parent can be overwritten
// into the main descriptor.
func TestOverwrittenFromParentByMain(t *testing.T) {
	parentContent := `
ekara:
  components:
    comp1:
      repository: comp1
`
	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	repParent := tester.CreateDir("parent")
	repComp1 := tester.CreateDir("comp1")
	repComp1Overwritten := tester.CreateDir("comp1Overwritten")
	repDesc := tester.CreateDir("descriptor")

	repParent.WriteCommit("ekara.yaml", parentContent)

	repComp1.WriteCommit("content.txt", "comp content from parent")
	repComp1Overwritten.WriteCommit("content.txt", "comp content overwritten in descriptor")

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: parent	
  components:
    comp1:
      repository: comp1Overwritten
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
	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator and provider
	tester.AssertComponentAvailable(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1")

	usableComp, err := tester.ComponentManager().Use(env.Orchestrator, tester.TemplateContext())
	assert.Nil(t, err)
	defer usableComp.Release()
	// Check that the comp1 used is the one defined into the main descriptor
	tester.AssertFileContent(usableComp, "content.txt", "comp content overwritten in descriptor")
}

// Test that a component defined into a parent can be overwritten
// into a child.
func TestOverwritenFromParentByChild(t *testing.T) {

	parent2Content := `
ekara:
  components:
    comp1:
      repository: comp1
`

	parent1Content := `
ekara:
  parent:
    repository: parent2	
  components:
    comp1:
      repository: comp1Overwritten
`

	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	repParent1 := tester.CreateDir("parent1")
	repParent2 := tester.CreateDir("parent2")
	repComp1 := tester.CreateDir("comp1")
	repComp1Overwritten := tester.CreateDir("comp1Overwritten")
	repDesc := tester.CreateDir("descriptor")

	repParent1.WriteCommit("ekara.yaml", parent1Content)
	repParent2.WriteCommit("ekara.yaml", parent2Content)

	repComp1.WriteCommit("content.txt", "comp content from parent")
	repComp1Overwritten.WriteCommit("content.txt", "comp content overwritten in descriptor")

	descContent := `
name: ekaraDemoVar
qualifier: dev
ekara:
  parent:
    repository: parent1	

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
	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator and provider
	tester.AssertComponentAvailable(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, model.MainComponentId+model.ParentComponentSuffix+model.ParentComponentSuffix, "comp1")

	usableComp, err := tester.ComponentManager().Use(env.Orchestrator, tester.TemplateContext())
	assert.Nil(t, err)
	defer usableComp.Release()
	// Check that the comp1 used is the one defined into the main descriptor
	tester.AssertFileContent(usableComp, "content.txt", "comp content overwritten in descriptor")
}

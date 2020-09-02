package engine

import (
	"github.com/ekara-platform/engine/model"
	"github.com/ekara-platform/engine/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

// A component declared into another component must be ignored
func TestComponentInComponent(t *testing.T) {

	comp2Content := `
ekara:
  components:
    comp1:
      repository: comp1
`

	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	tester.CreateDirEmptyDesc("parent")

	repComp2 := tester.CreateDir("comp2")
	repComp2.WriteCommit("ekara.yaml", comp2Content)

	repComp1 := tester.CreateDir("comp1")
	repComp1.WriteCommit("content.txt", "comp content from parent")

	repDesc := tester.CreateDir("descriptor")

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: parent
  components:
    comp2:
      repository: comp2
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
	repDesc.WriteCommit("ekara.yaml", descContent)
	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)
	tester.AssertComponentAvailable(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp2")
}

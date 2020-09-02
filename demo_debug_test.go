package engine

import (
	_ "log"
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/engine/model"
)

func TestDebugDemo(t *testing.T) {
	parent2Content := `
ekara:
  components:
    comp2:
      repository: comp2
`

	parent1Content := `
ekara:
  parent:
    repository: parent2
  components:
    comp1:
      repository: comp1
`

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: parent1
providers:
  p1:
    component: comp1
nodes:
  node1:
    instances: 1
    provider:
      name: p1
`

	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	repParent1 := tester.CreateDir("parent1")
	repParent2 := tester.CreateDir("parent2")
	tester.CreateDirEmptyDesc("comp1")
	tester.CreateDirEmptyDesc("comp2")
	repDesc := tester.CreateDir("descriptor")
	repParent1.WriteCommit("ekara.yaml", parent1Content)
	repParent2.WriteCommit("ekara.yaml", parent2Content)
	repDesc.WriteCommit("ekara.yaml", descContent)
	tester.Init(repDesc.AsRepository("master"))

	// comp1 should be downloaded because it's used as orchestrator into the parent
	// comp2 should not be downloaded because it's not referenced by a component
	tester.AssertComponentAvailable(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1")
	tester.AssertComponentMissing("comp2")
}

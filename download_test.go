package engine

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/engine/model"

	"github.com/stretchr/testify/assert"
)

func TestDownloadOnlyUsedComponents(t *testing.T) {
	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	repParent := tester.CreateDir("parent")
	tester.CreateDirEmptyDesc("comp1")
	tester.CreateDirEmptyDesc("comp2")
	tester.CreateDirEmptyDesc("comp3")
	tester.CreateDirEmptyDesc("comp4")
	repDesc := tester.CreateDir("descriptor")

	parentContent := `
ekara:
  components:
    comp1:
      repository: comp1
    comp2:
      repository: comp2
`
	repParent.WriteCommit("ekara.yaml", parentContent)

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: parent
  components:
    comp3:
      repository: comp3
    comp4:
      repository: comp4
`
	repDesc.WriteCommit("ekara.yaml", descContent)
	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)

	// comp1, comp2, comp3 and comp4 shouldn't be downloaded because they are not used into the descriptor
	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix)
}

func TestDownloadOnlyFromDescriptorAndParents(t *testing.T) {
	comp1Content := `
ekara:
 components:
   comp2:
     repository: comp2
`

	parentContent := `
ekara:
 components:
   comp1:
     repository: comp1
orchestrator:
 component: comp1
`
	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
 parent:
   repository: parent

# Following content just to force the download of comp1 and comp2
providers:
 p1:
   component: comp2

nodes:
 node1:
   instances: 1
   provider:
     name: p1
`
	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	repParent := tester.CreateDir("parent")
	repComp1 := tester.CreateDir("comp1")
	tester.CreateDirEmptyDesc("comp2")
	repDesc := tester.CreateDir("descriptor")

	repDesc.WriteCommit("ekara.yaml", descContent)
	repParent.WriteCommit("ekara.yaml", parentContent)
	repComp1.WriteCommit("ekara.yaml", comp1Content)

	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)

	// comp1 should be downloaded because it's used as orchestrator into the parent
	// comp2 should not be downloaded because it's referenced by a component
	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1")
}

func TestDownloadTwoParents(t *testing.T) {
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

# Following content just to force the download of comp1 and comp2
orchestrator:
 component: comp2
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
	env := tester.Env()
	assert.NotNil(t, env)

	// comp1 should be downloaded because it's used as orchestrator
	// comp2 should be also downloaded because it's used as provider
	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, model.MainComponentId+model.ParentComponentSuffix+model.ParentComponentSuffix, "comp1", "comp2")
}

func TestDownloadTwoParentsUpperUsed(t *testing.T) {
	parent2Content := `
ekara:
 components:
   comp2:
     repository: comp2
orchestrator:
 component: comp2
providers:
 p1:
   component: comp1
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

# Following content just to force the download of comp1 and comp2
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
	env := tester.Env()
	assert.NotNil(t, env)

	// comp1 should be downloaded because it's used as orchestrator
	// comp2 should be also downloaded because it's used as provider
	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, model.MainComponentId+model.ParentComponentSuffix+model.ParentComponentSuffix, "comp1", "comp2")
}

// The orchestrator and the providers, once defined into a parent
// can be customized lower into the hierarchy of descriptors.
// In this case the component will be omitted, because it has already been defined,
// and the should not affect the download process
func TestDownloadTwoParentsProviderUpperUsedRedefined(t *testing.T) {
	parent2Content := `
ekara:
 components:
   comp2:
     repository: comp2
orchestrator:
 component: comp2
providers:
 p1:
   component: comp1
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

orchestrator:
 params:
   dummy_key: dummy_value
providers:
 p1:
   params:
     dummy_key: dummy_value


# Following content just to force the download of comp1 and comp2
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
	env := tester.Env()
	assert.NotNil(t, env)

	// comp1 should be downloaded because it's used as orchestrator
	// comp2 should be also downloaded because it's used as provider
	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, model.MainComponentId+model.ParentComponentSuffix+model.ParentComponentSuffix, "comp1", "comp2")
}

// A stack  once defined into a parent
// can be customized lower into the hierarchy of descriptors.
// In this case the component will be omitted, because it has already been defined,
// and the should not affect the download process
func TestDownloadTwoParentsStackUpperUsedRedefined(t *testing.T) {
	parent2Content := `
ekara:
 components:
   comp1:
     repository: comp1
   comp2:
     repository: comp2
   comp3:
     repository: comp3

orchestrator:
 component: comp1
providers:
 p1:
   component: comp2
stacks:
 s1:
   component: comp3
`
	parent1Content := `
ekara:
 parent:
   repository: parent2
`

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
 parent:
   repository: parent1

stacks:
 s1:
   params:
     dummy_key: dummy_value

# Following content just to force the download of comp1 and comp2
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
	tester.CreateDirEmptyDesc("comp3")
	repDesc := tester.CreateDir("descriptor")

	repParent1.WriteCommit("ekara.yaml", parent1Content)
	repParent2.WriteCommit("ekara.yaml", parent2Content)
	repDesc.WriteCommit("ekara.yaml", descContent)

	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)

	// comp1 should be downloaded because it's used as orchestrator
	// comp2 should be also downloaded because it's used as provider
	// comp3 should be also downloaded because it's used as a stack
	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, model.MainComponentId+model.ParentComponentSuffix+model.ParentComponentSuffix, "comp1", "comp2", "comp3")
}

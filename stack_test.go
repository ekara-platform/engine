package engine

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/engine/model"

	"github.com/stretchr/testify/assert"
)

func TestStackFromDescriptorAndParent(t *testing.T) {
	parentContent := `
ekara:
  components:
    comp1:
      repository: comp1
stacks:
  parentStack:
`
	mainPath := "descriptor"
	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	repParent := tester.CreateDir("parent")
	tester.CreateDirEmptyDesc("comp1")
	repDesc := tester.CreateDir(mainPath)

	repParent.WriteCommit("ekara.yaml", parentContent)
	// write the compose/playbook content into the parent component
	repParent.WriteCommit("docker_compose.yml", "parent docker compose content")

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: parent
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
stacks:
  descriptorStack:
`

	repDesc.WriteCommit("ekara.yaml", descContent)
	// write the compose/playbook content into the descriptor component
	repDesc.WriteCommit("docker_compose.yml", "descriptor docker compose content")

	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator and provider
	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1")

	// Check that the environment has two stacks
	if assert.Equal(t, 2, len(env.Stacks)) {
		tester.CheckStack(model.MainComponentId, "descriptorStack", "descriptor docker compose content")
		tester.CheckStack(model.MainComponentId+model.ParentComponentSuffix, "parentStack", "parent docker compose content")
	}
}

func TestStackThroughParent(t *testing.T) {
	comp1Content := `
stacks:
  comp1Stack:
`

	parentContent := `
ekara:
  components:
    comp1:
      repository: comp1
`
	mainPath := "descriptor"
	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	repParent := tester.CreateDir("parent")
	repComp1 := tester.CreateDir("comp1")
	// write the compose/playbook content into the comp1 component
	repComp1.WriteCommit("docker_compose.yml", "comp1 docker compose content")
	repDesc := tester.CreateDir(mainPath)

	repComp1.WriteCommit("ekara.yaml", comp1Content)
	repParent.WriteCommit("ekara.yaml", parentContent)

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: parent
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
	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1")

	// Check that the environment has no stacks
	assert.Equal(t, 1, len(env.Stacks))
	tester.CheckStack("comp1", "comp1Stack", "comp1 docker compose content")
}

func TestStackThroughComponent(t *testing.T) {

	comp2Content := `
stacks:
  comp2Stack:
`

	parentContent := `
ekara:
  components:
    comp1:
      repository: comp1
`
	mainPath := "descriptor"
	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	repParent := tester.CreateDir("parent")
	tester.CreateDirEmptyDesc("comp1")
	repComp2 := tester.CreateDir("comp2")
	// write the compose/playbook content into the comp2 component
	repComp2.WriteCommit("docker_compose.yml", "comp2 docker compose content")
	repDesc := tester.CreateDir(mainPath)

	repComp2.WriteCommit("ekara.yaml", comp2Content)
	repParent.WriteCommit("ekara.yaml", parentContent)

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

	repDesc.WriteCommit("ekara.yaml", descContent)

	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used asprovider
	// comp2 should be downloaded because it's used as orchestrator
	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1", "comp2")

	// Check that the environment has no stacks
	assert.Equal(t, 1, len(env.Stacks))
	tester.CheckStack("comp2", "comp2Stack", "comp2 docker compose content")
}

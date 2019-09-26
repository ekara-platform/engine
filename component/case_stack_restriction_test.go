package component

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestStackFromDesciptorAndParent(t *testing.T) {

	parentContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
stacks:
  parentStack:
`
	mainPath := "./testdata/gittest/descriptor"
	c := util.CreateMockLaunchContext(mainPath, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	repParent := tester.CreateRep("./testdata/gittest/parent")
	tester.CreateRepDefaultDescriptor("./testdata/gittest/comp1")
	repDesc := tester.CreateRep(mainPath)

	repParent.WriteCommit("ekara.yaml", parentContent)
	// write the compose/playbook content into the parent component
	repParent.WriteCommit("docker_compose.yml", "parent docker compose content")

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent
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

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator and provider
	tester.AssertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", "comp1")

	// Chect that the environment has two stacks
	if assert.Equal(t, 2, len(env.Stacks)) {

		cm := tester.cM
		assert.NotNil(t, cm)
		tester.CheckStack(model.MainComponentId, "descriptorStack", "descriptor docker compose content")
		tester.CheckStack(model.EkaraComponentId+"1", "parentStack", "parent docker compose content")
	}
}

func TestIgnoredStackThroughParent(t *testing.T) {

	comp1Content := `
stacks:
  comp1Stack:
`

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
	// write the compose/playbook content into the comp1 component
	repComp1.WriteCommit("docker_compose.yml", "comp1 docker compose content")
	repDesc := tester.CreateRep(mainPath)

	repComp1.WriteCommit("ekara.yaml", comp1Content)
	repParent.WriteCommit("ekara.yaml", parentContent)

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent
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

	// Chect that the environment has no stacks
	assert.Equal(t, 0, len(env.Stacks))
}

func TestIgnoredStackThroughComponent(t *testing.T) {

	comp2Content := `
stacks:
  comp2Stack:
`

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
	tester.CreateRepDefaultDescriptor("./testdata/gittest/comp1")
	repComp2 := tester.CreateRep("./testdata/gittest/comp2")
	// write the compose/playbook content into the comp2 component
	repComp2.WriteCommit("docker_compose.yml", "comp2 docker compose content")
	repDesc := tester.CreateRep(mainPath)

	repComp2.WriteCommit("ekara.yaml", comp2Content)
	repParent.WriteCommit("ekara.yaml", parentContent)

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent
  components:
    comp2:
      repository: ./testdata/gittest/comp2
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

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used asprovider
	// comp2 should be downloaded because it's used as orchestrator
	tester.AssertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", "comp1", "comp2")

	// Chect that the environment has no stacks
	assert.Equal(t, 0, len(env.Stacks))
}

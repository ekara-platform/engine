package engine

import (
	"testing"

	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestStackFromDesciptorAndParent(t *testing.T) {

	distContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
stacks:
  parentStack:
`
	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/parent")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repDesc := tester.createRep(mainPath)

	repComp1.writeCommit(t, "ekara.yaml", "")
	repDist.writeCommit(t, "ekara.yaml", distContent)
	// write the compose/playbook content into the parent component
	repDist.writeCommit(t, "docker_compose.yml", "parent docker compose content")

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

	repDesc.writeCommit(t, "ekara.yaml", descContent)
	// write the compose/playbook content into the descriptor component
	repDesc.writeCommit(t, "docker_compose.yml", "descriptor docker compose content")

	err := tester.initEngine()
	assert.Nil(t, err)
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator and provider
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId, "comp1")

	// Chect that the enviroment has two stacks
	if assert.Equal(t, 2, len(env.Stacks)) {

		cm := c.Ekara().ComponentManager()
		assert.NotNil(t, cm)
		checkStack(t, env, cm, model.MainComponentId, "descriptorStack", "descriptor docker compose content")
		checkStack(t, env, cm, model.EkaraComponentId, "parentStack", "parent docker compose content")
	}
}

func TestIgnoredStackThroughParent(t *testing.T) {

	comp1Content := `
stacks:
  comp1Stack:
`

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
	// write the compose/playbook content into the comp1 component
	repComp1.writeCommit(t, "docker_compose.yml", "comp1 docker compose content")
	repDesc := tester.createRep(mainPath)

	repComp1.writeCommit(t, "ekara.yaml", comp1Content)
	repDist.writeCommit(t, "ekara.yaml", distContent)

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

	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator and provider
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId, "comp1")

	// Chect that the enviroment has no stacks
	assert.Equal(t, 0, len(env.Stacks))
}

func TestIgnoredStackThroughComponent(t *testing.T) {

	comp2Content := `
stacks:
  comp2Stack:
`

	distContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
`
	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
	//defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/parent")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp2 := tester.createRep("./testdata/gittest/comp2")
	// write the compose/playbook content into the comp2 component
	repComp2.writeCommit(t, "docker_compose.yml", "comp2 docker compose content")
	repDesc := tester.createRep(mainPath)

	repComp1.writeCommit(t, "ekara.yaml", "")
	repComp2.writeCommit(t, "ekara.yaml", comp2Content)
	repDist.writeCommit(t, "ekara.yaml", distContent)

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

	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used asprovider
	// comp2 should be downloaded because it's used as orchestrator
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId, "comp1", "comp2")

	// Chect that the enviroment has no stacks
	assert.Equal(t, 0, len(env.Stacks))
}

func checkStack(t *testing.T, env model.Environment, cm component.ComponentManager, holder, stackName, compose string) {
	stack, ok := env.Stacks[stackName]
	if assert.True(t, ok) {
		//Check that the self contained stack has been well built
		assert.Equal(t, stackName, stack.Name)
		stackC, err := stack.Component()
		assert.Nil(t, err)
		assert.NotNil(t, stackC)
		assert.Equal(t, holder, stackC.Id)

		// Check that the stack is usable and returns the correct component
		usableStack, err := cm.Use(stack)
		defer usableStack.Release()
		assert.Nil(t, err)
		assert.NotNil(t, usableStack)
		assert.False(t, usableStack.Templated())
		// Check that the stacks contains the compose/playbook file
		checkFile(t, usableStack, "docker_compose.yml", compose)
	}
}

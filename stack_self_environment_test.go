package engine

import (
	"testing"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentSelfStackNoComponent(t *testing.T) {

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  distribution:
    repository: ./testdata/gittest/distribution
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
  myStack:
`
	checkSelfStack(t, descContent)
}

func TestEnvironmentSelfStackLowDash(t *testing.T) {

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  distribution:
    repository: ./testdata/gittest/distribution
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
  myStack:
    component: "_"
`
	checkSelfStack(t, descContent)
}

func checkSelfStack(t *testing.T, descContent string) {
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

	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repDesc := tester.createRep(mainPath)

	repDist.writeCommit(t, "ekara.yaml", distContent)
	repComp1.writeCommit(t, "ekara.yaml", "")

	repDesc.writeCommit(t, "ekara.yaml", descContent)
	// write the compose/playbook content into the descriptor component
	repDesc.writeCommit(t, "docker_compose.yml", "docker compose content")

	err := tester.initEngine()
	assert.Nil(t, err)
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator and provider
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId, "comp1")

	// Chect that the enviroment has one self contained stack
	if assert.Equal(t, 1, len(env.Stacks)) {

		cm := c.Ekara().ComponentManager()
		assert.NotNil(t, cm)

		stack, ok := env.Stacks["myStack"]
		if assert.True(t, ok) {
			//Check that the self contained stack has been well built
			assert.Equal(t, "myStack", stack.Name)
			stackC, err := stack.Component()
			assert.Nil(t, err)
			assert.NotNil(t, stackC)
			assert.Equal(t, model.MainComponentId, stackC.Id)

			// Check that the stack is usable and returns the environent as component
			usableStack, err := cm.Use(stack)
			defer usableStack.Release()
			assert.Nil(t, err)
			assert.NotNil(t, usableStack)
			assert.False(t, usableStack.Templated())
			// Check that the stacks contains the compose/playbook file
			checkFile(t, usableStack, "docker_compose.yml", "docker compose content")
		}
	}
}

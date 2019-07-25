package engine

import (
	"testing"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestParentSelfStackNoComponent(t *testing.T) {

	distContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
stacks:
  myStack:
    params:
      myStack_param_key1: myStack_param_key1_value
      myStack_param_key2: myStack_param_key1_value
`
	checkSelfStackParent(t, distContent)
}

func TestParentSelfStackLowDash(t *testing.T) {
	distContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
stacks:
  myStack:
    component: "_"
    params:
      myStack_param_key1: myStack_param_key1_value
      myStack_param_key2: myStack_param_key1_value
`
	checkSelfStackParent(t, distContent)
}

func checkSelfStackParent(t *testing.T, distContent string) {

	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/parent")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp1")
	repDesc := tester.createRep(mainPath)

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
stacks:
  myStack:
    component: __ekara__1  
    params:
      myStack_param_key2: myStack_param_key2_value_overwritten
      myStack_param_key3: myStack_param_key3_value
`
	repDesc.writeCommit(t, "ekara.yaml", descContent)
	// write the compose/playbook content into the parent component
	repDist.writeCommit(t, "docker_compose.yml", "docker compose content")

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator and provider
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", "comp1")

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
			assert.Equal(t, model.EkaraComponentId+"1", stackC.Id)

			// Check that the stack is usable and returns the environent as component
			usableStack, err := cm.Use(stack)
			defer usableStack.Release()
			assert.Nil(t, err)
			assert.NotNil(t, usableStack)
			assert.False(t, usableStack.Templated())
			// Check that the stacks contains the compose/playbook file
			checkFile(t, usableStack, "docker_compose.yml", "docker compose content")

			//check the stack parameters inheritence
			if assert.Equal(t, 3, len(stack.Parameters)) {
				assert.Contains(t, stack.Parameters, "myStack_param_key1", "myStack_param_key1_value")
				assert.Contains(t, stack.Parameters, "myStack_param_key2", "myStack_param_key2_value_overwritten")
				assert.Contains(t, stack.Parameters, "myStack_param_key3", "myStack_param_key3_value")
			}
		}
	}
}

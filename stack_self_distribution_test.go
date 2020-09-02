package engine

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/engine/model"

	"github.com/stretchr/testify/assert"
)

func TestParentSelfStackNoComponent(t *testing.T) {
	parentContent := `
ekara:
  components:
    comp1:
      repository: comp1
stacks:
  myStack:
    params:
      myStack_param_key1: myStack_param_key1_value
      myStack_param_key2: myStack_param_key1_value
`
	checkSelfStackParent(t, parentContent)
}

func TestParentSelfStackLowDash(t *testing.T) {
	parentContent := `
ekara:
  components:
    comp1:
      repository: comp1
stacks:
  myStack:
    component: "_"
    params:
      myStack_param_key1: myStack_param_key1_value
      myStack_param_key2: myStack_param_key1_value
`
	checkSelfStackParent(t, parentContent)
}

func checkSelfStackParent(t *testing.T, parentContent string) {
	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	repParent := tester.CreateDir("parent")
	tester.CreateDirEmptyDesc("comp1")
	repDesc := tester.CreateDir("descriptor")

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
stacks:
  myStack:
    component: __main__parent  
    params:
      myStack_param_key2: myStack_param_key2_value_overwritten
      myStack_param_key3: myStack_param_key3_value
`
	repDesc.WriteCommit("ekara.yaml", descContent)
	// write the compose/playbook content into the parent component
	repParent.WriteCommit("docker_compose.yml", "docker compose content")

	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator and provider
	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1")

	// Check that the environment has one self contained stack
	if assert.Len(t, env.Stacks, 1) {
		stack, ok := env.Stacks["myStack"]
		if assert.True(t, ok) {
			//Check that the self contained stack has been well built
			assert.Equal(t, "myStack", stack.Name)
			stackC, err := stack.Component(tester.Model())
			assert.Nil(t, err)
			assert.NotNil(t, stackC)
			assert.Equal(t, model.MainComponentId+model.ParentComponentSuffix, stackC.ComponentId())

			// Check that the stack is usable and returns the environent as component
			usableStack, err := tester.ComponentManager().Use(stack, tester.TemplateContext())
			defer usableStack.Release()
			assert.Nil(t, err)
			assert.NotNil(t, usableStack)
			assert.False(t, usableStack.Templated())
			// Check that the stacks contains the compose/playbook file
			tester.AssertFileContent(usableStack, "docker_compose.yml", "docker compose content")

			//check the stack parameters inheritance
			if assert.Len(t, stack.Parameters(), 3) {
				assert.Contains(t, stack.Parameters(), "myStack_param_key1", "myStack_param_key1_value")
				assert.Contains(t, stack.Parameters(), "myStack_param_key2", "myStack_param_key2_value_overwritten")
				assert.Contains(t, stack.Parameters(), "myStack_param_key3", "myStack_param_key3_value")
			}
		}
	}
}

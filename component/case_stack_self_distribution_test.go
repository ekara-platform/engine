package component

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestParentSelfStackNoComponent(t *testing.T) {

	parentContent := `
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
	checkSelfStackParent(t, parentContent)
}

func TestParentSelfStackLowDash(t *testing.T) {
	parentContent := `
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
	checkSelfStackParent(t, parentContent)
}

func checkSelfStackParent(t *testing.T, parentContent string) {

	mainPath := "./testdata/gittest/descriptor"
	c := util.CreateMockLaunchContext(mainPath, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	repParent := tester.CreateRep("./testdata/gittest/parent")
	tester.CreateRepDefaultDescriptor("./testdata/gittest/comp1")
	repDesc := tester.CreateRep(mainPath)

	repParent.WriteCommit("ekara.yaml", parentContent)

	descContent := `
name: ekaraDemoVar
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
	repDesc.WriteCommit("ekara.yaml", descContent)
	// write the compose/playbook content into the parent component
	repParent.WriteCommit("docker_compose.yml", "docker compose content")

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator and provider
	tester.AssertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", "comp1")

	// Chect that the environment has one self contained stack
	if assert.Len(t, env.Stacks, 1) {

		cF := tester.cF
		assert.NotNil(t, cF)

		stack, ok := env.Stacks["myStack"]
		if assert.True(t, ok) {
			//Check that the self contained stack has been well built
			assert.Equal(t, "myStack", stack.Name)
			stackC, err := stack.Component()
			assert.Nil(t, err)
			assert.NotNil(t, stackC)
			assert.Equal(t, model.EkaraComponentId+"1", stackC.Id)

			// Check that the stack is usable and returns the environent as component
			usableStack, err := cF.Use(stack, *tester.tplC)
			defer usableStack.Release()
			assert.Nil(t, err)
			assert.NotNil(t, usableStack)
			assert.False(t, usableStack.Templated())
			// Check that the stacks contains the compose/playbook file
			tester.CheckFile(usableStack, "docker_compose.yml", "docker compose content")

			//check the stack parameters inheritance
			if assert.Len(t, stack.Parameters, 3) {
				assert.Contains(t, stack.Parameters, "myStack_param_key1", "myStack_param_key1_value")
				assert.Contains(t, stack.Parameters, "myStack_param_key2", "myStack_param_key2_value_overwritten")
				assert.Contains(t, stack.Parameters, "myStack_param_key3", "myStack_param_key3_value")
			}
		}
	}
}

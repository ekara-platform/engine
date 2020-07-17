package engine

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/engine/model"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentSelfStackNoComponent(t *testing.T) {

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
`
	checkSelfStack(t, descContent)
}

func TestEnvironmentSelfStackLowDash(t *testing.T) {

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
    component: "_"
`
	checkSelfStack(t, descContent)
}

func checkSelfStack(t *testing.T, descContent string) {
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
	repDesc := tester.CreateDir(mainPath)

	repParent.WriteCommit("ekara.yaml", parentContent)

	repDesc.WriteCommit("ekara.yaml", descContent)
	// write the compose/playbook content into the descriptor component
	repDesc.WriteCommit("docker_compose.yml", "docker compose content")

	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator and provider
	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1")

	// Chect that the environment has one self contained stack
	if assert.Len(t, env.Stacks, 1) {
		stack, ok := env.Stacks["myStack"]
		if assert.True(t, ok) {
			//Check that the self contained stack has been well built
			assert.Equal(t, "myStack", stack.Name)
			stackC, err := stack.Component(tester.Model())
			assert.Nil(t, err)
			assert.NotNil(t, stackC)
			assert.Equal(t, model.MainComponentId, stackC.ComponentId())

			// Check that the stack is usable and returns the environent as component
			usableStack, err := tester.ComponentManager().Use(stack, tester.TemplateContext())
			defer usableStack.Release()
			assert.Nil(t, err)
			assert.NotNil(t, usableStack)
			assert.False(t, usableStack.Templated())
			// Check that the stacks contains the compose/playbook file
			tester.AssertFileContent(usableStack, "docker_compose.yml", "docker compose content")
		}
	}
}

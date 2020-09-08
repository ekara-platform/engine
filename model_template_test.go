package engine

import (
	"github.com/ekara-platform/engine/model"
	"github.com/ekara-platform/engine/util"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"testing"
)

func TestTemplateOnReadOnlyModel(t *testing.T) {
	mainPath := "descriptor"
	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	tester.CreateDirEmptyDesc("parent")
	tester.CreateDirEmptyDesc("comp1")
	repStack := tester.CreateDir("stack")
	repDesc := tester.CreateDir(mainPath)

	stackContent := `
ekara:
 templates:
   - "/template.txt"
`
	repStack.WriteCommit("ekara.yaml", stackContent)
	repStack.WriteCommit("template.txt", "{{.Model.QName}}")

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
 parent:
   repository: parent
 components:
   comp1:
     repository: comp1
   stack1:
     repository: stack
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
 stack1:
   component: stack1
`

	repDesc.WriteCommit("ekara.yaml", descContent)

	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)

	assert.Equal(t, "ekaraDemoVar", env.QName.Name)
	assert.Equal(t, "dev", env.QName.Qualifier)

	tEnvironment := tester.TemplateContext().(*model.TemplateContext).Model
	assert.NotNil(t, tEnvironment)
	assert.Equal(t, "ekaraDemoVar_dev", tEnvironment.QName.String())

	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1", "stack1")
	if assert.Equal(t, 1, len(env.Stacks)) {
		stack, ok := env.Stacks["stack1"]
		if assert.True(t, ok) {
			c, err := stack.Component(env)
			assert.Nil(t, err)
			isTemplated, templateGlobs := c.GetTemplates()
			assert.True(t, isTemplated)
			assert.Equal(t, 1, len(templateGlobs))
			usableStack, err := tester.ComponentManager().Use(stack, tester.TemplateContext())
			assert.Nil(t, err)
			defer usableStack.Release()
			assert.True(t, usableStack.Templated())

			b, err := ioutil.ReadFile(path.Join(usableStack.RootPath(), "template.txt"))
			assert.Nil(t, err)
			assert.Equal(t, "ekaraDemoVar_dev", string(b))
		}
	}
}

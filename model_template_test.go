package engine

import (
	"testing"
)

func TestTemplateOnReadOnlyModel(t *testing.T) {
	/*
	   	mainPath := "./testdata/gittest/descriptor"
	   	c := &util.MockLaunchContext{locationContent: mainPath, templateContext: tester.context.engine.Context().data}
	   	tester := gitTester(t, c)
	   	defer tester.clean()

	   	tester.createRepDefaultDescriptor(t, "./testdata/gittest/parent")
	   	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp1")
	   	repStack := tester.createRep("./testdata/gittest/stack")
	   	repDesc := tester.createRep(mainPath)

	   	stackContent := `
	   templates:
	     - "/template.txt"
	   `
	   	repStack.writeCommit(t, "ekara.yaml", stackContent)
	   	repStack.writeCommit(t, "template.txt", "{{.Model.QualifiedName}}")

	   	descContent := `
	   name: ekaraDemoVar
	   qualifier: dev

	   ekara:
	     parent:
	       repository: ./testdata/gittest/parent
	     components:
	       comp1:
	         repository: ./testdata/gittest/comp1
	       stack1:
	         repository: ./testdata/gittest/stack
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

	   	repDesc.writeCommit(t, "ekara.yaml", descContent)

	   	err := tester.initEngine()
	   	assert.Nil(t, err)
	   	env := tester.env()
	   	assert.NotNil(t, env)

	   	tEnvironment := c.TemplateContext().Model
	   	assert.Equal(t, "ekaraDemoVar", env.Name)
	   	assert.Equal(t, "dev", env.Qualifier)

	   	assert.NotNil(t, tEnvironment)
	   	assert.Equal(t, "ekaraDemoVar", tEnvironment.Name())
	   	assert.Equal(t, "dev", tEnvironment.Qualifier())
	   	assert.Equal(t, "ekaraDemoVar_dev", tEnvironment.QualifiedName())

	   	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", "comp1", "stack1")
	   	if assert.Equal(t, 1, len(env.Stacks)) {
	   		stack, ok := env.Stacks["stack1"]
	   		if assert.True(t, ok) {
	   			c, err := stack.Component()
	   			assert.Nil(t, err)
	   			assert.Equal(t, 1, len(c.Templates.Content))
	   			cm := tester.context.Ekara().ComponentManager()
	   			assert.Equal(t, 4, tester.countComponent())
	   			usableStack, err := cm.Use(stack)
	   			assert.Nil(t, err)
	   			assert.Equal(t, 5, tester.countComponent())
	   			defer usableStack.Release()
	   			assert.True(t, usableStack.Templated())

	   			b, err := ioutil.ReadFile(path.Join(usableStack.RootPath(), "template.txt"))
	   			assert.Nil(t, err)
	   			assert.Equal(t, "ekaraDemoVar_dev", string(b))
	   		}
	   	}
	*/
}

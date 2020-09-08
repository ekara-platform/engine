package util

import (
	"github.com/GroupePSA/componentizer"
	"github.com/ekara-platform/engine/model"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

type EkaraComponentTester struct {
	*componentizer.ComponentTester
}

func CreateComponentTester(t *testing.T, extVars model.Parameters) EkaraComponentTester {
	return EkaraComponentTester{ComponentTester: componentizer.CreateComponentTester(componentizer.TestContext{
		T:              t,
		Logger:         log.New(os.Stdout, "TEST: ", log.LstdFlags),
		Directory:      "./testdata",
		DescriptorName: model.DefaultDescriptorName,
	}, model.CreateTemplateContext(extVars))}
}

func (t EkaraComponentTester) Init(repo componentizer.Repository) {
	err := t.ComponentTester.Init(model.CreateComponent(model.MainComponentId, repo))
	if err != nil {
		assert.Nil(t.T(), err, "Init error: %s", err.Error())
	}
	t.ComponentTester.TemplateContext().(*model.TemplateContext).Model = t.ComponentTester.Model().(model.Environment)
}

func (t EkaraComponentTester) Env() model.Environment {
	return t.ComponentTester.Model().(model.Environment)
}

func (t EkaraComponentTester) AssertParam(p model.Parameters, key, value string) {
	assert.Contains(t.T(), p, key)
	assert.Equal(t.T(), value, p[key])
}

func (t EkaraComponentTester) AssertEnvVar(p model.EnvVars, key, value string) {
	assert.Contains(t.T(), p, key)
	assert.Equal(t.T(), value, p[key])
}

func CreateFakeComponent(id string) componentizer.Component {
	return model.CreateComponent(id, componentizer.Repository{})
}

// CheckStack asserts than the environment contains the stack, then check that the stack's component
// is usable and finally check that the stack constraints the provide compose content .
func (t EkaraComponentTester) CheckStack(holder, stackName, composeContent string) {
	stack, ok := t.Model().(model.Environment).Stacks[stackName]
	if assert.True(t.T(), ok) {
		//Check that the self contained stack has been well built
		assert.Equal(t.T(), stackName, stack.Name)
		stackC, err := stack.Component(t.Model())
		assert.Nil(t.T(), err)
		assert.NotNil(t.T(), stackC)
		assert.Equal(t.T(), holder, stackC.ComponentId())

		// Check that the stack is usable and returns the correct component
		usableStack, err := t.ComponentManager().Use(stack, t.TemplateContext())
		assert.Nil(t.T(), err)
		defer usableStack.Release()
		assert.NotNil(t.T(), usableStack)
		assert.False(t.T(), usableStack.Templated())
		// Check that the stacks contains the compose/playbook file
		t.AssertFileContent(usableStack, "docker_compose.yml", composeContent)
	}
}

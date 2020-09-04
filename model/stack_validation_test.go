package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading an environment without stacks.
//
// The validation must complain only about the missing stacks
//- Warning: no stack specified @stacks
//
func TestValidationNoStacks(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "stacks", true)
	assert.False(t, vErrs.HasErrors())
	assert.True(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Warning, "no stack specified", "stacks"))
}

// Test loading an environment referencing an unknown stack.
//
// The validation must complain only about the reference on unknown stack
//
//- Error: reference to unknown component: dummy @stacks.monitoring.component
//
func TestValidationUnknownStack(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml("./testdata/yaml/grammar/unknown_stack.yaml", &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{Id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "no such component: dummy", "stacks.monitoring.component"))

}

// Test loading an environment referencing stack which depends on an unknown one.
//
// The validation must complain only about the dependency on unknown stack
//
//- Error: reference to unknown stack dependency: dummy @stacks.monitoring.dependencies.dummy
//
func TestValidationUnknownDependencies(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml("./testdata/yaml/grammar/stack_unknown_depends_on.yaml", &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{Id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "no such stack: dummy", "stacks.monitoring.dependencies"))
}

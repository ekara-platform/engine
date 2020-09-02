package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading an empty yaml file.
//
// The validation must complain about all root elements missing
//- Error: empty environment name @name
//- Error: empty component reference @orchestrator
//- Error: no provider specified @providers
//- Error: no node specified @nodes
//- Warning: no stack specified @stacks
//
// There is no message about a missing ekera platform because it has been defaulted
func TestValidationNoContent(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "content", false)
	assert.True(t, vErrs.HasErrors())
	assert.True(t, vErrs.HasWarnings())
	assert.Equal(t, 5, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "empty environment name", "name"))
	assert.True(t, vErrs.contains(Error, "empty component reference", "orchestrator.component"))
	assert.True(t, vErrs.contains(Error, "no provider specified", "providers"))
	assert.True(t, vErrs.contains(Error, "no node specified", "nodes"))
	assert.True(t, vErrs.contains(Warning, "no stack specified", "stacks"))
}

func testEmptyContent(t *testing.T, name string, onlyWarning bool) (ValidationErrors, Environment) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml(fmt.Sprintf("./testdata/yaml/grammar/no_%s.yaml", name), &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	vErrs := testValidate(t, env, onlyWarning)
	return vErrs, env
}

func testValidate(t *testing.T, env Environment, onlyWarning bool) ValidationErrors {
	vErrs := env.Validate()
	if onlyWarning {
		assert.True(t, vErrs.HasWarnings())
	} else {
		assert.True(t, vErrs.HasErrors())
	}
	return vErrs
}

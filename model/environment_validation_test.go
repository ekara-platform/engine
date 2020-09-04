package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading an environment without name.
//
// The validation must complain only about the missing name
//- Error: empty environment name @name
//
func TestValidationNoEnvironmentName(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "environment_name", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "empty environment name", "name"))
}

// Test loading an environment with an invalid name
//
// The validation must complain only about the invalid name
//- Error: the environment name or the qualifier contains a non alphanumeric character @name|qualifier
//
func TestValidateNoValidName(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml("./testdata/yaml/grammar/no_valid_name.yaml", &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{Id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "the environment name or the qualifier contains a non alphanumeric character", "name|qualifier"))
}

// Test loading an environment with an invalid qualifier
//
// The validation must complain only about the invalid qualifier
//- Error: the environment name or the qualifier contains a non alphanumeric character @name|qualifier
//
func TestValidateNoValidQualifier(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml("./testdata/yaml/grammar/no_valid_qualifier.yaml", &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{Id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "the environment name or the qualifier contains a non alphanumeric character", "name|qualifier"))
}

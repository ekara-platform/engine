package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading an environment without providers.
//
// The validation must complain only about the missing providers and the reference
// to a missing provider into the node set specification
//
//- Error: no provider specified @providers
//- Error: reference to unknown provider: aws @nodes.managers.provider
//
func TestValidationNoProviders(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "providers", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 2, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "no provider specified", "providers"))
	assert.True(t, vErrs.contains(Error, "no such provider: aws", "nodes.managers.provider"))
}

// Test loading an nodeset referencing an unknown provider.
//
// The validation must complain only about the reference on unknown provider
//
//- Error: reference to unknown provider: dummy @nodes.managers.provider
//
func TestValidationNodesUnknownProvider(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml("./testdata/yaml/grammar/nodes_unknown_provider.yaml", &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{Id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "no such provider: dummy", "nodes.managers.provider"))
}

// Test loading an node set without a reference on a provider.
//
// The validation must complain only about the missing provider reference
//- Error: empty provider reference @nodes.managers.provider
//
func TestValidationNoNodesProvider(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "nodes_provider", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "empty provider reference", "nodes.managers.provider"))
}

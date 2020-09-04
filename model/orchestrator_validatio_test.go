package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading an environment without orchestator.
//
// The validation must complain only about the missing orchestator
//- Error: empty component reference @orchestrator
//
func TestValidationNoOrchestrator(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "orchestrator", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "empty component reference", "orchestrator.component"))
}

// Test loading an environment referencing an unknown orchestrator.
//
// The validation must complain only about the reference on unknown orchestrator
//
//- Error: reference to unknown component: dummy @orchestrator
//
func TestValidationUnknownOrchestrator(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml("./testdata/yaml/grammar/unknown_orchestrator.yaml", &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{Id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "no such component: dummy", "orchestrator.component"))
}

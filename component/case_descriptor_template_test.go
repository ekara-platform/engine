package component

import (
	"testing"

	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	path := "./testdata/template/template-params.yaml"
	params, err := model.ParseParameters(path)
	assert.Nil(t, err)
	assert.NotNil(t, params)
	vars := model.CreateTemplateContext(params)

	url, err := model.CreateUrl("./testdata/template/descriptor.yaml")
	assert.Nil(t, err)

	yamlEnv, err := model.ParseYamlDescriptor(url, vars)
	assert.Nil(t, err)
	// Parsing the descriptor
	env, err := model.CreateEnvironment("", yamlEnv, model.MainComponentId)
	assert.Nil(t, err)

	// Parsing the descriptor
	assert.NotNil(t, env)
	assert.Equal(t, 2, len(env.Tasks))
	ta := env.Tasks["testhook_post"]
	assert.Equal(t, 2, len(ta.Parameters))
	val, ok := ta.Parameters["param1"]
	assert.True(t, ok)
	assert.Equal(t, "key2_value-key4_value", val)
	val, ok = ta.Parameters["param2"]
	assert.True(t, ok)
	assert.Equal(t, "key4_value", val)
}

func TestTemplateNoDot(t *testing.T) {

	path := "./testdata/template/template-params.yaml"
	params, err := model.ParseParameters(path)
	assert.Nil(t, err)
	assert.NotNil(t, params)
	vars := model.CreateTemplateContext(params)

	url, err := model.CreateUrl("./testdata/template/descriptor_no_dot.yaml")
	assert.Nil(t, err)

	// Parsing the descriptor
	_, err = model.ParseYamlDescriptor(url, vars)
	assert.NotNil(t, err)

}
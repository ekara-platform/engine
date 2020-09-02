package engine

import (
	"github.com/ekara-platform/engine/util"
	"testing"

	"github.com/ekara-platform/engine/model"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	path := "./testdata/template/params.yaml"
	params, err := model.ParseParameters(path)
	assert.Nil(t, err)
	assert.NotNil(t, params)
	m, err := util.CreateFakeComponent(model.MainComponentId).ParseModel("testdata/template/nominal", model.CreateTemplateContext(params))
	env := m.(model.Environment)

	assert.NotNil(t, env)
	assert.Equal(t, 2, len(env.Tasks))
	ta := env.Tasks["testhook_post"]
	assert.Equal(t, 2, len(ta.Parameters()))
	val, ok := ta.Parameters()["param1"]
	assert.True(t, ok)
	assert.Equal(t, "key2_value-key4_value", val)
	val, ok = ta.Parameters()["param2"]
	assert.True(t, ok)
	assert.Equal(t, "key4_value", val)
}

func TestTemplateNoDot(t *testing.T) {
	path := "./testdata/template/params.yaml"
	params, err := model.ParseParameters(path)
	assert.Nil(t, err)
	assert.NotNil(t, params)
	_, err = util.CreateFakeComponent(model.MainComponentId).ParseModel("testdata/template/no-dot", model.CreateTemplateContext(params))
	assert.NotNil(t, err)
}

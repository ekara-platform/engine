package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericNode(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml("./testdata/yaml/overwritten/ekara.yaml", &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	assert.Nil(t, e)
	if assert.Equal(t, len(env.NodeSets), 1) {
		n := env.NodeSets["managers"]
		p, e := n.Provider.Resolve(env)
		assert.Nil(t, e)
		if val, ok := p.Parameters()["generic_param1"]; ok {
			assert.Equal(t, val, "new_generic_param1")
		} else {
			assert.Fail(t, "missing generic param")
		}

		if val, ok := p.envVars["generic_env1"]; ok {
			assert.Equal(t, val, "new_generic_env1")
		} else {
			assert.Fail(t, "missing generic env var")
		}

		assert.Equal(t, p.proxy.NoProxy, "overwritten_aws_no_proxy")
		assert.Equal(t, p.proxy.Https, "generic_https_proxy")
		assert.Equal(t, p.proxy.Http, "aws_http_proxy")
	}
}

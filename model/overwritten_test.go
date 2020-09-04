package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOverwrittenProviderParam(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml("./testdata/yaml/overwritten/ekara.yaml", &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{Id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	aws := env.Providers["aws"]
	assert.NotNil(t, aws)
	assert.NotNil(t, aws.Parameters())
	assert.Equal(t, 2, len(aws.Parameters()))
	assert.Equal(t, "initial_param1", aws.Parameters()["param1"])
	assert.Equal(t, "initial_param3", aws.Parameters()["param3"])

	managers := env.NodeSets["managers"]
	assert.NotNil(t, managers)
	managersProvider, e := managers.Provider.Resolve(env)
	assert.Nil(t, e)
	params := managersProvider.Parameters()
	assert.NotNil(t, params)
	assert.Equal(t, 4, len(params))
	assert.Equal(t, "new_generic_param1", params["generic_param1"])
	assert.Equal(t, "overwritten_param1", params["param1"])
	assert.Equal(t, "new_param2", params["param2"])
	assert.Equal(t, "initial_param3", params["param3"])
}

func TestOverwrittenProviderEnv(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml("./testdata/yaml/overwritten/ekara.yaml", &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{Id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	aws := env.Providers["aws"]
	assert.NotNil(t, aws)
	assert.NotNil(t, aws.envVars)
	assert.Equal(t, 2, len(aws.envVars))
	assert.Equal(t, "initial_env1", aws.envVars["env1"])
	assert.Equal(t, "initial_env3", aws.envVars["env3"])

	managers := env.NodeSets["managers"]
	assert.NotNil(t, managers)
	managersProvider, e := managers.Provider.Resolve(env)
	assert.Nil(t, e)
	envs := managersProvider.envVars
	assert.NotNil(t, envs)
	assert.Equal(t, 4, len(envs))
	assert.Equal(t, "overwritten_env1", envs["env1"])
	assert.Equal(t, "new_env2", envs["env2"])
	assert.Equal(t, "initial_env3", envs["env3"])
}

func TestOverwrittenProviderProxy(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml("./testdata/yaml/overwritten/ekara.yaml", &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{Id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	aws := env.Providers["aws"]
	assert.NotNil(t, aws)
	assert.NotNil(t, aws.proxy)
	assert.Equal(t, "", aws.proxy.Https)
	assert.Equal(t, "aws_http_proxy", aws.proxy.Http)
	assert.Equal(t, "aws_no_proxy", aws.proxy.NoProxy)

	managers := env.NodeSets["managers"]
	assert.NotNil(t, managers)
	managersProvider, e := managers.Provider.Resolve(env)
	assert.Nil(t, e)
	pr := managersProvider.proxy
	assert.NotNil(t, pr)
	assert.Equal(t, "aws_http_proxy", pr.Http)
	assert.Equal(t, "generic_https_proxy", pr.Https)
	assert.Equal(t, "overwritten_aws_no_proxy", pr.NoProxy)
}

// TODO Add test for TaskRef ans Task and stack

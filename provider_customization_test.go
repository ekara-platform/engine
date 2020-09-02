package engine

import (
	"github.com/ekara-platform/engine/util"
	"testing"

	"github.com/ekara-platform/engine/model"
	"github.com/stretchr/testify/assert"
)

func TestMergeProvider(t *testing.T) {

	parent2Content := `
ekara:
  components:
    comp2:
      repository: comp2
providers:
  p1:
    component: comp2
    params:
      param_key1: parent2_param_key1_value
      param_key2: parent2_param_key2_value
      param_key3: parent2_param_key3_value
      param_key4: parent2_param_key4_value
    env:
      env_key1: parent2_env_key1_value
      env_key2: parent2_env_key2_value
      env_key3: parent2_env_key3_value
      env_key4: parent2_env_key4_value
    proxy:
      http_proxy: http_proxy_parent2
      #https_proxy: https_proxy_parent2
      #no_proxy: no_proxy_parent2
`

	comp1Content := `
providers:
  p1:
    params:
      param_key2: comp1_param_key2_value
      param_key3: comp1_param_key3_value
      param_key4: comp1_param_key4_value	  
    env:
      env_key2: comp1_env_key2_value
      env_key3: comp1_env_key3_value
      env_key4: comp1_env_key4_value
    proxy:
      https_proxy: https_proxy_comp1
      #no_proxy: no_proxy_comp1
`

	parent1Content := `
ekara:
  parent:
    repository: parent2
  components:
    comp1:
      repository: comp1

providers:
  p1:
    params:
      param_key3: parent1_param_key3_value
      param_key4: parent1_param_key4_value	  
    env:
      env_key3: parent1_env_key3_value
      env_key4: parent1_env_key4_value	  
    proxy:
      no_proxy: no_proxy_parent1
`

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: parent1
# Following content just to force the download of comp1 and comp2
orchestrator:
  component: comp1
providers:
  p1:
    component: comp2
    params:
      param_key4: desc_param_key4_value	  
    env:
      env_key4: desc_env_key4_value	    
nodes:
  node1:
    instances: 1
    provider:
      name: p1
`

	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	repParent1 := tester.CreateDir("parent1")
	repParent2 := tester.CreateDir("parent2")
	repComp1 := tester.CreateDir("comp1")
	tester.CreateDirEmptyDesc("comp2")
	repDesc := tester.CreateDir("descriptor")

	repComp1.WriteCommit("ekara.yaml", comp1Content)
	repParent1.WriteCommit("ekara.yaml", parent1Content)
	repParent2.WriteCommit("ekara.yaml", parent2Content)
	repDesc.WriteCommit("ekara.yaml", descContent)

	tester.Init(repDesc.AsRepository("master"))
	tester.AssertComponentAvailable(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, model.MainComponentId+model.ParentComponentSuffix+model.ParentComponentSuffix, "comp1", "comp2")

	env := tester.Env()
	assert.NotNil(t, env)

	assert.Equal(t, len(env.Providers), 1)
	p, ok := env.Providers["p1"]
	if assert.True(t, ok) {
		params := p.Parameters()
		if assert.Equal(t, 4, len(params)) {
			tester.AssertParam(params, "param_key1", "parent2_param_key1_value")
			tester.AssertParam(params, "param_key2", "comp1_param_key2_value")
			tester.AssertParam(params, "param_key3", "parent1_param_key3_value")
			tester.AssertParam(params, "param_key4", "desc_param_key4_value")
		}
		envVars := p.EnvVars()
		if assert.Equal(t, 4, len(envVars)) {
			tester.AssertEnvVar(envVars, "env_key1", "parent2_env_key1_value")
			tester.AssertEnvVar(envVars, "env_key2", "comp1_env_key2_value")
			tester.AssertEnvVar(envVars, "env_key3", "parent1_env_key3_value")
			tester.AssertEnvVar(envVars, "env_key4", "desc_env_key4_value")

		}
		assert.Equal(t, p.Proxy().Http, "http_proxy_parent2")
		assert.Equal(t, p.Proxy().Https, "https_proxy_comp1")
		assert.Equal(t, p.Proxy().NoProxy, "no_proxy_parent1")
	}
}

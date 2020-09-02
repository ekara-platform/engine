package engine

import (
	"github.com/ekara-platform/engine/util"
	"testing"

	"github.com/ekara-platform/engine/model"
	"github.com/stretchr/testify/assert"
)

func TestMergeOrchestrator(t *testing.T) {

	parent2Content := `
ekara:
  components:
    comp2:
      repository: comp2
orchestrator:
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
`

	comp1Content := `
orchestrator:
  params:
    param_key2: comp1_param_key2_value
    param_key3: comp1_param_key3_value
    param_key4: comp1_param_key4_value	  
  env:
    env_key2: comp1_env_key2_value
    env_key3: comp1_env_key3_value
    env_key4: comp1_env_key4_value
`

	parent1Content := `
ekara:
  parent:
    repository: parent2
  components:
    comp1:
      repository: comp1

orchestrator:
  params:
    param_key3: parent1_param_key3_value
    param_key4: parent1_param_key4_value	  
  env:
    env_key3: parent1_env_key3_value
    env_key4: parent1_env_key4_value	  
`

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: parent1
# Following content just to force the download of comp1 and comp2
orchestrator:
  params:
    param_key4: desc_param_key4_value	  
  env:
    env_key4: desc_env_key4_value	    
providers:
  p1:
    component: comp1
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
	oParams := env.Orchestrator.Parameters()
	if assert.Equal(t, 4, len(oParams)) {
		tester.AssertParam(oParams, "param_key1", "parent2_param_key1_value")
		tester.AssertParam(oParams, "param_key2", "comp1_param_key2_value")
		tester.AssertParam(oParams, "param_key3", "parent1_param_key3_value")
		tester.AssertParam(oParams, "param_key4", "desc_param_key4_value")
	}

	oEnvVars := env.Orchestrator.EnvVars()
	if assert.Equal(t, 4, len(oEnvVars)) {
		tester.AssertEnvVar(oEnvVars, "env_key1", "parent2_env_key1_value")
		tester.AssertEnvVar(oEnvVars, "env_key2", "comp1_env_key2_value")
		tester.AssertEnvVar(oEnvVars, "env_key3", "parent1_env_key3_value")
		tester.AssertEnvVar(oEnvVars, "env_key4", "desc_env_key4_value")
	}
}

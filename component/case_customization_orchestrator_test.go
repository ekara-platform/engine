package component

import (
	"github.com/ekara-platform/engine/util"
	"testing"

	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"
)

func TestMergeOrchestrator(t *testing.T) {

	parent2Content := `
ekara:
  components:
    comp2:
      repository: ./testdata/gittest/comp2
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
    repository: ./testdata/gittest/parent2
  components:
    comp1:
      repository: ./testdata/gittest/comp1

orchestrator:
  params:
    param_key3: parent1_param_key3_value
    param_key4: parent1_param_key4_value	  
  env:
    env_key3: parent1_env_key3_value
    env_key4: parent1_env_key4_value	  
`

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent1
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

	mainPath := "./testdata/gittest/descriptor"
	c := util.CreateMockLaunchContext(mainPath, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	repParent1 := tester.CreateRep("./testdata/gittest/parent1")
	repParent2 := tester.CreateRep("./testdata/gittest/parent2")
	repComp1 := tester.CreateRep("./testdata/gittest/comp1")
	tester.CreateRepDefaultDescriptor("./testdata/gittest/comp2")
	repDesc := tester.CreateRep(mainPath)

	repComp1.WriteCommit("ekara.yaml", comp1Content)
	repParent1.WriteCommit("ekara.yaml", parent1Content)
	repParent2.WriteCommit("ekara.yaml", parent2Content)
	repDesc.WriteCommit("ekara.yaml", descContent)

	err := tester.Init()
	assert.Nil(t, err)

	tester.AssertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", model.EkaraComponentId+"2", "comp1", "comp2")

	env := tester.Env()
	assert.NotNil(t, env)

	if assert.Equal(t, len(env.Orchestrator.Parameters), 4) {
		tester.CheckSpecificParameter(env.Orchestrator.Parameters, "param_key1", "parent2_param_key1_value")
		tester.CheckSpecificParameter(env.Orchestrator.Parameters, "param_key2", "comp1_param_key2_value")
		tester.CheckSpecificParameter(env.Orchestrator.Parameters, "param_key3", "parent1_param_key3_value")
		tester.CheckSpecificParameter(env.Orchestrator.Parameters, "param_key4", "desc_param_key4_value")
	}

	if assert.Equal(t, len(env.Orchestrator.EnvVars), 4) {
		tester.CheckSpecificEnvVar(env.Orchestrator.EnvVars, "env_key1", "parent2_env_key1_value")
		tester.CheckSpecificEnvVar(env.Orchestrator.EnvVars, "env_key2", "comp1_env_key2_value")
		tester.CheckSpecificEnvVar(env.Orchestrator.EnvVars, "env_key3", "parent1_env_key3_value")
		tester.CheckSpecificEnvVar(env.Orchestrator.EnvVars, "env_key4", "desc_env_key4_value")
	}
}

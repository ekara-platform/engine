package engine

import (
	"testing"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestTemplateOnMainVars(t *testing.T) {

	p, _ := model.CreateParameters(map[string]interface{}{
		"value1": map[interface{}]interface{}{
			"from": map[interface{}]interface{}{
				"cli": "value1.from.cli_value",
			},
		},
		"value2": "value2.from.cli_value",
	})
	mainPath := "./testdata/gittest/descriptor"
	tc := model.CreateContext(p)
	c := &MockLaunchContext{locationContent: mainPath, templateContext: tc}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repDesc := tester.createRep(mainPath)

	repComp1.writeCommit(t, "ekara.yaml", "")
	repDist.writeCommit(t, "ekara.yaml", "")

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  distribution:
    repository: ./testdata/gittest/distribution	
  components:
    comp1:
      repository: ./testdata/gittest/comp1	
vars:
  key1_descriptor: val1_descriptor
  key2_descriptor: "{{ .Vars.value1.from.cli }}"
orchestrator:
  component: comp1
providers:
  p1:
    component: comp1
    params:
      param1: {{ .Vars.key1_descriptor }}
      param2: {{ .Vars.key2_descriptor }}
      param3: {{ .Vars.value2 }} 
nodes:
  node1:
    instances: 1
    provider:
      name: p1
`

	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	// Check if the descriptor has been templated
	assert.Equal(t, len(env.Vars), 2)
	cp(t, env.Vars, "key1_descriptor", "val1_descriptor")
	cp(t, env.Vars, "key2_descriptor", "value1.from.cli_value")

	assert.Equal(t, len(env.Providers["p1"].Parameters), 3)
	cp(t, env.Providers["p1"].Parameters, "param1", "val1_descriptor")
	cp(t, env.Providers["p1"].Parameters, "param2", "value1.from.cli_value")
	cp(t, env.Providers["p1"].Parameters, "param3", "value2.from.cli_value")

}

func TestTemplateOnDistributionVars(t *testing.T) {

	p, _ := model.CreateParameters(map[string]interface{}{
		"value1": map[interface{}]interface{}{
			"from": map[interface{}]interface{}{
				"cli": map[interface{}]interface{}{
					"to_distribution": "value_from_cli_to_distribution",
					"to_comp1":        "value_from_cli_to_comp1",
				},
			},
		},
	})
	mainPath := "./testdata/gittest/descriptor"
	tc := model.CreateContext(p)
	c := &MockLaunchContext{locationContent: mainPath, templateContext: tc}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repDesc := tester.createRep(mainPath)

	distContent := `
ekara:
vars:
  key1_distribution: val1_distribution
  key2_distribution: "{{ .Vars.value1.from.cli.to_distribution }}"
  key3_distribution: "{{ .Vars.key1_environment }}"
`
	repDist.writeCommit(t, "ekara.yaml", distContent)

	comp1Content := `
ekara:
vars:
  key1_comp1: val1_comp1
  key2_comp1: "{{ .Vars.key1_distribution }}"
  key3_comp1: "{{ .Vars.key2_distribution }}"
  key4_comp1: "{{ .Vars.value1.from.cli.to_comp1 }}"
  key5_comp1: "{{ .Vars.key1_environment }}"
`
	repComp1.writeCommit(t, "ekara.yaml", comp1Content)
	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  distribution:
    repository: ./testdata/gittest/distribution	
  components:
    comp1:
      repository: ./testdata/gittest/comp1	
vars:
  key1_environment: val1_environment
providers:
  comp1:
    component: comp1
nodes:
  node1:
    instances: 1
    provider:
      name: comp1
`

	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	tester.assertComponentsContains(model.MainComponentId, model.EkaraComponentId, "comp1")

	assert.Equal(t, len(tc.Vars), 10)
	cp(t, tc.Vars, "key1_comp1", "val1_comp1")
	cp(t, tc.Vars, "key2_comp1", "val1_distribution")
	cp(t, tc.Vars, "key3_comp1", "value_from_cli_to_distribution")
	cp(t, tc.Vars, "key4_comp1", "value_from_cli_to_comp1")
	cp(t, tc.Vars, "key5_comp1", "val1_environment")
	cp(t, tc.Vars, "key3_distribution", "val1_environment")
}

func cp(t *testing.T, p model.Parameters, key, value string) {
	v, ok := p[key]
	if assert.True(t, ok) {
		assert.Equal(t, value, v)
	}
}

func TestVarsPrecedence(t *testing.T) {

	p, _ := model.CreateParameters(map[string]interface{}{
		"key1": "value1.from.cli",
	})

	comp1Content := `
vars:
  key1: val1_comp1
  key2: val2_comp1
  key3: val3_comp1
  keyY: val4_comp1
`

	distContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
vars:
  key1: val1_distribution			
  key2: val2_distribution			
  key3: val3_distribution			
  keyX: val4_distribution			
`
	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  distribution:
    repository: ./testdata/gittest/distribution
vars:
  key1: val1_descriptor					
  key3: val3_descriptor					
# Following content just to force the download of comp1 and comp2
orchestrator:
  component: comp1
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

	tc := model.CreateContext(p)

	c := &MockLaunchContext{locationContent: mainPath, templateContext: tc}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repDesc := tester.createRep(mainPath)

	repComp1.writeCommit(t, "ekara.yaml", comp1Content)
	repDist.writeCommit(t, "ekara.yaml", distContent)
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	tester.assertComponentsContains(model.MainComponentId, model.EkaraComponentId, "comp1")

	assert.Equal(t, len(tc.Vars), 5)
	// Cli var has precedence over descriptor/distribution/comp1
	cp(t, tc.Vars, "key1", "value1.from.cli")
	// Descriptor var has precedence over distribution
	cp(t, tc.Vars, "key3", "val3_descriptor")
	// Distribution var has precedence over comp1
	cp(t, tc.Vars, "key2", "val2_distribution")
	// Test accumation of vars from distribution and its components
	cp(t, tc.Vars, "keyY", "val4_comp1")
	cp(t, tc.Vars, "keyX", "val4_distribution")
}

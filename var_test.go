package engine

import (
	"log"
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
	tester := gitTester(t, c)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repDesc := tester.createRep(mainPath)

	distContent := `
ekara:
`
	repDist.writeCommit(t, "ekara.yaml", distContent)

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  distribution:
    repository: ./testdata/gittest/distribution	
vars:
  key1_descriptor: val1_descriptor
  key2_descriptor: "{{ .Vars.value1.from.cli }}"

providers:
  ek-aws:
    component: ek-aws
    params:
      param1: {{ .Vars.key1_descriptor }}
      param2: {{ .Vars.key2_descriptor }}
      param3: {{ .Vars.value2 }} 
`

	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	// Check if the descriptor has been templated
	assert.Equal(t, len(env.Vars), 2)
	cp(t, env.Vars, "key1_descriptor", "val1_descriptor")
	cp(t, env.Vars, "key2_descriptor", "value1.from.cli_value")

	assert.Equal(t, len(env.Providers["ek-aws"].Parameters), 3)
	cp(t, env.Providers["ek-aws"].Parameters, "param1", "val1_descriptor")
	cp(t, env.Providers["ek-aws"].Parameters, "param2", "value1.from.cli_value")
	cp(t, env.Providers["ek-aws"].Parameters, "param3", "value2.from.cli_value")

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
	tester := gitTester(t, c)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repDesc := tester.createRep(mainPath)

	distContent := `
ekara:
vars:
  key1_distribution: val1_distribution
  key2_distribution: "{{ .Vars.value1.from.cli.to_distribution }}"
`
	repDist.writeCommit(t, "ekara.yaml", distContent)

	comp1Content := `
ekara:
vars:
  key1_comp1: val1_comp1
  key2_comp1: "{{ .Vars.key1_distribution }}"
  key3_comp1: "{{ .Vars.key2_distribution }}"
  key4_comp1: "{{ .Vars.value1.from.cli.to_comp1 }}"
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
	env := tester.env()
	assert.NotNil(t, env)

	tester.assertComponentsContains("__main__", "__ekara__", "comp1")

	assert.Equal(t, len(tc.Vars), 7)
	cp(t, tc.Vars, "key1_comp1", "val1_comp1")
	cp(t, tc.Vars, "key2_comp1", "val1_distribution")
	cp(t, tc.Vars, "key3_comp1", "value_from_cli_to_distribution")
	cp(t, tc.Vars, "key4_comp1", "value_from_cli_to_comp1")
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
	tester := gitTester(t, c)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repDesc := tester.createRep(mainPath)

	repComp1.writeCommit(t, "ekara.yaml", comp1Content)
	repDist.writeCommit(t, "ekara.yaml", distContent)
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	tester.assertComponentsContains("__main__", "__ekara__", "comp1")
	log.Printf("--> GBE vars %v", tc.Vars)
	assert.Equal(t, len(tc.Vars), 3)
	// Cli var has precedence over descriptor/distribution/comp1
	cp(t, tc.Vars, "key1", "value1.from.cli")
	// Descriptor var has precedence over distribution
	cp(t, tc.Vars, "key3", "val3_descriptor")
	// Distribution var has precedence over comp1
	cp(t, tc.Vars, "key2", "val2_distribution")
}

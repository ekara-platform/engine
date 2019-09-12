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
	c := &MockLaunchContext{locationContent: mainPath, data: p}
	tester := gitTester(t, c, false)
	defer tester.clean()

	tester.createRepDefaultDescriptor(t, "./testdata/gittest/parent")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp1")
	repDesc := tester.createRep(mainPath)

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent
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
	env := tester.env()
	assert.NotNil(t, env)

	// Check if the descriptor has been templated
	assert.Equal(t, len(env.Vars), 2)
	//Original value defined into the descriptor
	cp(t, env.Vars, "key1_descriptor", "val1_descriptor")
	//Value templated using the parameter file
	cp(t, env.Vars, "key2_descriptor", "value1.from.cli_value")

	assert.Equal(t, len(env.Providers["p1"].Parameters), 3)
	//Value templated using a value defined into the descriptor
	cp(t, env.Providers["p1"].Parameters, "param1", "val1_descriptor")
	//Value templated using a value previously templated into the descriptor
	cp(t, env.Providers["p1"].Parameters, "param2", "value1.from.cli_value")
	//Value templated using the parameter file
	cp(t, env.Providers["p1"].Parameters, "param3", "value2.from.cli_value")
}

func TestTemplateOnParentVars(t *testing.T) {

	p, _ := model.CreateParameters(map[string]interface{}{
		"value1": map[interface{}]interface{}{
			"from": map[interface{}]interface{}{
				"cli": map[interface{}]interface{}{
					"to_parent":     "value_from_cli_to_parent",
					"to_comp1":      "value_from_cli_to_comp1",
					"to_descriptor": "value_from_cli_to_descriptor",
				},
			},
		},
	})
	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, data: p}
	tester := gitTester(t, c, true)
	defer tester.clean()

	repParent := tester.createRep("./testdata/gittest/parent")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repDesc := tester.createRep(mainPath)

	comp1Content := `
vars:
  key1_comp1: val1_comp1
  key2_comp1: "{{ .Vars.value1.from.cli.to_comp1 }}"
`
	repComp1.writeCommit(t, "ekara.yaml", comp1Content)

	parentContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1	
vars:
  key1_parent: val1_parent
  key2_parent: "{{ .Vars.value1.from.cli.to_parent }}"
  key3_parent: "{{ .Vars.key1_comp1 }}"
`
	repParent.writeCommit(t, "ekara.yaml", parentContent)

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent
vars:
  key1_descriptor: val1_descriptor
  key2_descriptor: "{{ .Vars.value1.from.cli.to_descriptor }}"
  key3_descriptor: "{{ .Vars.key1_comp1 }}"
  key4_descriptor: "{{ .Vars.key2_comp1 }}"
  key5_descriptor: "{{ .Vars.key1_parent }}"
  key6_descriptor: "{{ .Vars.key2_parent }}"
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

	tester.assertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1")

	rc := tester.context.engine.Context()
	assert.Equal(t, len(rc.data.Vars), 12)
	cp(t, rc.data.Vars, "key1_comp1", "val1_comp1")
	// Should be templated with the cli params content
	cp(t, rc.data.Vars, "key2_comp1", "value_from_cli_to_comp1")
	cp(t, rc.data.Vars, "key1_parent", "val1_parent")
	// Should be templated with the cli params content
	cp(t, rc.data.Vars, "key2_parent", "value_from_cli_to_parent")
	// Should be templated from comp1 parameter
	cp(t, rc.data.Vars, "key3_parent", "val1_comp1")
	cp(t, rc.data.Vars, "key1_descriptor", "val1_descriptor")
	// Should be templated with the cli params content
	cp(t, rc.data.Vars, "key2_descriptor", "value_from_cli_to_descriptor")
	// Should be templated from comp1 parameter
	cp(t, rc.data.Vars, "key3_descriptor", "val1_comp1")
	// Should be templated from comp1 parameter
	cp(t, rc.data.Vars, "key4_descriptor", "value_from_cli_to_comp1")
	// Should be templated from parent parameter
	cp(t, rc.data.Vars, "key5_descriptor", "val1_parent")
	// Should be templated from parent parameter
	cp(t, rc.data.Vars, "key6_descriptor", "value_from_cli_to_parent")
}

func cp(t *testing.T, p model.Parameters, key, value string) {
	v, ok := p[key]
	if assert.True(t, ok) {
		assert.Equal(t, value, v)
	}
}

func TestVarsPrecedence(t *testing.T) {

	p, _ := model.CreateParameters(map[string]interface{}{
		"keyCli": "value4.from.cli",
	})

	comp1Content := `
vars:
  key1: val1_comp1
  key2: val2_comp1
  key3: val3_comp1
  key4: val4_comp1
`

	parentContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
vars:
  key2: val2_parent
  key3: val3_parent
  key4: val4_parent
`
	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent
vars:
  key3: val3_descriptor
  key4: "{{ .Vars.keyCli }}"
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

	c := &MockLaunchContext{locationContent: mainPath, data: p}
	tester := gitTester(t, c, true)
	defer tester.clean()

	repParent := tester.createRep("./testdata/gittest/parent")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repDesc := tester.createRep(mainPath)

	repComp1.writeCommit(t, "ekara.yaml", comp1Content)
	repParent.writeCommit(t, "ekara.yaml", parentContent)
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)

	env := tester.env()
	assert.NotNil(t, env)

	tester.assertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1")

	rc := tester.context.engine.Context()

	assert.Equal(t, len(rc.data.Vars), 5)
	cp(t, rc.data.Vars, "key1", "val1_comp1")
	cp(t, rc.data.Vars, "key2", "val2_parent")
	cp(t, rc.data.Vars, "key3", "val3_descriptor")
	cp(t, rc.data.Vars, "key4", "value4.from.cli")

}

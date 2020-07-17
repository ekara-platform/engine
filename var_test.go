package engine

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/engine/model"

	"github.com/stretchr/testify/assert"
)

func TestTemplateOnMainVars(t *testing.T) {

	p := model.CreateParameters(map[string]interface{}{
		"value1": map[interface{}]interface{}{
			"from": map[interface{}]interface{}{
				"cli": "value1.from.cli_value",
			},
		},
		"value2": "value2.from.cli_value",
	})
	mainPath := "descriptor"
	tester := util.CreateComponentTester(t, p)
	defer tester.Clean()

	tester.CreateDirEmptyDesc("parent")
	tester.CreateDirEmptyDesc("comp1")
	repDesc := tester.CreateDir(mainPath)

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: parent
  components:
    comp1:
      repository: comp1	
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

	repDesc.WriteCommit("ekara.yaml", descContent)

	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)

	assert.Len(t, env.Providers["p1"].Parameters(), 3)
	//Value templated using a value defined into the descriptor
	tester.AssertParam(env.Providers["p1"].Parameters(), "param1", "val1_descriptor")
	//Value templated using a value previously templated into the descriptor
	tester.AssertParam(env.Providers["p1"].Parameters(), "param2", "value1.from.cli_value")
	//Value templated using the parameter file
	tester.AssertParam(env.Providers["p1"].Parameters(), "param3", "value2.from.cli_value")
}

func TestTemplateOnParentVars(t *testing.T) {

	p := model.CreateParameters(map[string]interface{}{
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
	mainPath := "descriptor"
	tester := util.CreateComponentTester(t, p)
	defer tester.Clean()

	repParent := tester.CreateDir("parent")
	repComp1 := tester.CreateDir("comp1")
	repDesc := tester.CreateDir(mainPath)

	comp1Content := `
vars:
  key1_comp1: val1_comp1
  key2_comp1: "{{ .Vars.value1.from.cli.to_comp1 }}"
`
	repComp1.WriteCommit("ekara.yaml", comp1Content)

	parentContent := `
ekara:
  components:
    comp1:
      repository: comp1	
vars:
  key1_parent: val1_parent
  key2_parent: "{{ .Vars.value1.from.cli.to_parent }}"
  key3_parent: "{{ .Vars.key1_comp1 }}"
`
	repParent.WriteCommit("ekara.yaml", parentContent)

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: parent
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

	repDesc.WriteCommit("ekara.yaml", descContent)

	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)

	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1")

	tplC := tester.TemplateContext().(*model.TemplateContext)
	assert.Len(t, tplC.Vars, 12)
	tester.AssertParam(tplC.Vars, "key1_comp1", "val1_comp1")
	// Should be templated with the cli params content
	tester.AssertParam(tplC.Vars, "key2_comp1", "value_from_cli_to_comp1")
	tester.AssertParam(tplC.Vars, "key1_parent", "val1_parent")
	// Should be templated with the cli params content
	tester.AssertParam(tplC.Vars, "key2_parent", "value_from_cli_to_parent")
	// Should be templated from comp1 parameter
	tester.AssertParam(tplC.Vars, "key3_parent", "val1_comp1")
	tester.AssertParam(tplC.Vars, "key1_descriptor", "val1_descriptor")
	// Should be templated with the cli params content
	tester.AssertParam(tplC.Vars, "key2_descriptor", "value_from_cli_to_descriptor")
	// Should be templated from comp1 parameter
	tester.AssertParam(tplC.Vars, "key3_descriptor", "val1_comp1")
	// Should be templated from comp1 parameter
	tester.AssertParam(tplC.Vars, "key4_descriptor", "value_from_cli_to_comp1")
	// Should be templated from parent parameter
	tester.AssertParam(tplC.Vars, "key5_descriptor", "val1_parent")
	// Should be templated from parent parameter
	tester.AssertParam(tplC.Vars, "key6_descriptor", "value_from_cli_to_parent")
}

func TestVarsPrecedence(t *testing.T) {

	p := model.CreateParameters(map[string]interface{}{
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
      repository: comp1
vars:
  key2: val2_parent
  key3: val3_parent
  key4: val4_parent
`
	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: parent
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

	mainPath := "descriptor"

	tester := util.CreateComponentTester(t, p)
	defer tester.Clean()

	repParent := tester.CreateDir("parent")
	repComp1 := tester.CreateDir("comp1")
	repDesc := tester.CreateDir(mainPath)

	repComp1.WriteCommit("ekara.yaml", comp1Content)
	repParent.WriteCommit("ekara.yaml", parentContent)
	repDesc.WriteCommit("ekara.yaml", descContent)

	tester.Init(repDesc.AsRepository("master"))

	env := tester.Env()
	assert.NotNil(t, env)

	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1")

	tplC := tester.TemplateContext().(*model.TemplateContext)

	assert.Len(t, tplC.Vars, 5)
	tester.AssertParam(tplC.Vars, "key1", "val1_comp1")
	tester.AssertParam(tplC.Vars, "key2", "val2_parent")
	tester.AssertParam(tplC.Vars, "key3", "val3_descriptor")
	tester.AssertParam(tplC.Vars, "key4", "value4.from.cli")

}

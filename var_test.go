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
	c := MockLaunchContext{locationContent: mainPath, templateContext: tc}
	tester := gitTester(t, c)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repDesc := tester.createRep(mainPath)

	distContent := `
ekara:
vars:
  key1_distribution: val1_distribution
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
	//      param4: {{ .Model.key1_distribution }}
	// The vars comming from the distribution cannot me used into the descriptor.

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
	//	cp(t, env.Providers["ek-aws"].Parameters, "param4", "val1_distribution")

}

func cp(t *testing.T, p model.Parameters, key, value string) {
	v, ok := p[key]
	if assert.True(t, ok) {
		assert.Equal(t, value, v)
	}
}

package engine

import (
	"testing"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestVarsAccumulation(t *testing.T) {

	p, _ := model.CreateParameters(map[string]interface{}{
		"value1": map[interface{}]interface{}{
			"from": map[interface{}]interface{}{
				"cli": "value1.from.cli_value",
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
	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

	comp2Content := `
vars:
  key1_comp2: val1_comp2
  key2_comp2: val2_comp2
`
	repComp2.writeCommit(t, "ekara.yaml", comp2Content)

	comp1Content := `
ekara:
  components:
    comp2:
      repository: ./testdata/gittest/comp2
vars:
  key1_comp1: val1_comp1
  key2_comp1: val2_comp1
`
	repComp1.writeCommit(t, "ekara.yaml", comp1Content)

	distContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
vars:
  key1_distribution: val1_distribution
  key2_distribution: val2_distribution
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
  key2_descriptor: val2_descriptor

# Following content just to force the download of comp1 and comp2
orchestrator:
  component: comp1 

providers:
  p1:
    component: comp2

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

	// Check that all vars have been accumulated
	// From the descriptor
	assert.Equal(t, len(tc.Vars), 9)
	// From comp2
	cp(t, tc.Vars, "key1_comp2", "val1_comp2")
	cp(t, tc.Vars, "key2_comp2", "val2_comp2")
	// From comp1
	cp(t, tc.Vars, "key1_comp1", "val1_comp1")
	cp(t, tc.Vars, "key2_comp1", "val2_comp1")
	// From distribution
	cp(t, tc.Vars, "key1_distribution", "val1_distribution")
	cp(t, tc.Vars, "key2_distribution", "val2_distribution")
	// From descriptor
	cp(t, tc.Vars, "key1_descriptor", "val1_descriptor")
	cp(t, tc.Vars, "key2_descriptor", "val2_descriptor")
	// From the client
	_, ok := tc.Vars["value1"]
	assert.True(t, ok)

}

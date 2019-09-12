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
	c := &MockLaunchContext{locationContent: mainPath, data: p}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repParent := tester.createRep("./testdata/gittest/parent")
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
vars:
  key1_comp1: val1_comp1
  key2_comp1: val2_comp1
`
	repComp1.writeCommit(t, "ekara.yaml", comp1Content)

	parentContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
    comp2:
      repository: ./testdata/gittest/comp2
vars:
  key1_parent: val1_parent
  key2_parent: val2_parent
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
	env := tester.env()
	assert.NotNil(t, env)

	rc := tester.context.engine.Context()

	// Check that all vars have been accumulated
	// From the descriptor
	assert.Equal(t, len(rc.data.Vars), 9)
	// From comp2
	cp(t, rc.data.Vars, "key1_comp2", "val1_comp2")
	cp(t, rc.data.Vars, "key2_comp2", "val2_comp2")
	// From comp1
	cp(t, rc.data.Vars, "key1_comp1", "val1_comp1")
	cp(t, rc.data.Vars, "key2_comp1", "val2_comp1")
	// From parent
	cp(t, rc.data.Vars, "key1_parent", "val1_parent")
	cp(t, rc.data.Vars, "key2_parent", "val2_parent")
	// From descriptor
	cp(t, rc.data.Vars, "key1_descriptor", "val1_descriptor")
	cp(t, rc.data.Vars, "key2_descriptor", "val2_descriptor")
	// From the client
	_, ok := rc.data.Vars["value1"]
	assert.True(t, ok)

}

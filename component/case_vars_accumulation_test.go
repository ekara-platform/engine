package component

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestVarsAccumulation(t *testing.T) {

	p := model.CreateParameters(map[string]interface{}{
		"value1": map[interface{}]interface{}{
			"from": map[interface{}]interface{}{
				"cli": "value1.from.cli_value",
			},
		},
	})
	mainPath := "./testdata/gittest/descriptor"
	c := util.CreateMockLaunchContextWithData(mainPath, p, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	repParent := tester.CreateRep("./testdata/gittest/parent")
	repComp1 := tester.CreateRep("./testdata/gittest/comp1")
	repComp2 := tester.CreateRep("./testdata/gittest/comp2")
	repDesc := tester.CreateRep(mainPath)

	comp2Content := `
vars:
  key1_comp2: val1_comp2
  key2_comp2: val2_comp2
`
	repComp2.WriteCommit("ekara.yaml", comp2Content)

	comp1Content := `
vars:
  key1_comp1: val1_comp1
  key2_comp1: val2_comp1
`
	repComp1.WriteCommit("ekara.yaml", comp1Content)

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
	repParent.WriteCommit("ekara.yaml", parentContent)

	descContent := `
name: ekaraDemoVar
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
	repDesc.WriteCommit("ekara.yaml", descContent)

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)

	tplC := tester.tplC

	// Check that all vars have been accumulated
	// From the descriptor
	if assert.Len(t, tplC.Vars, 9) {
		// From comp2
		tester.CheckParameter("key1_comp2", "val1_comp2")
		tester.CheckParameter("key2_comp2", "val2_comp2")
		// From comp1
		tester.CheckParameter("key1_comp1", "val1_comp1")
		tester.CheckParameter("key2_comp1", "val2_comp1")
		// From parent
		tester.CheckParameter("key1_parent", "val1_parent")
		tester.CheckParameter("key2_parent", "val2_parent")
		// From descriptor
		tester.CheckParameter("key1_descriptor", "val1_descriptor")
		tester.CheckParameter("key2_descriptor", "val2_descriptor")
		// From the client
		_, ok := tplC.Vars["value1"]
		assert.True(t, ok)
	}
}

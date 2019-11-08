package component

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestFetchOrderedAlphabetical(t *testing.T) {
	mainPath := "./testdata/gittest/descriptor"
	c := util.CreateMockLaunchContext(mainPath, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	repDesc := tester.CreateRep(mainPath)
	tester.CreateRepDefaultDescriptor("./testdata/gittest/parent")
	tester.CreateRepDefaultDescriptor("./testdata/gittest/comp1")
	tester.CreateRepDefaultDescriptor("./testdata/gittest/comp2")
	tester.CreateRepDefaultDescriptor("./testdata/gittest/comp3")
	tester.CreateRepDefaultDescriptor("./testdata/gittest/comp4")
	tester.CreateRepDefaultDescriptor("./testdata/gittest/comp5")
	tester.CreateRepDefaultDescriptor("./testdata/gittest/comp6")

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent
  components:
    comp6:
      repository: ./testdata/gittest/comp6
    comp5:
      repository: ./testdata/gittest/comp5
    comp1:
      repository: ./testdata/gittest/comp1
    comp3:
      repository: ./testdata/gittest/comp3
    comp4:
      repository: ./testdata/gittest/comp4
    comp2:
      repository: ./testdata/gittest/comp2
orchestrator:
  component: comp1
providers:
  p1:
    component: comp2
  p2:
    component: comp3
  p3:
    component: comp4
  p4:
    component: comp5
  p5:
    component: comp6
nodes:
  node1:
    instances: 1
    provider:
      name: p1
  node2:
    instances: 1
    provider:
      name: p2
  node3:
    instances: 1
    provider:
      name: p3
  node4:
    instances: 1
    provider:
      name: p4
  node5:
    instances: 1
    provider:
      name: p5
`

	repDesc.WriteCommit("ekara.yaml", descContent)

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)

	tester.AssertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1", "comp2", "comp3", "comp4", "comp5", "comp6")

	assert.Equal(t, len(tester.rM.sortedFetchedComponents), 8)
	checkFetchOrder(tester, t, model.EkaraComponentId+"1", "comp1", "comp2", "comp3", "comp4", "comp5", "comp6", model.MainComponentId)
}

//
// Descriptor
//   Parent
//     Comp1
//     Comp2
//     Comp4
//     Comp5
//   Components
//     Comp3
//     Comp6
//
//
func TestFetchOrderedTwoParents(t *testing.T) {

	p := model.CreateParameters(map[string]interface{}{})
	mainPath := "./testdata/gittest/descriptor"
	c := util.CreateMockLaunchContextWithData(mainPath, p, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	repParent2 := tester.CreateRep("./testdata/gittest/parent2")
	repParent1 := tester.CreateRep("./testdata/gittest/parent1")
	repDesc := tester.CreateRep(mainPath)

	c1Rep := tester.CreateRep("./testdata/gittest/comp1")
	comp1Content := `
vars:
  key1: val1_comp1
  key2: val2_comp1
  key3: val3_comp1
  key4: val4_comp1
  key5: val5_comp1
  key6: val6_comp1
  key7: val7_comp1
  key8: val8_comp1
  key9: val9_comp1
`
	c1Rep.WriteCommit("ekara.yaml", comp1Content)

	c2Rep := tester.CreateRep("./testdata/gittest/comp2")
	comp2Content := `
vars:
  key2: val2_comp2
  key3: val3_comp2
  key4: val4_comp2
  key5: val5_comp2
  key6: val6_comp2
  key7: val7_comp2
  key8: val8_comp2
  key9: val9_comp2
`
	c2Rep.WriteCommit("ekara.yaml", comp2Content)

	parent2Content := `
ekara:
  components:
    comp2:
      repository: ./testdata/gittest/comp2
    comp1:
      repository: ./testdata/gittest/comp1

vars:
  key3: val3_ekara2
  key4: val4_ekara2
  key5: val5_ekara2
  key6: val6_ekara2
  key7: val7_ekara2
  key8: val8_ekara2
  key9: val9_ekara2
`
	repParent2.WriteCommit("ekara.yaml", parent2Content)

	c4Rep := tester.CreateRep("./testdata/gittest/comp4")
	comp4Content := `
vars:
  key4: val4_comp4
  key5: val5_comp4
  key6: val6_comp4
  key7: val7_comp4
  key8: val8_comp4
  key9: val9_comp4
`
	c4Rep.WriteCommit("ekara.yaml", comp4Content)

	c5Rep := tester.CreateRep("./testdata/gittest/comp5")
	comp5Content := `
vars:
  key5: val5_comp5
  key6: val6_comp5
  key7: val7_comp5
  key8: val8_comp5
  key9: val9_comp5
`
	c5Rep.WriteCommit("ekara.yaml", comp5Content)

	parent1Content := `
ekara:
  parent:
    repository: ./testdata/gittest/parent2
  components:
    comp5:
      repository: ./testdata/gittest/comp5
    comp4:
      repository: ./testdata/gittest/comp4

vars:
  key6: val6_ekara1
  key7: val7_ekara1
  key8: val8_ekara1
  key9: val9_ekara1
`
	repParent1.WriteCommit("ekara.yaml", parent1Content)

	c3Rep := tester.CreateRep("./testdata/gittest/comp3")
	comp3Content := `
vars:
  key7: val7_comp3
  key8: val8_comp3
  key9: val9_comp3
`
	c3Rep.WriteCommit("ekara.yaml", comp3Content)

	c6Rep := tester.CreateRep("./testdata/gittest/comp6")
	comp6Content := `
vars:
  key8: val8_comp6
  key9: val9_comp6
`
	c6Rep.WriteCommit("ekara.yaml", comp6Content)

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent1
  components:
    comp6:
      repository: ./testdata/gittest/comp6
    comp3:
      repository: ./testdata/gittest/comp3

vars:
  key9: val9_main
          
orchestrator:
  component: comp1
providers:
  p1:
    component: comp2
  p2:
    component: comp4
  p3:
    component: comp3
  p4:
    component: comp5
  p5:
    component: comp6
nodes:
  node1:
    instances: 1
    provider:
      name: p1
  node2:
    instances: 1
    provider:
      name: p2
  node3:
    instances: 1
    provider:
      name: p3
  node4:
    instances: 1
    provider:
      name: p4
  node5:
    instances: 1
    provider:
      name: p5
`
	repDesc.WriteCommit("ekara.yaml", descContent)

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)

	tester.AssertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", model.EkaraComponentId+"2", "comp1", "comp2", "comp3", "comp4", "comp5", "comp6")

	assert.Equal(t, len(tester.rM.sortedFetchedComponents), 9)
	// We need to fetch:
	//- first the components referenced by parent2
	//- then parent2 itself
	//- later the components referenced by parent1
	//- then parent1 itself
	//- later the components referenced by the main descriptor
	//- then main descriptor itself
	checkFetchOrder(tester, t, "comp1", "comp2", model.EkaraComponentId+"2", "comp4", "comp5", model.EkaraComponentId+"1", "comp3", "comp6", model.MainComponentId)

	// Check that all vars have been accumulated
	assert.Equal(t, len(tester.tplC.Vars), 9)

	tester.CheckParameter("key1", "val1_comp1")
	tester.CheckParameter("key2", "val2_comp2")
	tester.CheckParameter("key3", "val3_ekara2")
	tester.CheckParameter("key4", "val4_comp4")
	tester.CheckParameter("key5", "val5_comp5")
	tester.CheckParameter("key6", "val6_ekara1")
	tester.CheckParameter("key7", "val7_comp3")
	tester.CheckParameter("key8", "val8_comp6")
	tester.CheckParameter("key9", "val9_main")
}

func checkFetchOrder(tester *ComponentTester, t *testing.T, names ...string) {
	for i, v := range names {
		assert.Equal(t, tester.rM.sortedFetchedComponents[i], v)
	}
}

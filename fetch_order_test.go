package engine

import (
	"testing"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestFetchOrderedAlphabetical(t *testing.T) {

	p, _ := model.CreateParameters(map[string]interface{}{})
	mainPath := "./testdata/gittest/descriptor"

	c := &MockLaunchContext{locationContent: mainPath, data: p}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDesc := tester.createRep(mainPath)
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/parent")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp1")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp2")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp3")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp4")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp5")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp6")

	descContent := `
name: ekara-demo-var
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

	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	tester.assertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1", "comp2", "comp3", "comp4", "comp5", "comp6")

	assert.Equal(t, len(tester.context.engine.ReferenceManager().SortedFetchedComponents), 8)
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

	p, _ := model.CreateParameters(map[string]interface{}{})
	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, data: p}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repParent2 := tester.createRep("./testdata/gittest/parent2")
	repParent1 := tester.createRep("./testdata/gittest/parent1")
	repDesc := tester.createRep(mainPath)

	c1Rep := tester.createRep("./testdata/gittest/comp1")
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
	c1Rep.writeCommit(t, "ekara.yaml", comp1Content)

	c2Rep := tester.createRep("./testdata/gittest/comp2")
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
	c2Rep.writeCommit(t, "ekara.yaml", comp2Content)

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
	repParent2.writeCommit(t, "ekara.yaml", parent2Content)

	c4Rep := tester.createRep("./testdata/gittest/comp4")
	comp4Content := `
vars:
  key4: val4_comp4
  key5: val5_comp4
  key6: val6_comp4
  key7: val7_comp4
  key8: val8_comp4
  key9: val9_comp4
`
	c4Rep.writeCommit(t, "ekara.yaml", comp4Content)

	c5Rep := tester.createRep("./testdata/gittest/comp5")
	comp5Content := `
vars:
  key5: val5_comp5
  key6: val6_comp5
  key7: val7_comp5
  key8: val8_comp5
  key9: val9_comp5
`
	c5Rep.writeCommit(t, "ekara.yaml", comp5Content)

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
	repParent1.writeCommit(t, "ekara.yaml", parent1Content)

	c3Rep := tester.createRep("./testdata/gittest/comp3")
	comp3Content := `
vars:
  key7: val7_comp3
  key8: val8_comp3
  key9: val9_comp3
`
	c3Rep.writeCommit(t, "ekara.yaml", comp3Content)

	c6Rep := tester.createRep("./testdata/gittest/comp6")
	comp6Content := `
vars:
  key8: val8_comp6
  key9: val9_comp6
`
	c6Rep.writeCommit(t, "ekara.yaml", comp6Content)

	descContent := `
name: ekara-demo-var
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
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	tester.assertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", model.EkaraComponentId+"2", "comp1", "comp2", "comp3", "comp4", "comp5", "comp6")

	assert.Equal(t, len(tester.context.engine.ReferenceManager().SortedFetchedComponents), 9)
	// We need to fetch:
	//- first the components referenced by parent2
	//- then parent2 itself
	//- later the components referenced by parent1
	//- then parent1 itself
	//- later the components referenced by the main descriptor
	//- then main descriptor itself
	checkFetchOrder(tester, t, "comp1", "comp2", model.EkaraComponentId+"2", "comp4", "comp5", model.EkaraComponentId+"1", "comp3", "comp6", model.MainComponentId)

	rc := tester.context.engine.Context()

	// Check that all vars have been accumulated
	assert.Equal(t, len(rc.data.Vars), 9)

	cp(t, rc.data.Vars, "key1", "val1_comp1")
	cp(t, rc.data.Vars, "key2", "val2_comp2")
	cp(t, rc.data.Vars, "key3", "val3_ekara2")
	cp(t, rc.data.Vars, "key4", "val4_comp4")
	cp(t, rc.data.Vars, "key5", "val5_comp5")
	cp(t, rc.data.Vars, "key6", "val6_ekara1")
	cp(t, rc.data.Vars, "key7", "val7_comp3")
	cp(t, rc.data.Vars, "key8", "val8_comp6")
	cp(t, rc.data.Vars, "key9", "val9_main")
}

func checkFetchOrder(tester *tester, t *testing.T, names ...string) {
	for i, v := range names {
		assert.Equal(t, tester.context.engine.ReferenceManager().SortedFetchedComponents[i], v)
	}
}

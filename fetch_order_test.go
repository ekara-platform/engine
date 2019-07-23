package engine

import (
	"testing"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestFetchOrderedAlphabetical(t *testing.T) {

	p, _ := model.CreateParameters(map[string]interface{}{})
	mainPath := "./testdata/gittest/descriptor"
	tc := model.CreateContext(p)
	c := &MockLaunchContext{locationContent: mainPath, templateContext: tc}
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
	checkFetchOrder(tester, t, model.MainComponentId, model.EkaraComponentId+"1", "comp1", "comp2", "comp3", "comp4", "comp5", "comp6")
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
func TestFetchOrderedBase(t *testing.T) {

	p, _ := model.CreateParameters(map[string]interface{}{})
	mainPath := "./testdata/gittest/descriptor"
	tc := model.CreateContext(p)
	c := &MockLaunchContext{locationContent: mainPath, templateContext: tc}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDesc := tester.createRep(mainPath)
	repDist := tester.createRep("./testdata/gittest/parent")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp1")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp2")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp3")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp4")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp5")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp6")

	distContent := `
ekara:
  components:
    comp4:
      repository: ./testdata/gittest/comp4
    comp5:
      repository: ./testdata/gittest/comp5
    comp1:
      repository: ./testdata/gittest/comp1
    comp2:
      repository: ./testdata/gittest/comp2
`
	repDist.writeCommit(t, "ekara.yaml", distContent)

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent
  components:
    comp3:
      repository: ./testdata/gittest/comp3
    comp6:
      repository: ./testdata/gittest/comp6
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

	tester.assertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1", "comp2", "comp3", "comp4", "comp5", "comp6")

	assert.Equal(t, len(tester.context.engine.ReferenceManager().SortedFetchedComponents), 8)
	checkFetchOrder(tester, t, model.MainComponentId, model.EkaraComponentId+"1", "comp1", "comp2", "comp4", "comp5", "comp3", "comp6")
}

func checkFetchOrder(tester *tester, t *testing.T, names ...string) {
	for i, v := range names {
		assert.Equal(t, tester.context.engine.ReferenceManager().SortedFetchedComponents[i], v)
	}

}

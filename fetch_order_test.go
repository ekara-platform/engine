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
	tester := gitTester(t, c)
	defer tester.clean()

	repDesc := tester.createRep(mainPath)
	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repComp3 := tester.createRep("./testdata/gittest/comp3")
	repComp4 := tester.createRep("./testdata/gittest/comp4")
	repComp5 := tester.createRep("./testdata/gittest/comp5")
	repComp6 := tester.createRep("./testdata/gittest/comp6")

	repComp1.writeCommit(t, "ekara.yaml", "")
	repComp2.writeCommit(t, "ekara.yaml", "")
	repComp3.writeCommit(t, "ekara.yaml", "")
	repComp4.writeCommit(t, "ekara.yaml", "")
	repComp5.writeCommit(t, "ekara.yaml", "")
	repComp6.writeCommit(t, "ekara.yaml", "")
	repDist.writeCommit(t, "ekara.yaml", "")

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  distribution:
    repository: ./testdata/gittest/distribution
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
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	tester.assertComponentsContains("__main__", "__ekara__", "comp1", "comp2", "comp3", "comp4", "comp5", "comp6")

	assert.Equal(t, len(env.Ekara.SortedFetchedComponents), 8)
	checkFetchOrder(env, t, "__main__", "__ekara__", "comp1", "comp2", "comp3", "comp4", "comp5", "comp6")
}

//
// Descriptor
//   Distribution
//     Comp1
//     Comp2
//       --> Comp4
//         --> Comp5
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
	tester := gitTester(t, c)
	defer tester.clean()

	repDesc := tester.createRep(mainPath)
	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repComp3 := tester.createRep("./testdata/gittest/comp3")
	repComp4 := tester.createRep("./testdata/gittest/comp4")
	repComp5 := tester.createRep("./testdata/gittest/comp5")
	repComp6 := tester.createRep("./testdata/gittest/comp6")

	repComp1.writeCommit(t, "ekara.yaml", "")
	repComp3.writeCommit(t, "ekara.yaml", "")
	repComp5.writeCommit(t, "ekara.yaml", "")
	repComp6.writeCommit(t, "ekara.yaml", "")

	comp4Content := `
ekara:
  components:
    comp5:
      repository: ./testdata/gittest/comp5
`
	repComp4.writeCommit(t, "ekara.yaml", comp4Content)

	comp2Content := `
ekara:
  components:
    comp4:
      repository: ./testdata/gittest/comp4
`
	repComp2.writeCommit(t, "ekara.yaml", comp2Content)

	distContent := `
ekara:
  components:
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
  distribution:
    repository: ./testdata/gittest/distribution
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
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	tester.assertComponentsContains("__main__", "__ekara__", "comp1", "comp2", "comp3", "comp4", "comp5", "comp6")

	assert.Equal(t, len(env.Ekara.SortedFetchedComponents), 8)
	checkFetchOrder(env, t, "__main__", "__ekara__", "comp1", "comp2", "comp4", "comp5", "comp3", "comp6")
}

//
// Descriptor
//   Distribution
//     Comp2
//     Comp1
//       --> Comp4
//         --> Comp5
//   Components
//     Comp6
//     Comp3
//
//
func TestFetchOrderedSwitched(t *testing.T) {

	p, _ := model.CreateParameters(map[string]interface{}{})
	mainPath := "./testdata/gittest/descriptor"
	tc := model.CreateContext(p)
	c := &MockLaunchContext{locationContent: mainPath, templateContext: tc}
	tester := gitTester(t, c)
	defer tester.clean()

	repDesc := tester.createRep(mainPath)
	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repComp3 := tester.createRep("./testdata/gittest/comp3")
	repComp4 := tester.createRep("./testdata/gittest/comp4")
	repComp5 := tester.createRep("./testdata/gittest/comp5")
	repComp6 := tester.createRep("./testdata/gittest/comp6")

	repComp2.writeCommit(t, "ekara.yaml", "")
	repComp3.writeCommit(t, "ekara.yaml", "")
	repComp5.writeCommit(t, "ekara.yaml", "")
	repComp6.writeCommit(t, "ekara.yaml", "")

	comp4Content := `
ekara:
  components:
    comp5:
      repository: ./testdata/gittest/comp5
`
	repComp4.writeCommit(t, "ekara.yaml", comp4Content)

	comp1Content := `
ekara:
  components:
    comp4:
      repository: ./testdata/gittest/comp4
`
	repComp1.writeCommit(t, "ekara.yaml", comp1Content)

	distContent := `
ekara:
  components:
    comp2:
      repository: ./testdata/gittest/comp2
    comp1:
      repository: ./testdata/gittest/comp1
`
	repDist.writeCommit(t, "ekara.yaml", distContent)

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  distribution:
    repository: ./testdata/gittest/distribution
  components:
    comp6:
      repository: ./testdata/gittest/comp6
    comp3:
      repository: ./testdata/gittest/comp3
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
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	tester.assertComponentsContains("__main__", "__ekara__", "comp1", "comp2", "comp3", "comp4", "comp5", "comp6")

	assert.Equal(t, len(env.Ekara.SortedFetchedComponents), 8)
	checkFetchOrder(env, t, "__main__", "__ekara__", "comp1", "comp4", "comp5", "comp2", "comp3", "comp6")
}

func checkFetchOrder(env model.Environment, t *testing.T, names ...string) {
	for i, v := range names {
		assert.Equal(t, env.Ekara.SortedFetchedComponents[i], v)
	}

}

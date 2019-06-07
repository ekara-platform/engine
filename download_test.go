package engine

import (
	"testing"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

// when the descriptor doesn't define its own specific distribution then
// the defaulted one should be used
func TestDownloadDefaultDistribution(t *testing.T) {
	p, _ := model.CreateParameters(map[string]interface{}{
		"ek": map[interface{}]interface{}{
			"aws": map[interface{}]interface{}{
				"region": "dummy",
				"accessKey": map[interface{}]interface{}{
					"id":     "dummy",
					"secret": "dummy",
				},
			},
		},
	})
	mainPath := "./testdata/gittest/descriptor"
	tc := model.CreateContext(p)
	c := &MockLaunchContext{locationContent: mainPath, templateContext: tc}
	tester := gitTester(t, c)
	defer tester.clean()

	repDesc := tester.createRep(mainPath)

	descContent := `
name: ekara-demo-var
qualifier: dev

# Following content just to force the download of ek-swam and ek-aws

nodes:
  node1:
    instances: 1
    provider:
      name: ek-aws
`
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	// The defaulted distribution should comme with ek-aws as provider
	// and ek-swarm as orchestrator
	tester.assertComponentsContains("__main__", "__ekara__", "ek-swarm", "ek-aws")
}

func TestDownloadCustomDistribution(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

	repComp2.writeCommit(t, "ekara.yaml", ``)
	repComp1.writeCommit(t, "ekara.yaml", ``)

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
	// comp1 and comp2 should be downloaded because they are used into the descriptor
	tester.assertComponentsContains("__main__", "__ekara__", "comp1", "comp2")

}

func TestDownloadOnlyUsedComponents(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

	repComp2.writeCommit(t, "ekara.yaml", ``)
	repComp1.writeCommit(t, "ekara.yaml", ``)

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
`
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	// comp1 and comp2 shouldn't be downloaded because they are not used into the descriptor
	tester.assertComponentsContainsExactly("__main__", "__ekara__")

}

func TestDonwloadComplex(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

	repComp2.writeCommit(t, "ekara.yaml", ``)

	comp1Content := `
ekara:
  components:
    comp2:
      repository: ./testdata/gittest/comp2
`
	repComp1.writeCommit(t, "ekara.yaml", comp1Content)

	distContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
orchestrator:
  component: comp1 
`
	repDist.writeCommit(t, "ekara.yaml", distContent)

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  distribution:
    repository: ./testdata/gittest/distribution	

# Following content just to force the download of comp1 and comp2
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
	// comp1 should be downloaded because it's used as orchestrator into the distribution
	// comp2 should be also downloaded because it's used as provider into the descriptor
	tester.assertComponentsContainsExactly("__main__", "__ekara__", "comp1", "comp2")
}

// When more than one ekara.yaml file define a distribution the one taken
// in account should the the one defined in the main descriptor
func TestDownloadFirstDistribution(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c)
	defer tester.clean()

	repDist1 := tester.createRep("./testdata/gittest/distribution1")
	repDist2 := tester.createRep("./testdata/gittest/distribution2")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repComp3 := tester.createRep("./testdata/gittest/comp3")
	repComp4 := tester.createRep("./testdata/gittest/comp4")
	repDesc := tester.createRep(mainPath)

	repComp4.writeCommit(t, "ekara.yaml", ``)
	repComp3.writeCommit(t, "ekara.yaml", ``)

	// Comp2 defines another distribution but this
	// one should be ignored
	comp2Content := `
ekara:
  distribution:
    repository: ./testdata/gittest/distribution2
`
	repComp2.writeCommit(t, "ekara.yaml", comp2Content)
	repComp1.writeCommit(t, "ekara.yaml", ``)

	distContent1 := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
    comp2:
      repository: ./testdata/gittest/comp2
`
	repDist1.writeCommit(t, "ekara.yaml", distContent1)

	distContent2 := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp3
    comp2:
      repository: ./testdata/gittest/comp4
`
	repDist2.writeCommit(t, "ekara.yaml", distContent2)

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  distribution:
    repository: ./testdata/gittest/distribution1

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
	// comp1 and comp2 should be downloaded because they are used into the descriptor
	tester.assertComponentsContains("__main__", "__ekara__", "comp1", "comp2")
	cpnts := env.Ekara.Components
	assert.Equal(t, len(cpnts), 2)
	assert.Contains(t, cpnts, "comp1")
	assert.Contains(t, cpnts, "comp2")

}

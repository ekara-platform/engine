package engine

import (
	"testing"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)



func TestDownloadOnlyUsedComponents(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repComp3 := tester.createRep("./testdata/gittest/comp3")
	repComp4 := tester.createRep("./testdata/gittest/comp4")
	repDesc := tester.createRep(mainPath)
	
	repComp1.writeCommit(t, "ekara.yaml", ``)
	repComp2.writeCommit(t, "ekara.yaml", ``)
	repComp3.writeCommit(t, "ekara.yaml", ``)
	repComp4.writeCommit(t, "ekara.yaml", ``)

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
      comp4:
        repository: ./testdata/gittest/comp4
`
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	// comp1, comp2, comp3 and comp4 shouldn't be downloaded because they are not used into the descriptor
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
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator into the distribution
	// comp2 should be also downloaded because it's used as provider into the descriptor
	tester.assertComponentsContainsExactly("__main__", "__ekara__", "comp1", "comp2")
}

func TestDonwloadComplexFarReference(t *testing.T) {

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
orchestrator:
  component: comp2 
`
	repComp1.writeCommit(t, "ekara.yaml", comp1Content)

	distContent := `
ekara:
  components:
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

# Following content just to force the download of comp1 and comp2
providers:
  p1:
    component: comp1

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
	// comp1 should be downloaded because it's used as orchestrator into comp1
	// comp2 should be also downloaded because it's used as provider into the descriptor
	tester.assertComponentsContainsExactly("__main__", "__ekara__", "comp1", "comp2")
}



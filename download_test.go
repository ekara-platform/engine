package engine

import (
	"testing"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestDownloadOnlyUsedComponents(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
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
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId)

}

func TestDonwloadComplex(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
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
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId, "comp1", "comp2")
}

func TestDonwloadComplexFarReference(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
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
  component: comp2 
providers:
  p1:
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
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId, "comp1", "comp2")
}

func TestDonwloadFarInDistribution(t *testing.T) {
	comp1Content := `
ekara:
  components:
    comp2:
      repository: ./testdata/gittest/comp2
`

	distContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
orchestrator:
  component: comp2 
providers:
  p1:
    component: comp1  
`
	CheckDonwloadComplexFarReference(t, comp1Content, distContent)
}

func TestDonwloadFarInComp1(t *testing.T) {

	comp1Content := `
ekara:
  components:
    comp2:
      repository: ./testdata/gittest/comp2
orchestrator:
  component: comp2 
`

	distContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
providers:
  p1:
    component: comp1  
`
	CheckDonwloadComplexFarReference(t, comp1Content, distContent)
}

func CheckDonwloadComplexFarReference(t *testing.T, comp1Content, distContent string) {

	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

	repComp2.writeCommit(t, "ekara.yaml", ``)

	repDist.writeCommit(t, "ekara.yaml", distContent)
	repComp1.writeCommit(t, "ekara.yaml", comp1Content)

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  distribution:
    repository: ./testdata/gittest/distribution	

# Following content just to force the download of comp1 and comp2

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
	// comp1 should be downloaded because it's used as orchestrator
	// comp2 should be also downloaded because it's used as provider
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId, "comp1", "comp2")
}

func TestDonwloadFarProviderSplitted(t *testing.T) {
	comp1Content := `
providers:
  p1:
    params:
      key_1: value_1
    env:
      env_1: value_1	  
`

	distContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
orchestrator:
  component: comp1 	  
providers:
  p1:
    component: comp1  
`
	CheckDonwloadSplitted(t, comp1Content, distContent)
}

func CheckDonwloadSplitted(t *testing.T, comp1Content, distContent string) {

	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

	repComp2.writeCommit(t, "ekara.yaml", ``)

	repDist.writeCommit(t, "ekara.yaml", distContent)
	repComp1.writeCommit(t, "ekara.yaml", comp1Content)

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  distribution:
    repository: ./testdata/gittest/distribution	

# Following content just to force the download of comp1

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
	// comp1 should be downloaded because it's used as orchestrator
	// comp2 should be also downloaded because it's used as provider
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId, "comp1")

	if assert.Equal(t, 1, len(env.Providers)) {
		p := env.Providers["p1"]
		if assert.Equal(t, 1, len(p.Parameters)) {
			cp(t, p.Parameters, "key_1", "value_1")
		}

		if assert.Equal(t, 1, len(p.EnvVars)) {
			val, ok := p.EnvVars["env_1"]
			assert.True(t, ok)
			assert.Equal(t, "value_1", val)
		}
	}
}

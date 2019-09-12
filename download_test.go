package engine

import (
	"testing"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestDownloadOnlyUsedComponents(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	c := &MockLaunchContext{locationContent: mainPath, data: model.Parameters{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repParent := tester.createRep("./testdata/gittest/parent")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp1")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp2")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp3")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp4")
	repDesc := tester.createRep(mainPath)

	parentContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
    comp2:
      repository: ./testdata/gittest/comp2
`
	repParent.writeCommit(t, "ekara.yaml", parentContent)

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent
    components:
      comp3:
        repository: ./testdata/gittest/comp3
      comp4:
        repository: ./testdata/gittest/comp4
`
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	refM := tester.context.engine.ReferenceManager()
	assert.Equal(t, len(refM.UsedReferences.Refs), 0)

	// comp1, comp2, comp3 and comp4 shouldn't be downloaded because they are not used into the descriptor
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1")

}

func TestDonwloadOnlyFromDescriptorAndParents(t *testing.T) {
	comp1Content := `
ekara:
  components:
    comp2:
      repository: ./testdata/gittest/comp2
`

	parentContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
orchestrator:
  component: comp1 
`
	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent	

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
	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, data: model.Parameters{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repParent := tester.createRep("./testdata/gittest/parent")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

	repDesc.writeCommit(t, "ekara.yaml", descContent)
	repParent.writeCommit(t, "ekara.yaml", parentContent)
	repComp1.writeCommit(t, "ekara.yaml", comp1Content)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	refM := tester.context.engine.ReferenceManager()
	assert.Equal(t, len(refM.UsedReferences.Refs), 2)
	assert.True(t, refM.UsedReferences.IdUsed("comp1"))
	assert.True(t, refM.UsedReferences.IdUsed("comp2"))

	// comp1 should be downloaded because it's used as orchestrator into the parent
	// comp2 should not be downloaded because it's referenced by a component
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", "comp1")
}

func TestDonwloadTwoParents(t *testing.T) {

	parent2Content := `
ekara:
  components:
    comp2:
      repository: ./testdata/gittest/comp2
`
	parent1Content := `
ekara:
  parent:
    repository: ./testdata/gittest/parent2
  components:
    comp1:
      repository: ./testdata/gittest/comp1
`
	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent1

# Following content just to force the download of comp1 and comp2
orchestrator:
  component: comp2 
providers:
  p1:
    component: comp1  

nodes:
  node1:
    instances: 1
    provider:
      name: p1
`
	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, data: model.Parameters{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repParent1 := tester.createRep("./testdata/gittest/parent1")
	repParent2 := tester.createRep("./testdata/gittest/parent2")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp1")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

	repParent1.writeCommit(t, "ekara.yaml", parent1Content)
	repParent2.writeCommit(t, "ekara.yaml", parent2Content)
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	refM := tester.context.engine.ReferenceManager()
	assert.Equal(t, len(refM.UsedReferences.Refs), 2)
	assert.True(t, refM.UsedReferences.IdUsed("comp1"))
	assert.True(t, refM.UsedReferences.IdUsed("comp2"))

	// comp1 should be downloaded because it's used as orchestrator
	// comp2 should be also downloaded because it's used as provider
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", model.EkaraComponentId+"2", "comp1", "comp2")
}

func TestDonwloadTwoParentsUpperUsed(t *testing.T) {

	parent2Content := `
ekara:
  components:
    comp2:
      repository: ./testdata/gittest/comp2
orchestrator:
  component: comp2 
providers:
  p1:
    component: comp1  
`
	parent1Content := `
ekara:
  parent:
    repository: ./testdata/gittest/parent2
  components:
    comp1:
      repository: ./testdata/gittest/comp1
`
	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent1

# Following content just to force the download of comp1 and comp2
nodes:
  node1:
    instances: 1
    provider:
      name: p1
`
	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, data: model.Parameters{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repParent1 := tester.createRep("./testdata/gittest/parent1")
	repParent2 := tester.createRep("./testdata/gittest/parent2")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp1")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

	repParent1.writeCommit(t, "ekara.yaml", parent1Content)
	repParent2.writeCommit(t, "ekara.yaml", parent2Content)
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	refM := tester.context.engine.ReferenceManager()
	assert.Equal(t, len(refM.UsedReferences.Refs), 2)
	assert.True(t, refM.UsedReferences.IdUsed("comp1"))
	assert.True(t, refM.UsedReferences.IdUsed("comp2"))

	// comp1 should be downloaded because it's used as orchestrator
	// comp2 should be also downloaded because it's used as provider
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", model.EkaraComponentId+"2", "comp1", "comp2")
}

/*
func TestDonwloadFarProviderSplitted(t *testing.T) {
	comp1Content := `
providers:
  p1:
    params:
      key_1: value_1
    env:
      env_1: value_1
`

	parentContent := `
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
	CheckDonwloadSplitted(t, comp1Content, parentContent)
}

func CheckDonwloadSplitted(t *testing.T, comp1Content, parentContent string) {

	mainPath := "./testdata/gittest/descriptor"
	c := &MockLaunchContext{locationContent: mainPath, templateContext: tester.context.engine.Context().data}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repParent := tester.createRep("./testdata/gittest/parent")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	tester.createRepDefaultDescriptor(t,"./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

	repParent.writeCommit(t, "ekara.yaml", parentContent)
	repComp1.writeCommit(t, "ekara.yaml", comp1Content)

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent

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
	env := tester.env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator
	// comp2 should be also downloaded because it's used as provider
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", "comp1")

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
*/

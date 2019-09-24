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

// The orchestrator and the providers, once defined into a parent
// can be customized lower into the hierarchy of descriptors.
// In this case the component will be ommited, because it has already been defined,
// and the should not affect the download process
func TestDonwloadTwoParentsProviderUpperUsedRedefined(t *testing.T) {

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

orchestrator:
  params:
    dummy_key: dummy_value
providers:
  p1:
    params:
      dummy_key: dummy_value


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

// A stack  once defined into a parent
// can be customized lower into the hierarchy of descriptors.
// In this case the component will be ommited, because it has already been defined,
// and the should not affect the download process
func TestDonwloadTwoParentsStackUpperUsedRedefined(t *testing.T) {

	parent2Content := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
    comp2:
      repository: ./testdata/gittest/comp2
    comp3:
      repository: ./testdata/gittest/comp3

orchestrator:
  component: comp1
providers:
  p1:
    component: comp2
stacks:
  s1:
    component: comp3
`
	parent1Content := `
ekara:
  parent:
    repository: ./testdata/gittest/parent2
`

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent1

stacks:
  s1:
    params:
      dummy_key: dummy_value

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
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp3")
	repDesc := tester.createRep(mainPath)

	repParent1.writeCommit(t, "ekara.yaml", parent1Content)
	repParent2.writeCommit(t, "ekara.yaml", parent2Content)
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	refM := tester.context.engine.ReferenceManager()
	assert.Equal(t, len(refM.UsedReferences.Refs), 3)
	assert.True(t, refM.UsedReferences.IdUsed("comp1"))
	assert.True(t, refM.UsedReferences.IdUsed("comp2"))
	assert.True(t, refM.UsedReferences.IdUsed("comp3"))

	// comp1 should be downloaded because it's used as orchestrator
	// comp2 should be also downloaded because it's used as provider
	// comp3 should be also downloaded because it's used as a stack
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", model.EkaraComponentId+"2", "comp1", "comp2", "comp3")
}

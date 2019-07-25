package engine

import (
	"testing"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

// when the descriptor doesn't define its own specific parent then
// the defaulted one should be used
func TestDownloadDefaultParent2(t *testing.T) {
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
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDesc := tester.createRep(mainPath)

	descContent := `
name: ekara-demo-var
qualifier: dev

# Following content just to force the download of ek-swam and ek-aws

providers:
  ek-aws:
    component: ek-aws
nodes:
  node1:
    instances: 1
    provider:
      name: ek-aws
`
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)

	refM := tester.context.engine.ReferenceManager()
	assert.Equal(t, len(refM.UsedReferences.Refs), 3)

	assert.True(t, refM.UsedReferences.IdUsed("ek-swarm"))
	assert.True(t, refM.UsedReferences.IdUsed("ek-aws"))
	assert.True(t, refM.UsedReferences.IdUsed("ek-core"))

	assert.Equal(t, len(refM.ReferencedComponents.Refs), 4)
	assert.True(t, refM.ReferencedComponents.IdReferenced("ek-swarm"))
	assert.True(t, refM.ReferencedComponents.IdReferenced("ek-aws"))
	assert.True(t, refM.ReferencedComponents.IdReferenced("ek-openstack"))
	assert.True(t, refM.ReferencedComponents.IdReferenced("ek-core"))

	assert.Len(t, refM.Parents, 1)
	// Check that the parent has been renamed base on its position
	assert.Equal(t, model.EkaraComponentId+"1", refM.Parents[0].Component.Id)

	// Check the referenced components has been cleaned
	refM.ReferencedComponents.Clean(*refM.UsedReferences)
	assert.Equal(t, len(refM.ReferencedComponents.Refs), 3)
	assert.True(t, refM.ReferencedComponents.IdReferenced("ek-swarm"))
	assert.True(t, refM.ReferencedComponents.IdReferenced("ek-aws"))
	assert.True(t, refM.ReferencedComponents.IdReferenced("ek-core"))

	env := tester.env()
	assert.NotNil(t, env)
	// The defaulted parent should comme with ek-aws as provider
	// and ek-swarm as orchestrator
	tester.assertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "ek-swarm", "ek-aws", "ek-core")
	cpnts := env.Platform().Components
	assert.Equal(t, len(cpnts), 5)
	assert.Contains(t, cpnts, "__main__")
	assert.Contains(t, cpnts, "__ekara__1")
	assert.Contains(t, cpnts, "ek-swarm")
	assert.Contains(t, cpnts, "ek-aws")
	assert.Contains(t, cpnts, "ek-core")
}

func TestDownloadCustomParent(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/parent")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp1")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

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
  parent:
    repository: ./testdata/gittest/parent

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

	refM := tester.context.engine.ReferenceManager()
	assert.Equal(t, 2, len(refM.UsedReferences.Refs))
	assert.True(t, refM.UsedReferences.IdUsed("comp1"))
	assert.True(t, refM.UsedReferences.IdUsed("comp2"))
	assert.Equal(t, 2, len(refM.ReferencedComponents.Refs))
	assert.True(t, refM.ReferencedComponents.IdReferenced("comp1"))
	assert.True(t, refM.ReferencedComponents.IdReferenced("comp2"))

	assert.Len(t, refM.Parents, 1)
	// Check that the parent has been renamed base on its position
	assert.Equal(t, model.EkaraComponentId+"1", refM.Parents[0].Component.Id)

	// Check the referenced components has not been cleaned
	refM.ReferencedComponents.Clean(*refM.UsedReferences)
	assert.Equal(t, len(refM.ReferencedComponents.Refs), 2)
	assert.True(t, refM.ReferencedComponents.IdReferenced("comp1"))
	assert.True(t, refM.ReferencedComponents.IdReferenced("comp2"))

	env := tester.env()
	assert.NotNil(t, env)
	// comp1 and comp2 should be downloaded because they are used into the descriptor
	tester.assertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1", "comp2")
	cpnts := env.Platform().Components
	assert.Equal(t, len(cpnts), 4)
	assert.Contains(t, cpnts, "__main__")
	assert.Contains(t, cpnts, "__ekara__1")
	assert.Contains(t, cpnts, "comp1")
	assert.Contains(t, cpnts, "comp2")
}

// When more than one ekara.yaml file define a parent the one taken
// in account should the the one defined in the main descriptor
func TestDownloadFTwoParents(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDist1 := tester.createRep("./testdata/gittest/parent1")
	repDist2 := tester.createRep("./testdata/gittest/parent2")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp1")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp2")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp3")
	tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp4")
	repDesc := tester.createRep(mainPath)

	distContent1 := `
ekara:
  parent:
    repository: ./testdata/gittest/parent2
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
    comp3:
      repository: ./testdata/gittest/comp3
    comp4:
      repository: ./testdata/gittest/comp4
`
	repDist2.writeCommit(t, "ekara.yaml", distContent2)

	descContent := `
name: ekara-demo-var
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent1

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

	refM := tester.context.engine.ReferenceManager()
	assert.Equal(t, len(refM.UsedReferences.Refs), 2)
	assert.True(t, refM.UsedReferences.IdUsed("comp1"))
	assert.True(t, refM.UsedReferences.IdUsed("comp2"))
	assert.Equal(t, len(refM.ReferencedComponents.Refs), 4)
	assert.True(t, refM.ReferencedComponents.IdReferenced("comp1"))
	assert.True(t, refM.ReferencedComponents.IdReferenced("comp2"))
	assert.True(t, refM.ReferencedComponents.IdReferenced("comp3"))
	assert.True(t, refM.ReferencedComponents.IdReferenced("comp4"))
	assert.Len(t, refM.Parents, 2)
	// Check that the parents has been renamed base on their position
	assert.Equal(t, model.EkaraComponentId+"1", refM.Parents[0].Component.Id)
	assert.Equal(t, model.EkaraComponentId+"2", refM.Parents[1].Component.Id)

	// Check the referenced components has been cleaned
	refM.ReferencedComponents.Clean(*refM.UsedReferences)
	assert.Equal(t, len(refM.ReferencedComponents.Refs), 2)
	assert.True(t, refM.ReferencedComponents.IdReferenced("comp1"))
	assert.True(t, refM.ReferencedComponents.IdReferenced("comp2"))

	env := tester.env()
	assert.NotNil(t, env)
	// comp1 and comp2 should be downloaded because they are used into the descriptor
	tester.assertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", model.EkaraComponentId+"2", "comp1", "comp2")
	cpnts := env.Platform().Components
	assert.Equal(t, len(cpnts), 5)
	assert.Contains(t, cpnts, "__main__")
	assert.Contains(t, cpnts, "__ekara__1")
	assert.Contains(t, cpnts, "__ekara__2")
	assert.Contains(t, cpnts, "comp1")
	assert.Contains(t, cpnts, "comp2")

}

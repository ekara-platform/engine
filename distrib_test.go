package engine

import (
	"github.com/ekara-platform/engine/model"
	"github.com/ekara-platform/engine/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDownloadNoParent(t *testing.T) {
	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()
	repDesc := tester.CreateDir("descriptor")
	descContent := `
name: ekaraDemoVar
qualifier: dev
`
	repDesc.WriteCommit("ekara.yaml", descContent)
	tester.Init(repDesc.AsRepository("master"))
	assert.Equal(t, 1, tester.ComponentCount())
	tester.AssertComponentAvailable(model.MainComponentId)
}

// when the descriptor doesn't define its own specific parent then
// the defaulted one should be used
func TestDownloadDistribution(t *testing.T) {
	p := model.CreateParameters(map[string]interface{}{
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
	tester := util.CreateComponentTester(t, p)
	defer tester.Clean()
	repDesc := tester.CreateDir("descriptor")
	repDesc.WriteCommit(
		"ekara.yaml",
		`
name: ekaraDemoVar
qualifier: dev

ekara:
 base: https://github.com
 parent:
   repository: ekara-platform/distribution
   ref: master

# Following content just to force the download of ek-swarm and ek-aws
providers:
 ek-aws:
   component: ek-aws
nodes:
 node1:
   instances: 1
   provider:
     name: ek-aws
`)
	tester.Init(repDesc.AsRepository("master"))
	tester.AssertComponentAvailable(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "ek-aws", "ek-swarm")
	env := tester.Env()
	assert.NotNil(t, env)

	// Looking for the availability of a the deploy.yaml playbook
	use, err := tester.ComponentManager().Use(env.Orchestrator, tester.TemplateContext())
	assert.Nil(t, err)
	//defer use.Release()
	tester.AssertFile(use, "deploy.yaml")
	_, path := use.ContainsFile("deploy.yaml")
	assert.Equal(t, "ek-swarm", path.Owner().Id())
}

//func TestDownloadCustomParent(t *testing.T) {
//
//	mainPath := "descriptor"
//
//	c := util.CreateMockLaunchContext(mainPath, false)
//	tester := util.CreateComponentTester(t, c)
//	defer tester.Clean()
//
//	repParent := tester.CreateDir("parent")
//	tester.CreateDirEmptyDesc("comp1")
//	tester.CreateDirEmptyDesc("comp2")
//	repDesc := tester.CreateDir(mainPath)
//
//	parentContent := `
//ekara:
// components:
//   comp1:
//     repository: comp1
//   comp2:
//     repository: comp2
//`
//	repParent.WriteCommit("ekara.yaml", parentContent)
//
//	descContent := `
//name: ekaraDemoVar
//qualifier: dev
//
//ekara:
// parent:
//   repository: parent
//
//# Following content just to force the download of comp1 and comp2
//orchestrator:
// component: comp1
//
//providers:
// p1:
//   component: comp2
//
//nodes:
// node1:
//   instances: 1
//   provider:
//     name: p1
//`
//	repDesc.WriteCommit("ekara.yaml", descContent)
//
//	err := tester.Init()
//	assert.Nil(t, err)
//
//	refM := tester.rM
//	assert.Equal(t, 2, len(refM.usedReferences.Refs))
//	assert.True(t, refM.usedReferences.IdUsed("comp1"))
//	assert.True(t, refM.usedReferences.IdUsed("comp2"))
//	assert.Equal(t, 2, len(refM.referencedComponents.Refs))
//	assert.True(t, refM.referencedComponents.IdReferenced("comp1"))
//	assert.True(t, refM.referencedComponents.IdReferenced("comp2"))
//
//	assert.Len(t, refM.parents, 1)
//	// Check that the parent has been renamed base on its position
//	assert.Equal(t, model.ParentComponentId+"1", refM.parents[0].comp.Id)
//
//	// Check the referenced components has not been cleaned
//	refM.referencedComponents.Clean(*refM.usedReferences)
//	assert.Equal(t, len(refM.referencedComponents.Refs), 2)
//	assert.True(t, refM.referencedComponents.IdReferenced("comp1"))
//	assert.True(t, refM.referencedComponents.IdReferenced("comp2"))
//
//	env := tester.Env()
//	assert.NotNil(t, env)
//	// comp1 and comp2 should be downloaded because they are used into the descriptor
//	tester.AssertComponentsContains(model.MainComponentId, model.ParentComponentId+"1", "comp1", "comp2")
//	cpnts := env.Ekara.Components
//	assert.Equal(t, len(cpnts), 4)
//	assert.Contains(t, cpnts, "__main__")
//	assert.Contains(t, cpnts, model.ParentComponentId+"1")
//	assert.Contains(t, cpnts, "comp1")
//	assert.Contains(t, cpnts, "comp2")
//}
//
//// When more than one ekara.yaml file define a parent the one taken
//// in account should the the one defined in the main descriptor
//func TestDownloadFTwoParents(t *testing.T) {
//
//	mainPath := "descriptor"
//
//	c := util.CreateMockLaunchContext(mainPath, false)
//	tester := util.CreateComponentTester(t, c)
//	defer tester.Clean()
//
//	repParent1 := tester.CreateDir("parent1")
//	repParent2 := tester.CreateDir("parent2")
//	tester.CreateDirEmptyDesc("comp1")
//	tester.CreateDirEmptyDesc("comp2")
//	tester.CreateDirEmptyDesc("comp3")
//	tester.CreateDirEmptyDesc("comp4")
//	repDesc := tester.CreateDir(mainPath)
//
//	parentContent1 := `
//ekara:
// parent:
//   repository: parent2
// components:
//   comp1:
//     repository: comp1
//   comp2:
//     repository: comp2
//`
//	repParent1.WriteCommit("ekara.yaml", parentContent1)
//
//	parentContent2 := `
//ekara:
// components:
//   comp3:
//     repository: comp3
//   comp4:
//     repository: comp4
//`
//	repParent2.WriteCommit("ekara.yaml", parentContent2)
//
//	descContent := `
//name: ekaraDemoVar
//qualifier: dev
//
//ekara:
// parent:
//   repository: parent1
//
//# Following content just to force the download of comp1 and comp2
//orchestrator:
// component: comp1
//
//providers:
// p1:
//   component: comp2
//
//nodes:
// node1:
//   instances: 1
//   provider:
//     name: p1
//`
//	repDesc.WriteCommit("ekara.yaml", descContent)
//
//	err := tester.Init()
//	assert.Nil(t, err)
//
//	refM := tester.rM
//	assert.Equal(t, len(refM.usedReferences.Refs), 2)
//	assert.True(t, refM.usedReferences.IdUsed("comp1"))
//	assert.True(t, refM.usedReferences.IdUsed("comp2"))
//	assert.Equal(t, len(refM.referencedComponents.Refs), 4)
//	assert.True(t, refM.referencedComponents.IdReferenced("comp1"))
//	assert.True(t, refM.referencedComponents.IdReferenced("comp2"))
//	assert.True(t, refM.referencedComponents.IdReferenced("comp3"))
//	assert.True(t, refM.referencedComponents.IdReferenced("comp4"))
//	assert.Len(t, refM.parents, 2)
//	// Check that the parents has been renamed base on their position
//	assert.Equal(t, model.ParentComponentId+"1", refM.parents[0].comp.Id)
//	assert.Equal(t, model.ParentComponentId+"2", refM.parents[1].comp.Id)
//
//	// Check the referenced components has been cleaned
//	refM.referencedComponents.Clean(*refM.usedReferences)
//	assert.Equal(t, len(refM.referencedComponents.Refs), 2)
//	assert.True(t, refM.referencedComponents.IdReferenced("comp1"))
//	assert.True(t, refM.referencedComponents.IdReferenced("comp2"))
//
//	env := tester.Env()
//	assert.NotNil(t, env)
//	// comp1 and comp2 should be downloaded because they are used into the descriptor
//	tester.AssertComponentsContains(model.MainComponentId, model.ParentComponentId+"1", model.ParentComponentId+"2", "comp1", "comp2")
//	cpnts := env.Ekara.Components
//	assert.Equal(t, len(cpnts), 5)
//	assert.Contains(t, cpnts, "__main__")
//	assert.Contains(t, cpnts, model.ParentComponentId+"1")
//	assert.Contains(t, cpnts, model.ParentComponentId+"2")
//	assert.Contains(t, cpnts, "comp1")
//	assert.Contains(t, cpnts, "comp2")
//}

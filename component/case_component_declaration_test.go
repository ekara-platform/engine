package component

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

// A component declared into another component must be ignored
func TestComponentInComponentIgnored(t *testing.T) {

	comp2Content := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
`

	mainPath := "./testdata/gittest/descriptor"
	c := util.CreateMockLaunchContext(mainPath, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	tester.CreateRepDefaultDescriptor("./testdata/gittest/parent")

	repComp2 := tester.CreateRep("./testdata/gittest/comp2")
	repComp2.WriteCommit("ekara.yaml", comp2Content)

	repComp1 := tester.CreateRep("./testdata/gittest/comp1")
	repComp1.WriteCommit("content.txt", "comp content from parent")

	repDesc := tester.CreateRep(mainPath)

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent
  components:
    comp2:
      repository: ./testdata/gittest/comp2
# Following content just to force the download of comp1
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
	repDesc.WriteCommit("ekara.yaml", descContent)

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator
	tester.AssertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", "comp2")
	rm := tester.rM
	assert.NotNil(t, rm)

	assert.True(t, rm.referencedComponents.IdReferenced("comp2"))
	//Comp1 must not be referenced because it's declared into another component
	assert.False(t, rm.referencedComponents.IdReferenced("comp1"))

	assert.True(t, rm.usedReferences.IdUsed("comp2"))
	//Comp1 must be referenced because it's used into the main descriptor
	assert.True(t, rm.usedReferences.IdUsed("comp1"))

}

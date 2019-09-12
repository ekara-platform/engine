package engine

import (
	"testing"

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
	c := &MockLaunchContext{locationContent: mainPath, data: model.Parameters{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	tester.createRepDefaultDescriptor(t, "./testdata/gittest/parent")

	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repComp2.writeCommit(t, "ekara.yaml", comp2Content)

	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp1.writeCommit(t, "content.txt", "comp content from parent")

	repDesc := tester.createRep(mainPath)

	descContent := `
name: ekara-demo-var
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
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	// comp1 should be downloaded because it's used as orchestrator
	tester.assertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", "comp2")
	rm := c.Ekara().ReferenceManager()
	assert.NotNil(t, rm)

	assert.True(t, rm.ReferencedComponents.IdReferenced("comp2"))
	//Comp1 must not be referenced because it's declared into another component
	assert.False(t, rm.ReferencedComponents.IdReferenced("comp1"))

	assert.True(t, rm.UsedReferences.IdUsed("comp2"))
	//Comp1 must be referenced because it's used into the main descriptor
	assert.True(t, rm.UsedReferences.IdUsed("comp1"))

}

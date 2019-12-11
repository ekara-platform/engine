package component

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)
func TestTemplateOnReferences(t *testing.T) {

	p := model.CreateParameters(map[string]interface{}{
		"refParent": "./testdata/gittest/parent",
	})
	mainPath := "./testdata/gittest/descriptor"
	c := util.CreateMockLaunchContextWithData(mainPath, p, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	repParent := tester.CreateRep("./testdata/gittest/parent")
	tester.CreateRepDefaultDescriptor("./testdata/gittest/comp1")
	repDesc := tester.CreateRep(mainPath)

	

	parentContent := `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1	
`
	repParent.WriteCommit("ekara.yaml", parentContent)

	descContent := `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: "{{ .Vars.refParent }}"
providers:
  comp1:
    component: comp1
nodes:
  node1:
    instances: 1
    provider:
      name: comp1
`

	repDesc.WriteCommit("ekara.yaml", descContent)

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)

	tester.AssertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1")

}

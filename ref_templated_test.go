package engine

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/engine/model"

	"github.com/stretchr/testify/assert"
)

func TestTemplateOnReferences(t *testing.T) {

	p := model.CreateParameters(map[string]interface{}{
		"refParent": "parent",
	})
	mainPath := "descriptor"
	tester := util.CreateComponentTester(t, p)
	defer tester.Clean()

	repParent := tester.CreateDir("parent")
	tester.CreateDirEmptyDesc("comp1")
	repDesc := tester.CreateDir(mainPath)

	parentContent := `
ekara:
  components:
    comp1:
      repository: comp1	
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

	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)

	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1")

}

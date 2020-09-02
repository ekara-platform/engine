package engine

import (
	"github.com/ekara-platform/engine/model"
	"testing"

	"github.com/ekara-platform/engine/util"
	"github.com/stretchr/testify/assert"
)

var (
	refparentContent = `
ekara:
  components:
`

	refMaster = `
name: ekaraDemoVar
qualifier: master
`

	refBranch = `
name: ekaraDemoVar
qualifier: newBranch
`

	refTag1 = `
name: ekaraDemoVar
qualifier: newTag1
`
	refTag2 = `
name: ekaraDemoVar
qualifier: newTag2
`

	refDescContent = `
ekara:
  parent:
    repository: parent
`
)

func TestEngineLocalWithBranchRef(t *testing.T) {
	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	repParent := tester.CreateDir("parent")
	repParent.WriteCommit("ekara.yaml", refparentContent)

	repDesc := tester.CreateDir("descriptor")
	// Commit the master
	repDesc.WriteCommit("ekara.yaml", refMaster+refDescContent)

	// Commit the new branch
	repDesc.CreateBranch("newBranch")
	repDesc.WriteCommit("ekara.yaml", refBranch+refDescContent)
	repDesc.Checkout("master")

	tester.Init(repDesc.AsRepository("newBranch"))
	env := tester.Env()
	assert.NotNil(t, env)
	//TODO Fix this
	//assert.Equal(t, "newBranch", env.Qualifier)
}

func TestEngineLocalWithTagRef(t *testing.T) {
	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	repParent := tester.CreateDir("parent")
	repParent.WriteCommit("ekara.yaml", refparentContent)

	repDesc := tester.CreateDir("descriptor")
	// Commit the master
	repDesc.WriteCommit("ekara.yaml", refMaster+refDescContent)

	// Commit the new branch
	repDesc.CreateBranch("newBranch")
	repDesc.WriteCommit("ekara.yaml", refBranch+refDescContent)
	repDesc.Checkout("master")

	// Commit a tag content
	repDesc.WriteCommit("ekara.yaml", refTag1+refDescContent)
	// Creating the tag
	repDesc.Tag("newTag1")
	// Content has change after the tag creation
	repDesc.WriteCommit("ekara.yaml", refTag2+refDescContent)
	tester.Init(repDesc.AsRepository("newTag1"))
	env := tester.Env()
	assert.NotNil(t, env)
	//TODO Fix this
	//assert.Equal(t, "newTag1", env.Qualifier)
}

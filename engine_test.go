package engine

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/engine/component"
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
    repository: ./testdata/gittest/parent
`
)

func TestEngineLocalWithBranchRef(t *testing.T) {
	c := util.CreateMockLaunchContext("./testdata/gittest/descriptor@newBranch", false)
	tester := component.CreateComponentTester(t, c)
	defer tester.Clean()

	repParent := tester.CreateRep("./testdata/gittest/parent")
	repParent.WriteCommit("ekara.yaml", refparentContent)

	repDesc := tester.CreateRep("./testdata/gittest/descriptor")
	// Commit the master
	repDesc.WriteCommit("ekara.yaml", refMaster+refDescContent)

	// Commit the new branch
	repDesc.CreateBranch("newBranch")
	repDesc.WriteCommit("ekara.yaml", refBranch+refDescContent)
	repDesc.Checkout("master")

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)
	//TODO Fix this
	//assert.Equal(t, "newBranch", env.Qualifier)
}

func TestEngineLocalWithTagRef(t *testing.T) {

	c := util.CreateMockLaunchContext("./testdata/gittest/descriptor@newTag1", false)
	tester := component.CreateComponentTester(t, c)
	defer tester.Clean()

	repParent := tester.CreateRep("./testdata/gittest/parent")
	repParent.WriteCommit("ekara.yaml", refparentContent)

	repDesc := tester.CreateRep("./testdata/gittest/descriptor")
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

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)
	//TODO Fix this
	//assert.Equal(t, "newTag1", env.Qualifier)
}

func TestRepositoryFlavor(t *testing.T) {

	a, b := util.RepositoryFlavor("aaa")
	assert.Equal(t, a, "aaa")
	assert.Equal(t, b, "")

	a, b = util.RepositoryFlavor("aaa@bbb")
	assert.Equal(t, a, "aaa")
	assert.Equal(t, b, "bbb")

	a, b = util.RepositoryFlavor("aaa@")
	assert.Equal(t, a, "aaa")
	assert.Equal(t, b, "")

	a, b = util.RepositoryFlavor("aaa@bbb@willbeignored")
	assert.Equal(t, a, "aaa")
	assert.Equal(t, b, "bbb")
}

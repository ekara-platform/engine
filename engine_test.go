package engine

import (
	"testing"

	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"
)

var (
	refDistContent = `
ekara:
  components:
`

	refMaster = `
name: ekara-demo-var
qualifier: master
`

	refBranch = `
name: ekara-demo-var
qualifier: newBranch
`

	refTag1 = `
name: ekara-demo-var
qualifier: newTag1
`
	refTag2 = `
name: ekara-demo-var
qualifier: newTag2
`

	refDescContent = `
ekara:
  parent:
    repository: ./testdata/gittest/parent
`
)

func TestEngineLocalWithBranchRef(t *testing.T) {

	c := &MockLaunchContext{locationContent: "./testdata/gittest/descriptor@newBranch", templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/parent")
	repDist.writeCommit(t, "ekara.yaml", refDistContent)

	repDesc := tester.createRep("./testdata/gittest/descriptor")
	// Commit the master
	repDesc.writeCommit(t, "ekara.yaml", refMaster+refDescContent)

	// Commit the new branch
	repDesc.createBranch(t, "newBranch")
	repDesc.writeCommit(t, "ekara.yaml", refBranch+refDescContent)
	repDesc.checkout(t, "master")

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	//TODO Fix this
	//assert.Equal(t, "newBranch", env.Qualifier)
}

func TestEngineLocalWithTagRef(t *testing.T) {

	c := &MockLaunchContext{locationContent: "./testdata/gittest/descriptor@newTag1", templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/parent")
	repDist.writeCommit(t, "ekara.yaml", refDistContent)

	repDesc := tester.createRep("./testdata/gittest/descriptor")
	// Commit the master
	repDesc.writeCommit(t, "ekara.yaml", refMaster+refDescContent)

	// Commit the new branch
	repDesc.createBranch(t, "newBranch")
	repDesc.writeCommit(t, "ekara.yaml", refBranch+refDescContent)
	repDesc.checkout(t, "master")

	// Commit a tag content
	repDesc.writeCommit(t, "ekara.yaml", refTag1+refDescContent)
	// Creating the tag
	repDesc.tag(t, "newTag1")
	// Content has change after the tag creation
	repDesc.writeCommit(t, "ekara.yaml", refTag2+refDescContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	//TODO Fix this
	//assert.Equal(t, "newTag1", env.Qualifier)
}

func TestRepositoryFlavor(t *testing.T) {

	a, b := repositoryFlavor("aaa")
	assert.Equal(t, a, "aaa")
	assert.Equal(t, b, "")

	a, b = repositoryFlavor("aaa@bbb")
	assert.Equal(t, a, "aaa")
	assert.Equal(t, b, "bbb")

	a, b = repositoryFlavor("aaa@")
	assert.Equal(t, a, "aaa")
	assert.Equal(t, b, "")

	a, b = repositoryFlavor("aaa@bbb@willbeignored")
	assert.Equal(t, a, "aaa")
	assert.Equal(t, b, "bbb")
}

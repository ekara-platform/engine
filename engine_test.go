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
  distribution:
    repository: ./testdata/gittest/distribution
`
)

/*
func TestEngineLocalNoRef(t *testing.T) {
	mainPath := "./testdata/gittest/descriptor"

	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repDist.writeCommit(t, "ekara.yaml", refDistContent)

	repDesc := tester.createRep(mainPath)
	// Commit the master
	repDesc.writeCommit(t, "ekara.yaml", refMaster+refDescContent)

	// Commit the new branch
	repDesc.createBranch(t, "newBranch")
	repDesc.writeCommit(t, "ekara.yaml", refBranch+refDescContent)
	repDesc.checkout(t, "master")

	c.locationContent = mainPath
	a, b := repositoryFlavor(c.locationContent)
	assert.Equal(t, a, mainPath)
	assert.Equal(t, b, "")
	err := tester.initEngine()
	assert.Nil(t, err)
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	assert.Equal(t, env.Qualifier, "master")

}
*/

func TestEngineLocalWithBranchRef(t *testing.T) {

	c := &MockLaunchContext{locationContent: "./testdata/gittest/descriptor@newBranch", templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
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
	err = tester.context.engine.ComponentManager().Ensure()
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

	repDist := tester.createRep("./testdata/gittest/distribution")
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
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	//TODO Fix this
	//assert.Equal(t, "newTag1", env.Qualifier)
}

/*
func TestEngineLocalWithTagRef(t *testing.T) {
	engine := createTestEngine()
	c := MockLaunchContext{locationContent: "testdata/sample@v1.0.0"}
	e := engine.Init(c)
	assertOnlyWarnings(t, e)
}

func TestEngineLocalWithRawRef(t *testing.T) {
	engine := createTestEngine()
	c := MockLaunchContext{locationContent: "testdata/sample@refs/remotes/origin/test"}
	e := engine.Init(c)
	assertOnlyWarnings(t, e)
}

*/
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

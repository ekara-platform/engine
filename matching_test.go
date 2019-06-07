package engine

import (
	"path/filepath"
	"testing"

	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

var (
	matchDistContent = `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
    comp2:
      repository: ./testdata/gittest/comp2
`

	matchDescContent = `
name: ekara-demo-var
qualifier: dev

ekara:
  distribution:
    repository: ./testdata/gittest/distribution

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
)

func TestComponentFolderMatching(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

	repComp1.writeCommit(t, "ekara.yaml", ``)
	repComp1.writeFolderCommit(t, "wantedfolder1", "test.yaml", `test content`)

	repComp2.writeCommit(t, "ekara.yaml", ``)
	repComp2.writeFolderCommit(t, "wantedfolder1", "test.yaml", `test content`)
	repComp2.writeFolderCommit(t, "wantedfolder2", "test.yaml", `test content`)
	repComp2.writeFolderCommit(t, "wantedfolder2/subFolder1/subfolder2", "test.yaml", `test content`)
	repComp2.writeFolderCommit(t, "wantedfolder3", "fileSearchedAsFolder.yaml", `test content`)

	repDist.writeCommit(t, "ekara.yaml", matchDistContent)
	repDesc.writeCommit(t, "ekara.yaml", matchDescContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	
	tester.assertComponentsContains("__main__", "__ekara__", "comp1", "comp2")

	cm := c.Ekara().ComponentManager()
	assert.NotNil(t, cm)

	valP1Comp2, ok := env.Providers["p1"]
	assert.True(t, ok)

	usableComp1 := cm.Use(env.Orchestrator)
	usableComp2 := cm.Use(valP1Comp2)

	//----------------------------------------------------------
	// Matching againts a given component
	//----------------------------------------------------------

	ok, match := usableComp1.ContainsDirectory("wantedfolder1")
	if assert.True(t, ok) {
		checkMatchingPath(t, match, usableComp1, "wantedfolder1")
	}

	ok, match = usableComp1.ContainsDirectory("wantedfolder2")
	assert.False(t, ok)

	ok, match = usableComp1.ContainsDirectory("wantedfolder2/subFolder1/subfolder2")
	assert.False(t, ok)

	ok, match = usableComp2.ContainsDirectory("wantedfolder1")
	if assert.True(t, ok) {
		checkMatchingPath(t, match, usableComp2, "wantedfolder1")
	}

	ok, match = usableComp2.ContainsDirectory("wantedfolder2")
	if assert.True(t, ok) {
		checkMatchingPath(t, match, usableComp2, "wantedfolder2")
	}

	ok, match = usableComp2.ContainsDirectory("wantedfolder2/subFolder1/subfolder2")
	if assert.True(t, ok) {
		checkMatchingPath(t, match, usableComp2, "wantedfolder2/subFolder1/subfolder2")
	}

	ok, match = usableComp2.ContainsDirectory("wantedfolder3/fileSearchedAsFolder.yaml")
	assert.False(t, ok)

	//----------------------------------------------------------
	// Matching against all components through the component manager
	//----------------------------------------------------------
	paths := cm.ContainsDirectory("wantedfolder1")
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedfolder1"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder1"))
	}

	paths = cm.ContainsDirectory("wantedfolder2")
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder2"))
	}

	paths = cm.ContainsDirectory("wantedfolder2/subFolder1/subfolder2")
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder2/subFolder1/subfolder2"))
	}

	paths = cm.ContainsDirectory("wantedfolder3/fileSearchedAsFolder.yaml")
	assert.Equal(t, 0, paths.Count())

	paths = cm.ContainsDirectory("mising")
	assert.Equal(t, 0, paths.Count())

	//----------------------------------------------------------
	// Matching against restricted components through the component manager
	//----------------------------------------------------------
	paths = cm.ContainsDirectory("wantedfolder1", env.Orchestrator, valP1Comp2)
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedfolder1"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder1"))
	}

	paths = cm.ContainsDirectory("wantedfolder1", valP1Comp2)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder1"))
	}

	paths = cm.ContainsDirectory("wantedfolder1", env.Orchestrator)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedfolder1"))
	}

}

func TestComponentFileMatching(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	c := &MockLaunchContext{locationContent: mainPath, templateContext: &model.TemplateContext{}}
	tester := gitTester(t, c)
	defer tester.clean()

	repDist := tester.createRep("./testdata/gittest/distribution")
	repComp1 := tester.createRep("./testdata/gittest/comp1")
	repComp2 := tester.createRep("./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

	// Files in component 1
	repComp1.writeCommit(t, "ekara.yaml", ``)
	repComp1.writeCommit(t, "wantedFile1.txt", `test content`)

	// Files in component 2
	repComp2.writeCommit(t, "ekara.yaml", ``)
	repComp2.writeCommit(t, "wantedFile1.txt", `test content`)
	repComp2.writeFolderCommit(t, "subfolder", "wantedSubFile1.txt", `test content`)
	repComp2.writeCommit(t, "wantedFile2.txt", `test content`)
	repComp2.writeFolderCommit(t, "folderSearchedAsFile.txt", "test.yaml", `test content`)

	// Files in distribution
	repDist.writeCommit(t, "ekara.yaml", matchDistContent)
	repDesc.writeCommit(t, "ekara.yaml", matchDescContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)
	
	tester.assertComponentsContains("__main__", "__ekara__", "comp1", "comp2")

	cm := c.Ekara().ComponentManager()
	assert.NotNil(t, cm)

	valP1Comp2, ok := env.Providers["p1"]
	assert.True(t, ok)

	usableComp1 := cm.Use(env.Orchestrator)
	usableComp2 := cm.Use(valP1Comp2)

	//----------------------------------------------------------
	// Matching againts a given component
	//----------------------------------------------------------

	ok, match := usableComp1.ContainsFile("wantedFile1.txt")
	if assert.True(t, ok) {
		checkMatchingPath(t, match, usableComp1, "wantedFile1.txt")
	}

	ok, match = usableComp1.ContainsFile("wantedFile2.txt")
	assert.False(t, ok)

	ok, match = usableComp1.ContainsFile("subfolder/wantedSubFile1.txt")
	assert.False(t, ok)

	ok, match = usableComp2.ContainsFile("wantedFile1.txt")
	if assert.True(t, ok) {
		checkMatchingPath(t, match, usableComp2, "wantedFile1.txt")
	}

	ok, match = usableComp2.ContainsFile("wantedFile2.txt")
	if assert.True(t, ok) {
		checkMatchingPath(t, match, usableComp2, "wantedFile2.txt")
	}

	ok, match = usableComp2.ContainsFile("subfolder/wantedSubFile1.txt")
	if assert.True(t, ok) {
		checkMatchingPath(t, match, usableComp2, "subfolder/wantedSubFile1.txt")
	}

	ok, match = usableComp2.ContainsFile("folderSearchedAsFile.txt")
	assert.False(t, ok)

	//----------------------------------------------------------
	// Matching against all components through the component manager
	//----------------------------------------------------------

	paths := cm.ContainsFile("wantedFile1.txt")
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedFile1.txt"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile1.txt"))
	}

	paths = cm.ContainsFile("wantedFile2.txt")
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile2.txt"))
	}

	paths = cm.ContainsFile("subfolder/wantedSubFile1.txt")
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "subfolder/wantedSubFile1.txt"))
	}

	paths = cm.ContainsFile("folderSearchedAsFile.txt")
	assert.Equal(t, 0, paths.Count())

	paths = cm.ContainsFile("mising")
	assert.Equal(t, 0, paths.Count())

	//----------------------------------------------------------
	// Matching against restricted components through the component manager
	//----------------------------------------------------------

	paths = cm.ContainsFile("wantedFile1.txt", env.Orchestrator, valP1Comp2)
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedFile1.txt"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile1.txt"))
	}

	paths = cm.ContainsFile("wantedFile1.txt", valP1Comp2)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile1.txt"))
	}

	paths = cm.ContainsFile("wantedFile1.txt", env.Orchestrator)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedFile1.txt"))
	}
}

func mergeMatchingPaths(matches component.MatchingPaths) []string {
	res := make([]string, 0, 0)
	for _, v := range matches.Paths {
		res = append(res, filepath.Join(v.Component().RootPath(), v.RelativePath()))
	}
	return res
}

func checkMatchingPath(t *testing.T, match component.MatchingPath, u component.UsableComponent, relative string) {
	assert.Equal(t, match.Component().RootPath(), u.RootPath())
	assert.Equal(t, match.RelativePath(), relative)
}

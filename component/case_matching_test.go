package component

import (
	"path/filepath"
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

var (
	matchParentContent = `
ekara:
  components:
    comp1:
      repository: ./testdata/gittest/comp1
    comp2:
      repository: ./testdata/gittest/comp2
`

	matchDescContent = `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: ./testdata/gittest/parent

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

	c := util.CreateMockLaunchContext(mainPath, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	repParent := tester.CreateRep("./testdata/gittest/parent")
	repComp1 := tester.CreateRepDefaultDescriptor("./testdata/gittest/comp1")
	repComp2 := tester.CreateRepDefaultDescriptor("./testdata/gittest/comp2")
	repDesc := tester.CreateRep(mainPath)

	repComp1.WriteFolderCommit("wantedfolder1", "test.yaml", `test content`)

	repComp2.WriteFolderCommit("wantedfolder1", "test.yaml", `test content`)
	repComp2.WriteFolderCommit("wantedfolder2", "test.yaml", `test content`)
	repComp2.WriteFolderCommit("wantedfolder2/subFolder1/subfolder2", "test.yaml", `test content`)
	repComp2.WriteFolderCommit("wantedfolder3", "fileSearchedAsFolder.yaml", `test content`)

	repParent.WriteCommit("ekara.yaml", matchParentContent)
	repDesc.WriteCommit("ekara.yaml", matchDescContent)

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)

	tester.AssertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1", "comp2")

	cF := tester.cF
	assert.NotNil(t, cF)

	valP1Comp2, ok := env.Providers["p1"]
	assert.True(t, ok)

	usableComp1, err := cF.Use(env.Orchestrator, *tester.tplC)
	assert.Nil(t, err)
	usableComp2, err := cF.Use(valP1Comp2, *tester.tplC)
	assert.Nil(t, err)

	//----------------------------------------------------------
	// Matching against a given component
	//----------------------------------------------------------

	ok, match := usableComp1.ContainsDirectory("wantedfolder1")
	if assert.True(t, ok) {
		checkMatchingPath(t, match, usableComp1, "wantedfolder1")
	}

	ok, _ = usableComp1.ContainsDirectory("wantedfolder2")
	assert.False(t, ok)

	ok, _ = usableComp1.ContainsDirectory("wantedfolder2/subFolder1/subfolder2")
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

	ok, _ = usableComp2.ContainsDirectory("wantedfolder3/fileSearchedAsFolder.yaml")
	assert.False(t, ok)

	//----------------------------------------------------------
	// Matching against all components through the component manager
	//----------------------------------------------------------
	paths := cF.ContainsDirectory("wantedfolder1", *tester.tplC)
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedfolder1"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder1"))
	}

	paths = cF.ContainsDirectory("wantedfolder2", *tester.tplC)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder2"))
	}

	paths = cF.ContainsDirectory("wantedfolder2/subFolder1/subfolder2", *tester.tplC)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder2/subFolder1/subfolder2"))
	}

	paths = cF.ContainsDirectory("wantedfolder3/fileSearchedAsFolder.yaml", *tester.tplC)
	assert.Equal(t, 0, paths.Count())

	paths = cF.ContainsDirectory("mising", *tester.tplC)
	assert.Equal(t, 0, paths.Count())

	//----------------------------------------------------------
	// Matching against restricted components through the component manager
	//----------------------------------------------------------
	paths = cF.ContainsDirectory("wantedfolder1", *tester.tplC, env.Orchestrator, valP1Comp2)
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedfolder1"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder1"))
	}

	paths = cF.ContainsDirectory("wantedfolder1", *tester.tplC, valP1Comp2)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder1"))
	}

	paths = cF.ContainsDirectory("wantedfolder1", *tester.tplC, env.Orchestrator)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedfolder1"))
	}

}

func TestComponentFileMatching(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	c := util.CreateMockLaunchContext(mainPath, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	repParent := tester.CreateRep("./testdata/gittest/parent")
	repComp1 := tester.CreateRepDefaultDescriptor("./testdata/gittest/comp1")
	repComp2 := tester.CreateRepDefaultDescriptor("./testdata/gittest/comp2")
	repDesc := tester.CreateRep(mainPath)

	// Files in component 1
	repComp1.WriteCommit("wantedFile1.txt", `test content`)

	// Files in component 2
	repComp2.WriteCommit("wantedFile1.txt", `test content`)
	repComp2.WriteFolderCommit("subfolder", "wantedSubFile1.txt", `test content`)
	repComp2.WriteCommit("wantedFile2.txt", `test content`)
	repComp2.WriteFolderCommit("folderSearchedAsFile.txt", "test.yaml", `test content`)

	// Files in parent
	repParent.WriteCommit("ekara.yaml", matchParentContent)
	repDesc.WriteCommit("ekara.yaml", matchDescContent)

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)

	tester.AssertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1", "comp2")

	cF := tester.cF
	assert.NotNil(t, cF)

	valP1Comp2, ok := env.Providers["p1"]
	assert.True(t, ok)

	usableComp1, err := cF.Use(env.Orchestrator, *tester.tplC)
	assert.Nil(t, err)
	usableComp2, err := cF.Use(valP1Comp2, *tester.tplC)
	assert.Nil(t, err)

	//----------------------------------------------------------
	// Matching against a given component
	//----------------------------------------------------------

	ok, match := usableComp1.ContainsFile("wantedFile1.txt")
	if assert.True(t, ok) {
		checkMatchingPath(t, match, usableComp1, "wantedFile1.txt")
	}

	ok, _ = usableComp1.ContainsFile("wantedFile2.txt")
	assert.False(t, ok)

	ok, _ = usableComp1.ContainsFile("subfolder/wantedSubFile1.txt")
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

	ok, _ = usableComp2.ContainsFile("folderSearchedAsFile.txt")
	assert.False(t, ok)

	//----------------------------------------------------------
	// Matching against all components through the component manager
	//----------------------------------------------------------

	paths := cF.ContainsFile("wantedFile1.txt", *tester.tplC)
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedFile1.txt"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile1.txt"))
	}

	paths = cF.ContainsFile("wantedFile2.txt", *tester.tplC)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile2.txt"))
	}

	paths = cF.ContainsFile("subfolder/wantedSubFile1.txt", *tester.tplC)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "subfolder/wantedSubFile1.txt"))
	}

	paths = cF.ContainsFile("folderSearchedAsFile.txt", *tester.tplC)
	assert.Equal(t, 0, paths.Count())

	paths = cF.ContainsFile("mising", *tester.tplC)
	assert.Equal(t, 0, paths.Count())

	//----------------------------------------------------------
	// Matching against restricted components through the component manager
	//----------------------------------------------------------

	paths = cF.ContainsFile("wantedFile1.txt", *tester.tplC, env.Orchestrator, valP1Comp2)
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedFile1.txt"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile1.txt"))
	}

	paths = cF.ContainsFile("wantedFile1.txt", *tester.tplC, valP1Comp2)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile1.txt"))
	}

	paths = cF.ContainsFile("wantedFile1.txt", *tester.tplC, env.Orchestrator)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedFile1.txt"))
	}
}

func mergeMatchingPaths(matches MatchingPaths) []string {
	res := make([]string, 0, 0)
	for _, v := range matches.Paths {
		res = append(res, filepath.Join(v.Component().RootPath(), v.RelativePath()))
	}
	return res
}

func checkMatchingPath(t *testing.T, match MatchingPath, u UsableComponent, relative string) {
	assert.Equal(t, match.Component().RootPath(), u.RootPath())
	assert.Equal(t, match.RelativePath(), relative)
}

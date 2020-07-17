package engine

import (
	"github.com/GroupePSA/componentizer"
	"path/filepath"
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/engine/model"

	"github.com/stretchr/testify/assert"
)

var (
	matchParentContent = `
ekara:
  components:
    comp1:
      repository: comp1
    comp2:
      repository: comp2
`

	matchDescContent = `
name: ekaraDemoVar
qualifier: dev

ekara:
  parent:
    repository: parent

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
	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	repParent := tester.CreateDir("parent")
	repComp1 := tester.CreateDirEmptyDesc("comp1")
	repComp2 := tester.CreateDirEmptyDesc("comp2")
	repDesc := tester.CreateDir("descriptor")

	repComp1.WriteFolderCommit("wantedfolder1", "test.yaml", `test content`)

	repComp2.WriteFolderCommit("wantedfolder1", "test.yaml", `test content`)
	repComp2.WriteFolderCommit("wantedfolder2", "test.yaml", `test content`)
	repComp2.WriteFolderCommit("wantedfolder2/subFolder1/subfolder2", "test.yaml", `test content`)
	repComp2.WriteFolderCommit("wantedfolder3", "fileSearchedAsFolder.yaml", `test content`)

	repParent.WriteCommit("ekara.yaml", matchParentContent)
	repDesc.WriteCommit("ekara.yaml", matchDescContent)

	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)

	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1", "comp2")

	cM := tester.ComponentManager()
	assert.NotNil(t, cM)

	valP1Comp2, ok := env.Providers["p1"]
	assert.True(t, ok)

	usableComp1, err := cM.Use(env.Orchestrator, tester.TemplateContext())
	assert.Nil(t, err)
	usableComp2, err := cM.Use(valP1Comp2, tester.TemplateContext())
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
	paths := cM.ContainsDirectory("wantedfolder1", tester.TemplateContext())
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedfolder1"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder1"))
	}

	paths = cM.ContainsDirectory("wantedfolder2", tester.TemplateContext())
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder2"))
	}

	paths = cM.ContainsDirectory("wantedfolder2/subFolder1/subfolder2", tester.TemplateContext())
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder2/subFolder1/subfolder2"))
	}

	paths = cM.ContainsDirectory("wantedfolder3/fileSearchedAsFolder.yaml", tester.TemplateContext())
	assert.Equal(t, 0, paths.Count())

	paths = cM.ContainsDirectory("mising", tester.TemplateContext())
	assert.Equal(t, 0, paths.Count())

	//----------------------------------------------------------
	// Matching against restricted components through the component manager
	//----------------------------------------------------------
	paths = cM.ContainsDirectory("wantedfolder1", tester.TemplateContext(), env.Orchestrator, valP1Comp2)
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedfolder1"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder1"))
	}

	paths = cM.ContainsDirectory("wantedfolder1", tester.TemplateContext(), valP1Comp2)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder1"))
	}

	paths = cM.ContainsDirectory("wantedfolder1", tester.TemplateContext(), env.Orchestrator)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedfolder1"))
	}

}

func TestComponentFileMatching(t *testing.T) {
	tester := util.CreateComponentTester(t, model.CreateEmptyParameters())
	defer tester.Clean()

	repParent := tester.CreateDir("parent")
	repComp1 := tester.CreateDirEmptyDesc("comp1")
	repComp2 := tester.CreateDirEmptyDesc("comp2")
	repDesc := tester.CreateDir("descriptor")

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

	tester.Init(repDesc.AsRepository("master"))
	env := tester.Env()
	assert.NotNil(t, env)

	tester.AssertComponentsExactly(model.MainComponentId, model.MainComponentId+model.ParentComponentSuffix, "comp1", "comp2")

	cM := tester.ComponentManager()
	assert.NotNil(t, cM)

	valP1Comp2, ok := env.Providers["p1"]
	assert.True(t, ok)

	usableComp1, err := cM.Use(env.Orchestrator, tester.TemplateContext())
	assert.Nil(t, err)
	usableComp2, err := cM.Use(valP1Comp2, tester.TemplateContext())
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

	paths := cM.ContainsFile("wantedFile1.txt", tester.TemplateContext())
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedFile1.txt"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile1.txt"))
	}

	paths = cM.ContainsFile("wantedFile2.txt", tester.TemplateContext())
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile2.txt"))
	}

	paths = cM.ContainsFile("subfolder/wantedSubFile1.txt", tester.TemplateContext())
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "subfolder/wantedSubFile1.txt"))
	}

	paths = cM.ContainsFile("folderSearchedAsFile.txt", tester.TemplateContext())
	assert.Equal(t, 0, paths.Count())

	paths = cM.ContainsFile("mising", tester.TemplateContext())
	assert.Equal(t, 0, paths.Count())

	//----------------------------------------------------------
	// Matching against restricted components through the component manager
	//----------------------------------------------------------

	paths = cM.ContainsFile("wantedFile1.txt", tester.TemplateContext(), env.Orchestrator, valP1Comp2)
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedFile1.txt"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile1.txt"))
	}

	paths = cM.ContainsFile("wantedFile1.txt", tester.TemplateContext(), valP1Comp2)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile1.txt"))
	}

	paths = cM.ContainsFile("wantedFile1.txt", tester.TemplateContext(), env.Orchestrator)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedFile1.txt"))
	}
}

func mergeMatchingPaths(matches componentizer.MatchingPaths) []string {
	res := make([]string, 0, 0)
	for _, v := range matches.Paths {
		res = append(res, filepath.Join(v.Owner().RootPath(), v.RelativePath()))
	}
	return res
}

func checkMatchingPath(t *testing.T, match componentizer.MatchingPath, u componentizer.UsableComponent, relative string) {
	assert.Equal(t, match.Owner().RootPath(), u.RootPath())
	assert.Equal(t, match.RelativePath(), relative)
}

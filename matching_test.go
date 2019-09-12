package engine

import (
	"path/filepath"
	"testing"

	"github.com/ekara-platform/engine/component"
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
name: ekara-demo-var
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

	c := &MockLaunchContext{locationContent: mainPath, data: model.Parameters{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repParent := tester.createRep("./testdata/gittest/parent")
	repComp1 := tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp1")
	repComp2 := tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

	repComp1.writeFolderCommit(t, "wantedfolder1", "test.yaml", `test content`)

	repComp2.writeFolderCommit(t, "wantedfolder1", "test.yaml", `test content`)
	repComp2.writeFolderCommit(t, "wantedfolder2", "test.yaml", `test content`)
	repComp2.writeFolderCommit(t, "wantedfolder2/subFolder1/subfolder2", "test.yaml", `test content`)
	repComp2.writeFolderCommit(t, "wantedfolder3", "fileSearchedAsFolder.yaml", `test content`)

	repParent.writeCommit(t, "ekara.yaml", matchParentContent)
	repDesc.writeCommit(t, "ekara.yaml", matchDescContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	tester.assertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1", "comp2")

	cm := c.Ekara().ComponentManager()
	assert.NotNil(t, cm)

	valP1Comp2, ok := env.Providers["p1"]
	assert.True(t, ok)

	usableComp1, err := cm.Use(env.Orchestrator, &model.TemplateContext{})
	assert.Nil(t, err)
	usableComp2, err := cm.Use(valP1Comp2, &model.TemplateContext{})
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
	paths := cm.ContainsDirectory("wantedfolder1", &model.TemplateContext{})
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedfolder1"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder1"))
	}

	paths = cm.ContainsDirectory("wantedfolder2", &model.TemplateContext{})
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder2"))
	}

	paths = cm.ContainsDirectory("wantedfolder2/subFolder1/subfolder2", &model.TemplateContext{})
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder2/subFolder1/subfolder2"))
	}

	paths = cm.ContainsDirectory("wantedfolder3/fileSearchedAsFolder.yaml", &model.TemplateContext{})
	assert.Equal(t, 0, paths.Count())

	paths = cm.ContainsDirectory("mising", &model.TemplateContext{})
	assert.Equal(t, 0, paths.Count())

	//----------------------------------------------------------
	// Matching against restricted components through the component manager
	//----------------------------------------------------------
	paths = cm.ContainsDirectory("wantedfolder1", &model.TemplateContext{}, env.Orchestrator, valP1Comp2)
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedfolder1"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder1"))
	}

	paths = cm.ContainsDirectory("wantedfolder1", &model.TemplateContext{}, valP1Comp2)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedfolder1"))
	}

	paths = cm.ContainsDirectory("wantedfolder1", &model.TemplateContext{}, env.Orchestrator)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedfolder1"))
	}

}

func TestComponentFileMatching(t *testing.T) {

	mainPath := "./testdata/gittest/descriptor"

	c := &MockLaunchContext{locationContent: mainPath, data: model.Parameters{}}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repParent := tester.createRep("./testdata/gittest/parent")
	repComp1 := tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp1")
	repComp2 := tester.createRepDefaultDescriptor(t, "./testdata/gittest/comp2")
	repDesc := tester.createRep(mainPath)

	// Files in component 1
	repComp1.writeCommit(t, "wantedFile1.txt", `test content`)

	// Files in component 2
	repComp2.writeCommit(t, "wantedFile1.txt", `test content`)
	repComp2.writeFolderCommit(t, "subfolder", "wantedSubFile1.txt", `test content`)
	repComp2.writeCommit(t, "wantedFile2.txt", `test content`)
	repComp2.writeFolderCommit(t, "folderSearchedAsFile.txt", "test.yaml", `test content`)

	// Files in parent
	repParent.writeCommit(t, "ekara.yaml", matchParentContent)
	repDesc.writeCommit(t, "ekara.yaml", matchDescContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	tester.assertComponentsContains(model.MainComponentId, model.EkaraComponentId+"1", "comp1", "comp2")

	cm := c.Ekara().ComponentManager()
	assert.NotNil(t, cm)

	valP1Comp2, ok := env.Providers["p1"]
	assert.True(t, ok)

	usableComp1, err := cm.Use(env.Orchestrator, &model.TemplateContext{})
	assert.Nil(t, err)
	usableComp2, err := cm.Use(valP1Comp2, &model.TemplateContext{})
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

	paths := cm.ContainsFile("wantedFile1.txt", &model.TemplateContext{})
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedFile1.txt"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile1.txt"))
	}

	paths = cm.ContainsFile("wantedFile2.txt", &model.TemplateContext{})
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile2.txt"))
	}

	paths = cm.ContainsFile("subfolder/wantedSubFile1.txt", &model.TemplateContext{})
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "subfolder/wantedSubFile1.txt"))
	}

	paths = cm.ContainsFile("folderSearchedAsFile.txt", &model.TemplateContext{})
	assert.Equal(t, 0, paths.Count())

	paths = cm.ContainsFile("mising", &model.TemplateContext{})
	assert.Equal(t, 0, paths.Count())

	//----------------------------------------------------------
	// Matching against restricted components through the component manager
	//----------------------------------------------------------

	paths = cm.ContainsFile("wantedFile1.txt", &model.TemplateContext{}, env.Orchestrator, valP1Comp2)
	if assert.Equal(t, 2, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp1.RootPath(), "wantedFile1.txt"))
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile1.txt"))
	}

	paths = cm.ContainsFile("wantedFile1.txt", &model.TemplateContext{}, valP1Comp2)
	if assert.Equal(t, 1, paths.Count()) {
		located := mergeMatchingPaths(paths)
		assert.Contains(t, located, filepath.Join(usableComp2.RootPath(), "wantedFile1.txt"))
	}

	paths = cm.ContainsFile("wantedFile1.txt", &model.TemplateContext{}, env.Orchestrator)
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

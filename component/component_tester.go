package component

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type (
	ComponentTester struct {
		workdir string
		paths   []string
		context util.LaunchContext
		t       *testing.T
		cM      *manager
	}

	testRepo struct {
		path     string
		rep      *git.Repository
		lastHash plumbing.Hash
	}
)

func CreateComponentTester(t *testing.T, lC util.LaunchContext) *ComponentTester {
	tester := &ComponentTester{
		workdir: "testdata/gitwork",
		context: lC,
		t:       t,
		paths:   make([]string, 0, 0),
	}
	tester.Clean()
	tester.cM = CreateComponentManager(lC, tester.workdir).(*manager)
	return tester
}

func (t *ComponentTester) Clean() {
	os.RemoveAll("./testdata/gittest/")
	os.RemoveAll(t.workdir)
	for _, v := range t.paths {
		os.RemoveAll(v)
	}
}

func (t *ComponentTester) Init() error {
	wdURL, err := model.GetCurrentDirectoryURL(t.context.Log())
	if err != nil {
		return err
	}
	repo, ref := util.RepositoryFlavor(t.context.Location())
	mainRep, err := model.CreateRepository(model.Base{Url: wdURL}, repo, ref, t.context.DescriptorName())
	if err != nil {
		return err
	}
	mainComponent := model.CreateComponent(model.MainComponentId, mainRep)
	err = t.cM.Init(mainComponent)
	if err != nil {
		return err
	}
	return t.cM.Ensure()
	// TODO Link here the generated interface with the context
	// ONce this is done into the engine then remove this link.
	//t.context.templateContext.Model = t.Env()
}

func (t *ComponentTester) Env() model.Environment {
	return *t.cM.Environment()
}

func (t *ComponentTester) AssertComponentsContainsExactly(paths ...string) bool {
	files, _ := ioutil.ReadDir(t.workdir + "/components/")
	assert.Equal(t.t, len(paths), len(files))
	return t.AssertComponentsContains(paths...)
}

func (t *ComponentTester) AssertComponentsContains(paths ...string) bool {
	for _, v := range paths {
		_, err := os.Stat(filepath.Join(t.workdir, "components", v))
		if err == nil {
			continue
		}
		log.Println(err.Error())
		if os.IsNotExist(err) {
			return assert.Fail(t.t, fmt.Sprintf("workdir doesn't contains %s", v))
		}
	}
	return true
}

func (t *ComponentTester) CreateRepDefaultDescriptor(path string) *testRepo {
	t.paths = append(t.paths, path)
	rep, err := git.PlainInit(path, false)
	assert.NotNil(t.t, rep)
	assert.Nil(t.t, err)
	res := &testRepo{
		path: path,
		rep:  rep,
	}
	res.WriteCommit("ekara.yaml", ``)
	return res
}

func (t *ComponentTester) CreateRep(path string) *testRepo {
	t.paths = append(t.paths, path)
	rep, err := git.PlainInit(path, false)
	assert.NotNil(t.t, rep)
	assert.Nil(t.t, err)
	return &testRepo{
		path: path,
		rep:  rep,
	}
}

func (t *ComponentTester) ComponentCount() int {
	files, _ := ioutil.ReadDir(filepath.Join(t.workdir, "components"))
	res := 0
	for _, f := range files {
		if f.IsDir() {
			res++
		}
	}
	return res
}

func (t *ComponentTester) RootContainsComponent(path string) bool {
	if info, err := os.Stat(path); os.IsNotExist(err) || !info.IsDir() {
		return false
	}
	return true
}

func (r *testRepo) WriteCommit(name, content string) {
	co := &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test Create Branchh",
			Email: "dummy@stupid.com",
			When:  time.Now(),
		},
	}

	err := ioutil.WriteFile(filepath.Join(r.path, name), []byte(content), 0644)
	if err != nil {
		panic(err)
	}
	wt, err := r.rep.Worktree()
	if err != nil {
		panic(err)
	}
	_, err = wt.Add(".")
	if err != nil {
		panic(err)
	}
	r.lastHash, err = wt.Commit("Written", co)
	if err != nil {
		panic(err)
	}
}

func (r *testRepo) WriteFolderCommit(folder, name, content string) {
	co := &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test Create Branch",
			Email: "dummy@stupid.com",
			When:  time.Now(),
		},
	}

	path := filepath.Join(r.path, folder)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0777)
	}

	namePath := filepath.Join(path, name)
	err := ioutil.WriteFile(namePath, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
	wt, err := r.rep.Worktree()
	if err != nil {
		panic(err)
	}
	_, err = wt.Add(".")
	if err != nil {
		panic(err)
	}
	r.lastHash, err = wt.Commit("Written", co)
	if err != nil {
		panic(err)
	}
}

func (r *testRepo) CreateBranch(branch string) {

	branch = fmt.Sprintf("refs/heads/%s", branch)
	bo := &git.CheckoutOptions{
		Create: true,
		Force:  true,
		Branch: plumbing.ReferenceName(branch),
	}
	wt, err := r.rep.Worktree()
	if err != nil {
		panic(err)
	}

	err = wt.Checkout(bo)
	if err != nil {
		panic(err)
	}
}

func (r *testRepo) Checkout(branch string) {

	branch = fmt.Sprintf("refs/heads/%s", branch)
	bo := &git.CheckoutOptions{
		Create: false,
		Force:  false,
		Branch: plumbing.ReferenceName(branch),
	}
	wt, err := r.rep.Worktree()
	if err != nil {
		panic(err)
	}

	err = wt.Checkout(bo)
	if err != nil {
		panic(err)
	}
}

func (r *testRepo) Tag(tag string) {

	to := &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  "Test Create Tag",
			Email: "dummy@stupid.com",
			When:  time.Now(),
		},
		Message: "Tag created",
	}
	_, err := r.rep.CreateTag(tag, r.lastHash, to)
	if err != nil {
		panic(err)
	}
}

func (t *ComponentTester) CheckFile(u UsableComponent, file, wanted string) {
	b, err := ioutil.ReadFile(filepath.Join(u.RootPath(), file))
	assert.Nil(t.t, err)
	assert.Equal(t.t, wanted, string(b))
}

func (t *ComponentTester) CheckParameter(key, value string) {
	v, ok := t.cM.TemplateContext().Vars[key]
	if assert.True(t.t, ok) {
		assert.Equal(t.t, value, v)
	}
}

func (t *ComponentTester) CheckSpecificParameter(p model.Parameters, key, value string) {
	v, ok := p[key]
	if assert.True(t.t, ok) {
		assert.Equal(t.t, value, v)
	}
}

func (t *ComponentTester) CheckSpecificEnvVar(env model.EnvVars, key, value string) {
	v, ok := env[key]
	if assert.True(t.t, ok) {
		assert.Equal(t.t, value, v)
	}
}

func (t *ComponentTester) CheckStack(holder, stackName, compose string) {
	stack, ok := t.cM.Environment().Stacks[stackName]
	if assert.True(t.t, ok) {
		//Check that the self contained stack has been well built
		assert.Equal(t.t, stackName, stack.Name)
		stackC, err := stack.Component()
		assert.Nil(t.t, err)
		assert.NotNil(t.t, stackC)
		assert.Equal(t.t, holder, stackC.Id)

		// Check that the stack is usable and returns the correct component
		usableStack, err := t.cM.Use(stack)
		defer usableStack.Release()
		assert.Nil(t.t, err)
		assert.NotNil(t.t, usableStack)
		assert.False(t.t, usableStack.Templated())
		// Check that the stacks contains the compose/playbook file
		t.CheckFile(usableStack, "docker_compose.yml", compose)
	}
}

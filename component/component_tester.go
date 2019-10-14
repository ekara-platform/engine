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
	//ComponentTester is an helper user to run unit tests based on local GIT repositories
	ComponentTester struct {
		workdir     string
		paths       []string
		context     util.LaunchContext
		t           *testing.T
		environment *model.Environment
		cM          *Manager
		rM          *ReferenceManager
	}

	testRepo struct {
		path     string
		rep      *git.Repository
		lastHash plumbing.Hash
	}
)

//CreateComponentTester creates a new ComponentTester
func CreateComponentTester(t *testing.T, lC util.LaunchContext) *ComponentTester {
	tester := &ComponentTester{
		workdir:     "testdata/gitwork",
		context:     lC,
		t:           t,
		paths:       make([]string, 0, 0),
		environment: model.InitEnvironment(),
	}
	tester.Clean()
	tester.cM = CreateComponentManager(lC.Log(), lC.ExternalVars(), tester.environment.Platform(), tester.workdir)
	tester.rM = CreateReferenceManager(lC.Log())
	return tester
}

//Clean deletes all the content created locally during a test
func (t *ComponentTester) Clean() {
	os.RemoveAll("./testdata/gittest/")
	os.RemoveAll(t.workdir)
	for _, v := range t.paths {
		os.RemoveAll(v)
	}
}

//Init initializes the ComponentTester and build the environment
// bases on the launch context used during the tester's creation.
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
	err = t.rM.Init(mainComponent, t.cM)
	if err != nil {
		return err
	}

	return t.rM.Ensure(t.environment, t.cM)
	// TODO Link here the generated interface with the context
	// ONce this is done into the engine then remove this link.
	//t.context.templateContext.Model = t.Env()
}

//CreateRepDefaultDescriptor creates a new component folder corresponding
// to the given path and write into it an empty descriptor
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

//CreateRep creates a new component folder corresponding
// to the given path
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

//WriteCommit commit into the repo the given content into
// a file named as the provided name
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

//WriteCommit create the desired folder and then commit into it the given content into
// a file named as the provided name
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

//CreateBranch creates a branch into the repo
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

//Checkout Switch to the desired branch
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

//Tag creates a tag into the repo
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

//Env returns the environment built during the test
func (t *ComponentTester) Env() model.Environment {
	return *t.environment
}

//AssertComponentsContainsExactly asserts that the fetched component folder
// contains exactly the components identified by the provided ids
func (t *ComponentTester) AssertComponentsContainsExactly(ids ...string) bool {
	files, _ := ioutil.ReadDir(t.workdir + "/components/")
	assert.Equal(t.t, len(ids), len(files))
	return t.AssertComponentsContains(ids...)
}

//AssertComponentsContains asserts that the fetched component folder
// contains at least the components identified by the provided ids
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

//ComponentCount returns the number of components fetched into the
// component folder
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

//RootContainsComponent returns true the components tester contains
// the given path
func (t *ComponentTester) RootContainsComponent(path string) bool {
	if info, err := os.Stat(path); os.IsNotExist(err) || !info.IsDir() {
		return false
	}
	return true
}

//CheckFile asserts than a usable component contains a file with the desirec content
func (t *ComponentTester) CheckFile(u UsableComponent, file, desiredContent string) {
	b, err := ioutil.ReadFile(filepath.Join(u.RootPath(), file))
	assert.Nil(t.t, err)
	assert.Equal(t.t, desiredContent, string(b))
}

//CheckParameter asserts that the template context contains the given parameter
func (t *ComponentTester) CheckParameter(key, value string) {
	v, ok := t.cM.TemplateContext().Vars[key]
	if assert.True(t.t, ok) {
		assert.Equal(t.t, value, v)
	}
}

//CheckSpecificParameter asserts that a parameter contains the value
func (t *ComponentTester) CheckSpecificParameter(p model.Parameters, key, value string) {
	v, ok := p[key]
	if assert.True(t.t, ok) {
		assert.Equal(t.t, value, v)
	}
}

//CheckSpecificEnvVar asserts that an environment variable contains the value
func (t *ComponentTester) CheckSpecificEnvVar(env model.EnvVars, key, value string) {
	v, ok := env[key]
	if assert.True(t.t, ok) {
		assert.Equal(t.t, value, v)
	}
}

// CheckStack asserts than the environment contains the stack, then check that the stack's component
// is usable and finally check that the stack containts the provide compose content .
func (t *ComponentTester) CheckStack(holder, stackName, composeContent string) {
	stack, ok := t.environment.Stacks[stackName]
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
		t.CheckFile(usableStack, "docker_compose.yml", composeContent)
	}
}

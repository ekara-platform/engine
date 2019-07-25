package engine

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type (
	tester struct {
		workdir string
		paths   []string
		context *MockLaunchContext
		t       *testing.T
	}

	testRepo struct {
		path     string
		rep      *git.Repository
		lastHash plumbing.Hash
	}
)

func gitTester(t *testing.T, ctx *MockLaunchContext, withLog bool) *tester {
	tester := &tester{
		workdir: "testdata/gitwork",
		context: ctx,
		t:       t,
		paths:   make([]string, 0, 0),
	}
	tester.context.engine = createGitEngine(ctx.TemplateContext(), withLog)
	tester.clean()
	return tester
}

func createGitEngine(data *model.TemplateContext, withLog bool) Engine {
	var ekara Engine
	var e error
	if withLog {
		ekara, e = Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "testdata/gitwork", data)
	} else {
		ekara, e = Create(log.New(ioutil.Discard, "TEST: ", log.Ldate|log.Ltime), "testdata/gitwork", data)
	}

	if e != nil {
		panic(e)
	}
	return ekara
}

func (t *tester) clean() {
	os.RemoveAll("./testdata/gittest/")
	os.RemoveAll(t.workdir)
	for _, v := range t.paths {
		os.RemoveAll(v)
	}
}

func (t *tester) initEngine() error {
	err := t.context.engine.Init(t.context)
	// TODO Link here the generated interface with the context
	// ONce this is done into the engine then remove this link.
	//t.context.templateContext.Model = t.Env()
	return err
}

func (t *tester) env() model.Environment {
	return *t.context.engine.ComponentManager().Environment()
}

func (t *tester) assertComponentsContainsExactly(paths ...string) bool {
	files, _ := ioutil.ReadDir(t.workdir + "/components/")
	assert.Equal(t.t, len(paths), len(files))
	return t.assertComponentsContains(paths...)
}

func (t *tester) assertComponentsContains(paths ...string) bool {
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

func (t *tester) createRepDefaultDescriptor(tt *testing.T, path string) *testRepo {
	t.paths = append(t.paths, path)
	rep, err := git.PlainInit(path, false)
	assert.NotNil(t.t, rep)
	assert.Nil(t.t, err)
	res := &testRepo{
		path: path,
		rep:  rep,
	}
	res.writeCommit(tt, "ekara.yaml", ``)
	return res
}

func (t *tester) createRep(path string) *testRepo {
	t.paths = append(t.paths, path)
	rep, err := git.PlainInit(path, false)
	assert.NotNil(t.t, rep)
	assert.Nil(t.t, err)
	return &testRepo{
		path: path,
		rep:  rep,
	}
}

func (t *tester) countComponent() int {
	files, _ := ioutil.ReadDir(filepath.Join(t.workdir, "components"))
	res := 0
	for _, f := range files {
		if f.IsDir() {
			res++
		}
	}
	return res
}

func (t *tester) rootContainsComponent(path string) bool {
	if info, err := os.Stat(path); os.IsNotExist(err) || !info.IsDir() {
		return false
	}
	return true
}

func (r *testRepo) writeCommit(t *testing.T, name, content string) {
	co := &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test Create Branchh",
			Email: "dummy@stupid.com",
			When:  time.Now(),
		},
	}

	err := ioutil.WriteFile(filepath.Join(r.path, name), []byte(content), 0644)
	assert.Nil(t, err)
	wt, err := r.rep.Worktree()
	wt.Add(".")
	assert.Nil(t, err)
	r.lastHash, err = wt.Commit("Written", co)
	assert.Nil(t, err)
}

func (r *testRepo) writeFolderCommit(t *testing.T, folder, name, content string) {
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
	assert.Nil(t, err)
	wt, err := r.rep.Worktree()
	assert.Nil(t, err)
	wt.Add(".")
	r.lastHash, err = wt.Commit("Written", co)
	assert.Nil(t, err)
}

func (r *testRepo) createBranch(t *testing.T, branch string) {

	branch = fmt.Sprintf("refs/heads/%s", branch)
	bo := &git.CheckoutOptions{
		Create: true,
		Force:  true,
		Branch: plumbing.ReferenceName(branch),
	}
	wt, err := r.rep.Worktree()
	assert.Nil(t, err)

	err = wt.Checkout(bo)
	assert.Nil(t, err)
}

func (r *testRepo) checkout(t *testing.T, branch string) {

	branch = fmt.Sprintf("refs/heads/%s", branch)
	bo := &git.CheckoutOptions{
		Create: false,
		Force:  false,
		Branch: plumbing.ReferenceName(branch),
	}
	wt, err := r.rep.Worktree()
	assert.Nil(t, err)

	err = wt.Checkout(bo)
	assert.Nil(t, err)
}

func (r *testRepo) tag(t *testing.T, tag string) {

	to := &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  "Test Create Tag",
			Email: "dummy@stupid.com",
			When:  time.Now(),
		},
		Message: "Tag created",
	}
	_, err := r.rep.CreateTag(tag, r.lastHash, to)
	assert.Nil(t, err)
}

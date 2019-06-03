package engine

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type (
	tester struct {
		workdir string
		paths   []string
		context *MockLaunchContext
		engine  Engine
		t       *testing.T
	}

	testRepo struct {
		path string
		rep  *git.Repository
	}
)

func gitTester(t *testing.T, ctx MockLaunchContext) *tester {
	return &tester{
		workdir: "testdata/gitwork",
		context: &ctx,
		engine:  createGitEngine(ctx.TemplateContext()),
		t:       t,
		paths:   make([]string, 0, 0),
	}
}

func (t *tester) clean() {
	os.RemoveAll(t.workdir)
	for _, v := range t.paths {
		os.RemoveAll(v)
	}
}

func (t *tester) initEngine() error {
	engine := t.engine.Init(t.context)
	// TODO Link here the generated interface with the context
	// ONce this is done into the engine then remove this link.
	//t.context.templateContext.Model = t.Env()
	return engine
}

func (t *tester) env() model.Environment {
	return t.engine.ComponentManager().Environment()
}

func (t *tester) assertComponentsContainsExactly(paths ...string) bool {
	files, _ := ioutil.ReadDir(t.workdir + "/components/")
	assert.Equal(t.t, len(paths), len(files))
	return t.assertComponentsContains(paths...)
}

func (t *tester) assertComponentsContains(paths ...string) bool {
	for _, v := range paths {
		_, err := os.Stat(t.workdir + "/components/" + v)
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

func (t *tester) createRep(path string) testRepo {
	t.paths = append(t.paths, path)
	rep, err := git.PlainInit(path, false)
	assert.NotNil(t.t, rep)
	assert.Nil(t.t, err)
	return testRepo{
		path: path,
		rep:  rep,
	}
}

func createGitEngine(data *model.TemplateContext) Engine {
	ekara, e := Create(log.New(ioutil.Discard, "TEST: ", log.Ldate|log.Ltime), "testdata/gitwork", data)
	//ekara, e := Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "testdata/gitwork", data)
	if e != nil {
		panic(e)
	}
	return ekara
}

func (r testRepo) writeCommit(t *testing.T, name, content string) {
	co := &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test Create Branchh",
			Email: "dummy@stupid.com",
			When:  time.Now(),
		},
	}

	err := ioutil.WriteFile(r.path+"/"+name, []byte(content), 0644)
	assert.Nil(t, err)
	wt, err := r.rep.Worktree()
	assert.Nil(t, err)
	_, err = wt.Commit("Written", co)
	assert.Nil(t, err)
}

package engine

import (
	"log"
	"os"
	"testing"

	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"
)

func TestEngineLocalWithTagRef(t *testing.T) {
	engine := createTestEngine()
	c := MockLaunchContext{locationContent: "testdata/sample@v1.0.0"}
	e := engine.Init(c)
	assertOnlyWarnings(t, e)
}

func TestEngineLocalNoRef(t *testing.T) {
	engine := createTestEngine()
	c := MockLaunchContext{locationContent: "testdata/sample"}
	e := engine.Init(c)
	assertOnlyWarnings(t, e)
}

func TestEngineLocalWithRawRef(t *testing.T) {
	engine := createTestEngine()
	c := MockLaunchContext{locationContent: "testdata/sample@refs/remotes/origin/test"}
	e := engine.Init(c)
	assertOnlyWarnings(t, e)
}

func TestEngineLocalWithBranchRef(t *testing.T) {
	engine := createTestEngine()
	c := MockLaunchContext{locationContent: "testdata/sample@test"}
	e := engine.Init(c)
	assertOnlyWarnings(t, e)
}

func createTestEngine() Engine {
	os.RemoveAll("testdata/work")
	ekara, e := Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "testdata/work", map[string]interface{}{})
	if e != nil {
		panic(e)
	}
	return ekara
}

func assertOnlyWarnings(t *testing.T, e error) {
	if e != nil {
		switch e.(type) {
		case model.ValidationErrors:
			assert.False(t, e.(model.ValidationErrors).HasErrors())
		default:
			assert.Nil(t, e)
		}
	}
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

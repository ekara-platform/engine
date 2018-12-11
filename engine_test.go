package engine

import (
	"log"
	"os"
	"testing"

	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"
)

/*
func TestEngineRemoteNoTag(t *testing.T) {
	os.RemoveAll("testdata/work")
	engine1, e := Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "testdata/work", "https://github.com/ekara-platform/demomode/l", "")
	assertOnlyWarnings(t, e)
	assert.NotNil(t, engine1)

	engine2, e := Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "testdata/work", "https://github.com/ekara-platform/demomodel/", "")
	assertOnlyWarnings(t, e)
	assert.NotNil(t, engine2)
}
*/
func TestEngineLocalWithTagRef(t *testing.T) {
	engine := createTestEngine()
	e := engine.Init("testdata/sample", "1.0.0", "")
	assertOnlyWarnings(t, e)
}

func TestEngineLocalNoRef(t *testing.T) {
	engine := createTestEngine()
	e := engine.Init("testdata/sample", "", "")
	assertOnlyWarnings(t, e)
}

func TestEngineLocalWithRawRef(t *testing.T) {
	engine := createTestEngine()
	e := engine.Init("testdata/sample", "refs/remotes/origin/test", "")
	assertOnlyWarnings(t, e)
}

func TestEngineLocalWithBranchRef(t *testing.T) {
	engine := createTestEngine()
	e := engine.Init("testdata/sample", "test", "")
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

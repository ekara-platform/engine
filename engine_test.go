package engine

import (
	"log"
	"os"
	"testing"

	"github.com/lagoon-platform/model"
	"github.com/stretchr/testify/assert"
)

/*
func TestEngineRemoteNoTag(t *testing.T) {
	os.RemoveAll("testdata/work")
	engine1, e := Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "testdata/work", "https://github.com/lagoon-platform/demomode/l", "")
	assertOnlyWarnings(t, e)
	assert.NotNil(t, engine1)

	engine2, e := Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "testdata/work", "https://github.com/lagoon-platform/demomodel/", "")
	assertOnlyWarnings(t, e)
	assert.NotNil(t, engine2)
}
*/
func TestEngineLocalWithTagRef(t *testing.T) {
	engine := createTestEngine()
	e := engine.Init("testdata/sample", "v1.0.0", map[string]interface{}{})
	assertOnlyWarnings(t, e)
}

func TestEngineLocalNoRef(t *testing.T) {
	engine := createTestEngine()
	e := engine.Init("testdata/sample", "", map[string]interface{}{})
	assertOnlyWarnings(t, e)
}

func TestEngineLocalWithRawRef(t *testing.T) {
	engine := createTestEngine()
	e := engine.Init("testdata/sample", "refs/remotes/origin/test", map[string]interface{}{})
	assertOnlyWarnings(t, e)
}

func TestEngineLocalWithBranchRef(t *testing.T) {
	engine := createTestEngine()
	e := engine.Init("testdata/sample", "test", map[string]interface{}{})
	assertOnlyWarnings(t, e)
}

func createTestEngine() Lagoon {
	os.RemoveAll("testdata/work")
	lagoon, e := Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "testdata/work")
	if e != nil {
		panic(e)
	}
	return lagoon
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

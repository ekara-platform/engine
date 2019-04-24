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
	e := engine.Init("testdata/sample", "v1.0.0", "")
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

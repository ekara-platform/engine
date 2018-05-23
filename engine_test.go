package engine

import (
	"github.com/lagoon-platform/model"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestEngine(t *testing.T) {
	os.RemoveAll("testdata/work")
	engine, e := Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "testdata/work", "testdata/sample", "v1.0.0")
	assertOnlyWarnings(t, e)
	engine.ComponentManager().Ensure()
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

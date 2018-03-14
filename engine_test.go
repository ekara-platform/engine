package engine_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"github.com/lagoon-platform/engine"
)

func TestCreateEngine(t *testing.T) {
	lagoon, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "testdata/lagoon.yaml")
	assert.Nil(t, e) // no error occurred

	env := lagoon.GetEnvironment()
	assert.Equal(t, "testEnvironment", env.GetName())                               // importing file have has precedence
	assert.Equal(t, "This is my awesome Lagoon environment.", env.GetDescription()) // imported files are merged
	assert.Equal(t, []string{"tag1", "tag2"}, env.GetLabels().AsStrings())
	// FIXME assert.Contains(t, "task1", "task2", "task3", env.Hooks.Provision.After)        // order matters
}

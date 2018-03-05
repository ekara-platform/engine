package descriptor_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/lagoon-platform/descriptor"
	"log"
	"os"
)

func TestNew(t *testing.T) {
	lagoon, e := descriptor.Parse(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "testdata/lagoon.yaml")
	assert.Nil(t, e) // no error occurred

	desc := lagoon.GetDescriptor()
	assert.Nil(t, desc.Imports)                                                 // imports have been processed
	assert.Equal(t, desc.Name, "testEnvironment")                               // importing file have has precedence
	assert.Equal(t, desc.Description, "This is my awesome Lagoon environment.") // imported files are merged
	assert.Contains(t, desc.Hooks.Provision.After, "task1", "task2", "task3")   // order matters
}

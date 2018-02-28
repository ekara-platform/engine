package descriptor

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	desc, e := ParseDescriptor("testdata/lagoon.yaml")

	assert.Nil(t, e)                                                             // no error occurred
	assert.Nil(t, desc.Imports)                                                  // imports have been processed
	assert.Equal(t, desc.Name, "testEnvironment")                                // importing file have has precedence
	assert.Equal(t, desc.Description, "This is my awesome Lagoon environment.")  // imported files are merged
	assert.Contains(t, desc.Hooks.Provisioning.After, "task1", "task2", "task3") // order matters
}

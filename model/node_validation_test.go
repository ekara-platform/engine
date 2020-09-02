package model

import (
	"testing"
	_ "testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
)

// Test loading an environment without nodes.
//
// The validation must complain only about the missing nodes
//- Error: no node specified @nodes
//
func TestValidationNoNodes(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "nodes", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "no node specified", "nodes"))
}

// Test loading an node set without instances.
//
// The validation must complain only about the instances number being a positive number
//
//- Error: instances must be a positive number @nodes.managers.instances
//
func TestValidationNoNodesInstance(t *testing.T) {
	vErrs, _ := testEmptyContent(t, "nodes_instance", false)
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 1, len(vErrs.Errors))
	assert.True(t, vErrs.contains(Error, "instances must be a positive number", "nodes.managers.instances"))
}

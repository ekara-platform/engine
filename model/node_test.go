package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeDescType(t *testing.T) {
	n := NodeSet{}
	assert.Equal(t, n.DescType(), "NodeSet")
}

func TestNodeDescName(t *testing.T) {
	n := NodeSet{Name: "my_name"}
	assert.Equal(t, n.DescName(), "my_name")
}

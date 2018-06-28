package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildEquals(t *testing.T) {
	m := make(map[string]string)
	m["key1"] = "val1"
	m["key2"] = "val2"
	m["key3"] = "val3"

	r := BuildEquals(m)

	assert.Equal(t, r, "key1=val1 key2=val2 key3=val3 ")
}

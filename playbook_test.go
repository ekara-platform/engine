package engine

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildEquals(t *testing.T) {
	m := make(map[string]string)
	m["key1"] = "val1"
	m["key2"] = "val2"
	m["key3"] = "val3"

	r := BuildEquals(m)

	assert.Equal(t, true, strings.Contains(r, "key1=val1"))
	assert.Equal(t, true, strings.Contains(r, "key2=val2"))
	assert.Equal(t, true, strings.Contains(r, "key3=val3"))
}

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkMap(t *testing.T, m map[string]string, key, val string) {
	l, ok := m[key]
	assert.True(t, ok)
	assert.Equal(t, val, l)
}

func checkMapInterface(t *testing.T, m map[string]interface{}, key, val string) {
	l, ok := m[key]
	assert.True(t, ok)
	assert.Equal(t, val, l)
}

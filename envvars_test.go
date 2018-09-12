package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueEnvVars(t *testing.T) {
	ev := BuildEnvVars()
	assert.Equal(t, 0, len(ev.Content))
	ev.Add("key1", "value1")
	assert.Equal(t, 1, len(ev.Content))

	val, ok := ev.Content["key1"]
	assert.True(t, ok)
	assert.Equal(t, val, "value1")

	ev.Add("key2", "value2")
	assert.Equal(t, 2, len(ev.Content))

	val, ok = ev.Content["key2"]
	assert.True(t, ok)
	assert.Equal(t, val, "value2")
}

func TestBufferedEnvVars(t *testing.T) {
	ev := BuildEnvVars()
	assert.Equal(t, 0, len(ev.Content))

	b := CreateBuffer()
	b.Envvars["key1"] = "value1"
	b.Envvars["key2"] = "value2"

	ev.AddBuffer(b)
	assert.Equal(t, 2, len(ev.Content))

	val, ok := ev.Content["key1"]
	assert.True(t, ok)
	assert.Equal(t, val, "value1")

	val, ok = ev.Content["key2"]
	assert.True(t, ok)
	assert.Equal(t, val, "value2")
}

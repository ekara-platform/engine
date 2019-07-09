package ansible

import (
	"os"
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

func TestOsEnvVars(t *testing.T) {
	os.Setenv("HOSTNAME", "1")
	os.Setenv("PATH", "2")
	os.Setenv("TERM", "3")
	os.Setenv("HOME", "4")

	ev := BuildEnvVars()
	ev.AddDefaultOsVars()

	assert.Equal(t, 4, len(ev.Content))
	val, ok := ev.Content["HOSTNAME"]
	assert.True(t, ok)
	assert.Equal(t, val, "1")
	val, ok = ev.Content["PATH"]
	assert.True(t, ok)
	assert.Equal(t, val, "2")
	val, ok = ev.Content["TERM"]
	assert.True(t, ok)
	assert.Equal(t, val, "3")
	val, ok = ev.Content["HOME"]
	assert.True(t, ok)
	assert.Equal(t, val, "4")
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

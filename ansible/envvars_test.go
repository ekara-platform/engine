package ansible

import (
	"os"
	"testing"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestUniqueEnvVars(t *testing.T) {
	ev := BuildEnvVars()
	assert.Len(t, ev.Content, 0)
	ev.Add("key1", "value1")
	assert.Len(t, ev.Content, 1)

	val, ok := ev.Content["key1"]
	assert.True(t, ok)
	assert.Equal(t, val, "value1")

	ev.Add("key2", "value2")
	assert.Len(t, ev.Content, 2)

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

	assert.Len(t, ev.Content,4)
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

func TestProxyEnvVars(t *testing.T) {
	ev := BuildEnvVars()

	ev.AddProxy(model.Proxy{
		Http:    "testHttpProxy",
		Https:   "testHttpsProxy",
		NoProxy: "testNoProxy"})

	val, ok := ev.Content["http_proxy"]
	assert.True(t, ok)
	assert.Equal(t, val, "testHttpProxy")
	val, ok = ev.Content["https_proxy"]
	assert.True(t, ok)
	assert.Equal(t, val, "testHttpsProxy")
	val, ok = ev.Content["no_proxy"]
	assert.True(t, ok)
	assert.Equal(t, val, "testNoProxy")

	ev.AddProxy(model.Proxy{})

	val, ok = ev.Content["http_proxy"]
	assert.False(t, ok)
	val, ok = ev.Content["https_proxy"]
	assert.False(t, ok)
	val, ok = ev.Content["no_proxy"]
	assert.False(t, ok)
}

func TestBufferedEnvVars(t *testing.T) {
	ev := BuildEnvVars()
	assert.Len(t, ev.Content, 0)

	b := CreateBuffer()
	b.Envvars["key1"] = "value1"
	b.Envvars["key2"] = "value2"

	ev.AddBuffer(b)
	assert.Len(t, ev.Content, 2)

	val, ok := ev.Content["key1"]
	assert.True(t, ok)
	assert.Equal(t, val, "value1")

	val, ok = ev.Content["key2"]
	assert.True(t, ok)
	assert.Equal(t, val, "value2")
}

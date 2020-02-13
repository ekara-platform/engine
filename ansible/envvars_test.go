package ansible

import (
	"os"
	"testing"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
)

func TestUniqueEnvVars(t *testing.T) {
	ev := createEnvVars()
	assert.Len(t, ev.Content, 0)
	ev.add("key1", "value1")
	assert.Len(t, ev.Content, 1)

	val, ok := ev.Content["key1"]
	assert.True(t, ok)
	assert.Equal(t, val, "value1")

	ev.add("key2", "value2")
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

	ev := createEnvVars()
	ev.addDefaultOsVars()

	assert.Len(t, ev.Content, 4)
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

func TestAppendToVar(t *testing.T) {
	ev := createEnvVars()
	ev.appendToVar("TOTO", "/some/path")
	val, ok := ev.Content["TOTO"]
	assert.True(t, ok)
	assert.Equal(t, val, "/some/path")
	ev.appendToVar("TOTO", "/other/path")
	val, ok = ev.Content["TOTO"]
	assert.True(t, ok)
	assert.Equal(t, "/some/path"+string(os.PathListSeparator)+"/other/path", val)
}

func TestPrependToVar(t *testing.T) {
	ev := createEnvVars()
	ev.prependToVar("TOTO", "/some/path")
	val, ok := ev.Content["TOTO"]
	assert.True(t, ok)
	assert.Equal(t, val, "/some/path")
	ev.prependToVar("TOTO", "/other/path")
	val, ok = ev.Content["TOTO"]
	assert.True(t, ok)
	assert.Equal(t, "/other/path"+string(os.PathListSeparator)+"/some/path", val)
}

func TestProxyEnvVars(t *testing.T) {
	ev := createEnvVars()

	ev.addProxy(model.Proxy{
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

	ev.addProxy(model.Proxy{})

	_, ok = ev.Content["http_proxy"]
	assert.False(t, ok)
	_, ok = ev.Content["https_proxy"]
	assert.False(t, ok)
	_, ok = ev.Content["no_proxy"]
	assert.False(t, ok)
}

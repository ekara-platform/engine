package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEnvVars(t *testing.T) {
	e := CreateEnvVars(map[string]string{
		"key1": "val1",
		"key2": "val2",
	})
	if assert.Len(t, e, 2) {
		val, ok := e["key1"]
		assert.True(t, ok)
		assert.Equal(t, "val1", val)

		val, ok = e["key2"]
		assert.True(t, ok)
		assert.Equal(t, "val2", val)
	}
}

func TestEnvVarsInherits(t *testing.T) {
	ls := make(EnvVars)
	ls["key1"] = "val1"
	ls["key2"] = "val2"

	os := make(EnvVars)
	os["key2"] = "val2_overwritten"
	os["key3"] = "val3"
	os["key4"] = "val4"

	ins := ls.Override(os)
	assert.NotNil(t, ins)
	assert.Len(t, ins, 4)

	val, ok := ins["key1"]
	assert.True(t, ok)
	assert.Equal(t, "val1", val)

	val, ok = ins["key2"]
	assert.True(t, ok)
	assert.Equal(t, "val2_overwritten", val)

	val, ok = ins["key3"]
	assert.True(t, ok)
	assert.Equal(t, "val3", val)

	val, ok = ins["key4"]
	assert.True(t, ok)
	assert.Equal(t, "val4", val)

}

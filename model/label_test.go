package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabelsInherits(t *testing.T) {

	ls := make(Labels)
	ls["key1"] = "val1"
	ls["key2"] = "val2"

	os := make(Labels)
	os["key2"] = "val2_overwritten"
	os["key3"] = "val3"
	os["key4"] = "val4"

	ins := os.override(ls)
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

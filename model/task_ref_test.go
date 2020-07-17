package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTaskRefOrigin(ref string) TaskRef {
	tr1 := TaskRef{
		ref: ref,
	}
	tr1.envVars = make(map[string]string)
	tr1.envVars["key1"] = "val1_target"
	tr1.envVars["key2"] = "val2_target"
	tr1.params = make(map[string]interface{})
	tr1.params["key1"] = "val1_target"
	tr1.params["key2"] = "val2_target"

	return tr1
}

func getTaskRefOther(ref string) TaskRef {
	other := TaskRef{
		ref: ref,
	}
	other.envVars = make(map[string]string)
	other.envVars["key2"] = "val2_other"
	other.envVars["key3"] = "val3_other"
	other.params = make(map[string]interface{})
	other.params["key2"] = "val2_other"
	other.params["key3"] = "val3_other"

	return other
}

func checkTaskRefMerge(t *testing.T, ta TaskRef, expectedRef string) {
	assert.Equal(t, ta.ref, expectedRef)
	if assert.Len(t, ta.envVars, 3) {
		checkMap(t, ta.envVars, "key1", "val1_target")
		checkMap(t, ta.envVars, "key2", "val2_other")
		checkMap(t, ta.envVars, "key3", "val3_other")
	}

	if assert.Len(t, ta.params, 3) {
		checkMapInterface(t, ta.params, "key1", "val1_target")
		checkMapInterface(t, ta.params, "key2", "val2_other")
		checkMapInterface(t, ta.params, "key3", "val3_other")
	}
}

func TestTaskRefMerge(t *testing.T) {
	o := getTaskRefOrigin("my_ref")
	err := o.merge(getTaskRefOther("my_name"))
	if assert.Nil(t, err) {
		checkTaskRefMerge(t, o, "my_ref")
	}
}

func TestTaskRefMergeEmptyLocation(t *testing.T) {
	o := getTaskRefOrigin("")
	err := o.merge(getTaskRefOther("other_ref"))
	if assert.Nil(t, err) {
		checkTaskRefMerge(t, o, "other_ref")
	}
}

func TestMergeTaskRefItself(t *testing.T) {
	o := getTaskRefOrigin("my_ref")
	oi := o
	err := o.merge(o)
	if assert.Nil(t, err) {
		assert.Equal(t, oi, o)
	}
}

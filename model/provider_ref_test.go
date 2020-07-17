package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getProviderRefOrigin(ref string) ProviderRef {
	tr1 := ProviderRef{
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

func getProviderRefOther(ref string) ProviderRef {
	other := ProviderRef{
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

func checkProviderRefMerge(t *testing.T, ta ProviderRef, expectedRef string) {
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

func TestProviderRefMerge(t *testing.T) {
	o := getProviderRefOrigin("my_ref")
	o.merge(getProviderRefOther("my_name"))
	checkProviderRefMerge(t, o, "my_ref")
}

func TestProviderRefMergeEmptyLocation(t *testing.T) {
	o := getProviderRefOrigin("")
	o.merge(getProviderRefOther("other_ref"))
	checkProviderRefMerge(t, o, "other_ref")
}

func TestMergeProviderRefItself(t *testing.T) {
	o := getProviderRefOrigin("my_ref")
	oi := o
	o.merge(o)
	assert.Equal(t, oi, o)
}

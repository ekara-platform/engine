package model

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceConcatenation(t *testing.T) {
	parent := CreateParameters(map[string]interface{}{
		"slice": []int{1, 2, 3},
	})
	child := CreateParameters(map[string]interface{}{
		"slice": []int{4, 5, 6},
	})
	res := parent.Override(child)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, res["slice"])
}

func TestCloneParameters(t *testing.T) {
	m1 := map[string]interface{}{
		"a": "bbb",
		"b": map[string]interface{}{
			"c": 123,
		},
	}

	m2 := CloneParameters(m1)

	m1["a"] = "zzz"
	delete(m1, "b")

	require.Equal(t, map[string]interface{}{"a": "zzz"}, m1)
	require.Equal(t, Parameters{
		"a": "bbb",
		"b": Parameters{
			"c": 123,
		},
	}, m2)
}

func TestSliceMismatch(t *testing.T) {
	parent := CreateParameters(map[string]interface{}{
		"slice": []string{"a", "b", "c"},
	})
	child := CreateParameters(map[string]interface{}{
		"slice": []int{4, 5, 6},
	})
	res := parent.Override(child)
	assert.Equal(t, []int{4, 5, 6}, res["slice"])
}

func TestMapMerging(t *testing.T) {
	parent := CreateParameters(map[string]interface{}{
		"key1": map[interface{}]interface{}{
			"key11": "someValue",
			"key12": "otherValue",
		},
	})
	child := CreateParameters(map[string]interface{}{
		"key1": map[interface{}]interface{}{
			"key13": "thirdValue",
		},
		"key2": "unrelatedValue",
	})
	res := child.Override(parent)
	assert.Equal(t, "someValue", (res["key1"]).(map[interface{}]interface{})["key11"])
	assert.Equal(t, "otherValue", (res["key1"]).(map[interface{}]interface{})["key12"])
	assert.Equal(t, "thirdValue", (res["key1"]).(map[interface{}]interface{})["key13"])
	assert.Equal(t, "unrelatedValue", res["key2"])
}

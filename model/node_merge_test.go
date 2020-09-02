package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeMergeInstance(t *testing.T) {
	n1 := NodeSet{Name: "n1", Instances: 1}
	n2 := NodeSet{Name: "n2", Instances: 2}
	n3 := NodeSet{Name: "n3"}

	no1 := NodeSet{Name: "n1", Instances: 11}
	no2 := NodeSet{Name: "n2", Instances: 1}
	no3 := NodeSet{Name: "n3", Instances: 13}
	no4 := NodeSet{Name: "n4", Instances: 15}

	origin := make(NodeSets)
	origin[n1.Name] = n1
	origin[n2.Name] = n2
	origin[n3.Name] = n3

	other := make(NodeSets)
	other[no1.Name] = no1
	other[no2.Name] = no2
	other[no3.Name] = no3
	other[no4.Name] = no4

	origin.merge(other)
	assert.Equal(t, len(origin), 4)

	n, ok := origin["n1"]
	assert.True(t, ok)
	if assert.Equal(t, n.Name, "n1") {
		// The "other" instance number has priority when higher
		assert.Equal(t, n.Instances, 11)
	}
	n, ok = origin["n2"]
	assert.True(t, ok)
	if assert.Equal(t, n.Name, "n2") {
		assert.Equal(t, n.Instances, 1)
	}
	n, ok = origin["n3"]
	assert.True(t, ok)
	if assert.Equal(t, n.Name, "n3") {
		// If missing the instances should be added
		assert.Equal(t, n.Instances, 13)
	}
	n, ok = origin["n4"]
	assert.True(t, ok)
	if assert.Equal(t, n.Name, "n4") {
		// The new node should come with its instances
		assert.Equal(t, n.Instances, 15)
	}
}

func TestNodeMergeLabels(t *testing.T) {
	n1 := NodeSet{Name: "n1"}
	n1.Labels = make(map[string]string)
	n1.Labels["n1_lab1_k"] = "n1_lab1_untouched_v"
	n1.Labels["n1_lab2_k"] = "n1_lab2_untouched_v"

	n2 := NodeSet{Name: "n2"}
	n2.Labels = make(map[string]string)
	n2.Labels["n2_lab1_k"] = "n2_lab1_untouched_v"
	n2.Labels["n2_lab2_k"] = "n2_lab2_untouched_v"

	no1 := NodeSet{Name: "n1"}
	no1.Labels = make(map[string]string)
	no1.Labels["n1_lab1_k"] = "no1_lab1_v"

	no2 := NodeSet{Name: "n2"}
	no2.Labels = make(map[string]string)
	no2.Labels["n2_lab1_k"] = "no2_lab1_v"
	no2.Labels["n2_lab3_k"] = "no2_lab3_v"

	no3 := NodeSet{Name: "n3"}
	no3.Labels = make(map[string]string)
	no3.Labels["no3_lab1_k"] = "no3_lab1_v"
	no3.Labels["no3_lab2_k"] = "no3_lab2_v"

	origin := make(NodeSets)
	origin[n1.Name] = n1
	origin[n2.Name] = n2

	other := make(NodeSets)
	other[no1.Name] = no1
	other[no2.Name] = no2
	other[no3.Name] = no3

	origin.merge(other)
	assert.Equal(t, len(origin), 3)

	n, ok := origin["n1"]
	assert.True(t, ok)
	if assert.Equal(t, n.Name, "n1") {
		assert.Equal(t, len(n.Labels), 2)
		// The defined labels should remain unchanged
		checkMap(t, n.Labels, "n1_lab1_k", "n1_lab1_untouched_v")
		checkMap(t, n.Labels, "n1_lab2_k", "n1_lab2_untouched_v")
	}
	n, ok = origin["n2"]
	assert.True(t, ok)
	if assert.Equal(t, n.Name, "n2") {
		assert.Equal(t, len(n.Labels), 3)
		checkMap(t, n.Labels, "n2_lab1_k", "n2_lab1_untouched_v")
		checkMap(t, n.Labels, "n2_lab2_k", "n2_lab2_untouched_v")
		// The new label should be added
		checkMap(t, n.Labels, "n2_lab3_k", "no2_lab3_v")

	}
	n, ok = origin["n3"]
	assert.True(t, ok)
	if assert.Equal(t, n.Name, "n3") {
		// The new node should come with its labels
		assert.Equal(t, len(n.Labels), 2)
		checkMap(t, n.Labels, "no3_lab1_k", "no3_lab1_v")
		checkMap(t, n.Labels, "no3_lab2_k", "no3_lab2_v")
	}
}

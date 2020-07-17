package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCopies(t *testing.T) {
	copies := map[string]yamlCopy{
		"cp1": {
			yamlLabel: yamlLabel{map[string]string{
				"key1": "path1_lab1",
				"key2": "path1_lab2",
			}},
			Path:    "path1",
			Sources: []string{"path1_pattern1", "path1_pattern2"},
		},
		"cp2": {
			yamlLabel: yamlLabel{map[string]string{
				"key1": "path2_lab1",
				"key2": "path2_lab2",
			}},
			Path:    "path2",
			Sources: []string{"path2_pattern1", "path2_pattern2"},
		},
	}
	c := createCopies(copies)
	assert.NotNil(t, c)
	assert.Equal(t, len(copies), len(c.Content))

	val, ok := c.Content["cp1"]
	if assert.True(t, ok) {
		assert.Equal(t, val.Path, "path1")
		assert.Contains(t, val.Sources, "path1_pattern1")
		assert.Contains(t, val.Sources, "path1_pattern2")
		lab, ok := val.Labels["key1"]
		assert.True(t, ok)
		assert.Equal(t, lab, "path1_lab1")
		lab, ok = val.Labels["key2"]
		assert.True(t, ok)
		assert.Equal(t, lab, "path1_lab2")
	}
	val, ok = c.Content["cp2"]
	if assert.True(t, ok) {
		assert.Equal(t, val.Path, "path2")
		assert.Contains(t, val.Sources, "path2_pattern1")
		assert.Contains(t, val.Sources, "path2_pattern2")
		lab, ok := val.Labels["key1"]
		assert.True(t, ok)
		assert.Equal(t, lab, "path2_lab1")
		lab, ok = val.Labels["key2"]
		assert.True(t, ok)
		assert.Equal(t, lab, "path2_lab2")
	}
}

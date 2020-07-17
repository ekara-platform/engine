package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnionStringSlice(t *testing.T) {
	a := []string{"1", "2", "2"}
	b := []string{"2", "3", "4", "4"}

	res := union(a, b)
	assert.NotNil(t, res)

	expected := []string{"1", "2", "3", "4"}
	if assert.Equal(t, len(expected), len(res)) {
		for _, v := range res {
			assert.Contains(t, expected, v)
		}
	}
}

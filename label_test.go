package engine_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/lagoon-platform/engine"
)

type fixture struct {
	values []string
}

func (f fixture) Labels() engine.Labels {
	return engine.CreateLabels(f.values...)
}

func TestCreateLabels(t *testing.T) {
	assert.Equal(t, []string{"label1", "label2"}, fixture{[]string{"label1", "label2"}}.Labels().AsStrings())
}

func TestContainsLabels(t *testing.T) {
	f := fixture{[]string{"label1", "label2"}}.Labels()
	assert.Equal(t, true, f.Contains("label1"))
	assert.Equal(t, true, f.Contains("label2"))
	assert.Equal(t, false, f.Contains("label3"))
}

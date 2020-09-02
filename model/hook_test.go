package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasNoTask(t *testing.T) {
	h := Hook{}
	assert.False(t, h.HasTasks())
}

func TestHasTask(t *testing.T) {
	h := Hook{}
	h.Before = append(h.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskBefore(t *testing.T) {
	h := Hook{}
	h.Before = append(h.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfter(t *testing.T) {
	h := Hook{}
	h.After = append(h.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestMergeHookBefore(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := Hook{}
	h.Before = append(h.Before, task1)
	o := Hook{}
	o.Before = append(o.Before, task2)

	h.merge(o)
	assert.True(t, h.HasTasks())
	if assert.Equal(t, 2, len(h.Before)) {
		assert.Equal(t, 0, len(h.After))
		// Check if the appended content is after the genuine one
		assert.Equal(t, task1.ref, h.Before[0].ref)
		assert.Equal(t, task2.ref, h.Before[1].ref)
	}
}

func TestMergeHookAfter(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := Hook{}
	h.After = append(h.After, task1)
	o := Hook{}
	o.After = append(o.After, task2)

	h.merge(o)
	assert.True(t, h.HasTasks())
	if assert.Equal(t, 2, len(h.After)) {
		assert.Equal(t, 0, len(h.Before))
		// Check if the appended content is after the genuine one
		assert.Equal(t, task1.ref, h.After[0].ref)
		assert.Equal(t, task2.ref, h.After[1].ref)
	}
}

func TestMergeHookItself(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	h := Hook{}
	h.After = append(h.After, task1)

	h.merge(h)
	assert.True(t, h.HasTasks())
	assert.Equal(t, 0, len(h.Before))
	// Just on after because de hook cannot customize its own content
	assert.Equal(t, 1, len(h.After))
	assert.Equal(t, task1.ref, h.After[0].ref)
}

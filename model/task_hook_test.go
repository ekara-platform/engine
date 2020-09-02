package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading a task with unknown hooks
//
// The validation must complain only about 2 hooks pointing on unknown tasks
//
//- Error: empty volume path @nodes.managers.volumes.path
//
func TestValidateUnknownTaskHooks(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml("./testdata/yaml/grammar/task_unknown_hook.yaml", &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 2, len(vErrs.Errors))

	assert.True(t, vErrs.contains(Error, "no such task: unknown", "tasks.task1.hooks.execute.before[0].task"))
	assert.True(t, vErrs.contains(Error, "no such task: unknown", "tasks.task1.hooks.execute.after[0].task"))
}

func TestHasNoTaskTask(t *testing.T) {
	h := TaskHooks{}
	assert.False(t, h.HasTasks())
}

func TestHasTaskBeforeTaskExecute(t *testing.T) {
	h := TaskHooks{}
	h.Execute.Before = append(h.Execute.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterTaskExecute(t *testing.T) {
	h := TaskHooks{}
	h.Execute.After = append(h.Execute.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestMergeTaskHookBefore(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := TaskHooks{}
	h.Execute.Before = append(h.Execute.Before, task1)
	o := TaskHooks{}
	o.Execute.Before = append(o.Execute.Before, task2)

	h.merge(o)
	assert.True(t, h.HasTasks())
	if assert.Equal(t, 2, len(h.Execute.Before)) {
		assert.Equal(t, 0, len(h.Execute.After))
		assert.Equal(t, task1.ref, h.Execute.Before[0].ref)
		assert.Equal(t, task2.ref, h.Execute.Before[1].ref)
	}
}

func TestMergeTaskHookAfter(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := TaskHooks{}
	h.Execute.After = append(h.Execute.After, task1)
	o := TaskHooks{}
	o.Execute.After = append(o.Execute.After, task2)

	h.merge(o)
	assert.True(t, h.HasTasks())
	if assert.Equal(t, 2, len(h.Execute.After)) {
		assert.Equal(t, 0, len(h.Execute.Before))
		assert.Equal(t, task1.ref, h.Execute.After[0].ref)
		assert.Equal(t, task2.ref, h.Execute.After[1].ref)
	}
}

func TestMergeTaskHookItself(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	h := TaskHooks{}
	h.Execute.After = append(h.Execute.After, task1)

	h.merge(h)
	assert.True(t, h.HasTasks())
	assert.Equal(t, 0, len(h.Execute.Before))
	assert.Equal(t, 1, len(h.Execute.After))
	assert.Equal(t, task1.ref, h.Execute.After[0].ref)
}

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading s stack with unknown hooks
//
// The validation must complain only about 4 hooks pointing on unknown tasks
//
//- Error: empty volume path @nodes.managers.volumes.path
//
func TestValidateUnknownStackHooks(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml("./testdata/yaml/grammar/stack_unknown_hook.yaml", &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{Id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 2, len(vErrs.Errors))

	assert.True(t, vErrs.contains(Error, "no such task: unknown", "stacks.monitoring.hooks.deploy.before[0].task"))
	assert.True(t, vErrs.contains(Error, "no such task: unknown", "stacks.monitoring.hooks.deploy.after[0].task"))

}

func TestHasNoTaskStack(t *testing.T) {
	h := StackHooks{}
	assert.False(t, h.HasTasks())
}

func TestHasTaskBeforeStackDeploy(t *testing.T) {
	h := StackHooks{}
	h.Deploy.Before = append(h.Deploy.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterStackDeploy(t *testing.T) {
	h := StackHooks{}
	h.Deploy.After = append(h.Deploy.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestMergeStackHookBefore(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := StackHooks{}
	h.Deploy.Before = append(h.Deploy.Before, task1)
	o := StackHooks{}
	o.Deploy.Before = append(o.Deploy.Before, task2)

	h.merge(o)
	assert.True(t, h.HasTasks())
	if assert.Equal(t, 2, len(h.Deploy.Before)) {
		assert.Equal(t, 0, len(h.Deploy.After))
		assert.Equal(t, task1.ref, h.Deploy.Before[0].ref)
		assert.Equal(t, task2.ref, h.Deploy.Before[1].ref)
	}
}

func TestMergeStackHookAfter(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := StackHooks{}
	h.Deploy.After = append(h.Deploy.After, task1)
	o := StackHooks{}
	o.Deploy.After = append(o.Deploy.After, task2)

	h.merge(o)
	assert.True(t, h.HasTasks())
	if assert.Equal(t, 2, len(h.Deploy.After)) {
		assert.Equal(t, 0, len(h.Deploy.Before))
		assert.Equal(t, task1.ref, h.Deploy.After[0].ref)
		assert.Equal(t, task2.ref, h.Deploy.After[1].ref)
	}
}

func TestMergeStackHookItself(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	h := StackHooks{}
	h.Deploy.After = append(h.Deploy.After, task1)

	h.merge(h)
	assert.True(t, h.HasTasks())
	assert.Equal(t, 0, len(h.Deploy.Before))
	assert.Equal(t, 1, len(h.Deploy.After))
	assert.Equal(t, task1.ref, h.Deploy.After[0].ref)
}

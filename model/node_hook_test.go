package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading an nodeset with unknown hooks
//
// The validation must complain only about 2 hooks pointing on unknown tasks
//
//- Error: empty volume path @nodes.managers.volumes.path
//
func TestValidationNodesUnknownHook(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml("./testdata/yaml/grammar/nodes_unknown_hook.yaml", &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 2, len(vErrs.Errors))

	assert.True(t, vErrs.contains(Error, "no such task: unknown", "nodes.managers.hooks.create.before[0].task"))
	assert.True(t, vErrs.contains(Error, "no such task: unknown", "nodes.managers.hooks.create.after[0].task"))
}

// Test loading an nodeset with valid hooks
func TestValidationNodesKnownHook(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml("./testdata/yaml/grammar/nodes_known_hook.yaml", &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	vErrs := env.Validate()
	assert.NotNil(t, vErrs)
	assert.False(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 0, len(vErrs.Errors))
}

func TestHasNoTaskNode(t *testing.T) {
	h := NodeHooks{}
	assert.False(t, h.HasTasks())
}

func TestHasTaskBeforeNodeCreate(t *testing.T) {
	h := NodeHooks{}
	h.Create.Before = append(h.Create.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterNodeCreate(t *testing.T) {
	h := NodeHooks{}
	h.Create.After = append(h.Create.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskBeforeNodeDestroy(t *testing.T) {
	h := NodeHooks{}
	h.Destroy.Before = append(h.Destroy.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterNodeDestroy(t *testing.T) {
	h := NodeHooks{}
	h.Destroy.After = append(h.Destroy.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestMergeNodeHookBefore(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := NodeHooks{}
	h.Create.Before = append(h.Create.Before, task1)
	h.Destroy.Before = append(h.Destroy.Before, task1)

	o := NodeHooks{}
	o.Create.Before = append(o.Create.Before, task2)
	o.Destroy.Before = append(o.Destroy.Before, task2)

	h.merge(o)
	assert.True(t, h.HasTasks())
	if assert.Equal(t, 2, len(h.Create.Before)) {
		assert.Equal(t, 0, len(h.Create.After))
		assert.Equal(t, task1.ref, h.Create.Before[0].ref)
		assert.Equal(t, task2.ref, h.Create.Before[1].ref)
	}
	if assert.Equal(t, 2, len(h.Destroy.Before)) {
		assert.Equal(t, 0, len(h.Destroy.After))
		assert.Equal(t, task1.ref, h.Destroy.Before[0].ref)
		assert.Equal(t, task2.ref, h.Destroy.Before[1].ref)
	}

}

func TestMergeNodeHookAfter(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := NodeHooks{}
	h.Create.After = append(h.Create.After, task1)
	h.Destroy.After = append(h.Destroy.After, task1)
	o := NodeHooks{}
	o.Create.After = append(o.Create.After, task2)
	o.Destroy.After = append(o.Destroy.After, task2)

	h.merge(o)
	assert.True(t, h.HasTasks())
	if assert.Equal(t, 2, len(h.Create.After)) {
		assert.Equal(t, 0, len(h.Create.Before))
		assert.Equal(t, task1.ref, h.Create.After[0].ref)
		assert.Equal(t, task2.ref, h.Create.After[1].ref)
	}
	if assert.Equal(t, 2, len(h.Destroy.After)) {
		assert.Equal(t, 0, len(h.Destroy.Before))
		assert.Equal(t, task1.ref, h.Destroy.After[0].ref)
		assert.Equal(t, task2.ref, h.Destroy.After[1].ref)
	}
}

func TestMergeNodeHookItself(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	h := NodeHooks{}
	h.Create.After = append(h.Create.After, task1)
	h.Destroy.Before = append(h.Destroy.Before, task1)

	h.merge(h)
	assert.True(t, h.HasTasks())
	assert.Equal(t, 0, len(h.Create.Before))
	assert.Equal(t, 1, len(h.Create.After))
	assert.Equal(t, 1, len(h.Destroy.Before))
	assert.Equal(t, 0, len(h.Destroy.After))
	assert.Equal(t, task1.ref, h.Create.After[0].ref)
	assert.Equal(t, task1.ref, h.Destroy.Before[0].ref)
}

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test loading an environment with unknown global hooks
//
// The validation must complain only about 10 hooks pointing on unknown tasks
//
//- Error: empty volume path @nodes.managers.volumes.path
//

var oneTask = TaskRef{}

func TestValidateUnknownGlobalHooks(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	e := parseYaml("./testdata/yaml/grammar/unknown_global_hook.yaml", &TemplateContext{}, &yamlEnv)
	assert.Nil(t, e)
	env, e := CreateEnvironment(component{id: MainComponentId}, yamlEnv)
	assert.Nil(t, e)
	assert.NotNil(t, env)
	vErrs := env.Validate()
	assert.True(t, vErrs.HasErrors())
	assert.False(t, vErrs.HasWarnings())
	assert.Equal(t, 10, len(vErrs.Errors))

	assert.True(t, vErrs.contains(Error, "no such task: unknown", "hooks.init.before[0].task"))
	assert.True(t, vErrs.contains(Error, "no such task: unknown", "hooks.init.after[0].task"))
	assert.True(t, vErrs.contains(Error, "no such task: unknown", "hooks.create.before[0].task"))
	assert.True(t, vErrs.contains(Error, "no such task: unknown", "hooks.create.after[0].task"))
	assert.True(t, vErrs.contains(Error, "no such task: unknown", "hooks.install.before[0].task"))
	assert.True(t, vErrs.contains(Error, "no such task: unknown", "hooks.install.after[0].task"))
	assert.True(t, vErrs.contains(Error, "no such task: unknown", "hooks.deploy.before[0].task"))
	assert.True(t, vErrs.contains(Error, "no such task: unknown", "hooks.deploy.after[0].task"))
	assert.True(t, vErrs.contains(Error, "no such task: unknown", "hooks.delete.before[0].task"))
	assert.True(t, vErrs.contains(Error, "no such task: unknown", "hooks.delete.after[0].task"))
}

func TestHasNoTaskEnv(t *testing.T) {
	h := EnvironmentHooks{}
	assert.False(t, h.HasTasks())
}

func TestHasTaskBeforeEnvInit(t *testing.T) {
	h := EnvironmentHooks{}
	h.Init.Before = append(h.Init.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterEnvInit(t *testing.T) {
	h := EnvironmentHooks{}
	h.Init.After = append(h.Init.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskBeforeEnvCreate(t *testing.T) {
	h := EnvironmentHooks{}
	h.Create.Before = append(h.Create.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterEnvCreate(t *testing.T) {
	h := EnvironmentHooks{}
	h.Create.After = append(h.Create.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskBeforeEnvInstall(t *testing.T) {
	h := EnvironmentHooks{}
	h.Install.Before = append(h.Install.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterEnvInstall(t *testing.T) {
	h := EnvironmentHooks{}
	h.Install.After = append(h.Install.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskBeforeEnvDeploy(t *testing.T) {
	h := EnvironmentHooks{}
	h.Deploy.Before = append(h.Deploy.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterEnvDeploy(t *testing.T) {
	h := EnvironmentHooks{}
	h.Deploy.After = append(h.Deploy.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskBeforeEnvDelete(t *testing.T) {
	h := EnvironmentHooks{}
	h.Destroy.Before = append(h.Destroy.Before, oneTask)
	assert.True(t, h.HasTasks())
}

func TestHasTaskAfterEnvDelete(t *testing.T) {
	h := EnvironmentHooks{}
	h.Destroy.After = append(h.Destroy.After, oneTask)
	assert.True(t, h.HasTasks())
}

func TestMergeEnvironmentHookBefore(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := EnvironmentHooks{}
	h.Init.Before = append(h.Init.Before, task1)
	h.Create.Before = append(h.Create.Before, task1)
	h.Install.Before = append(h.Install.Before, task1)
	h.Deploy.Before = append(h.Deploy.Before, task1)
	h.Destroy.Before = append(h.Destroy.Before, task1)
	o := EnvironmentHooks{}
	o.Init.Before = append(o.Init.Before, task2)
	o.Create.Before = append(o.Create.Before, task2)
	o.Install.Before = append(o.Install.Before, task2)
	o.Deploy.Before = append(o.Deploy.Before, task2)
	o.Destroy.Before = append(o.Destroy.Before, task2)

	h.merge(o)
	assert.True(t, h.HasTasks())

	if assert.Equal(t, 2, len(h.Init.Before)) {
		assert.Equal(t, 0, len(h.Init.After))
		assert.Equal(t, task1.ref, h.Init.Before[0].ref)
		assert.Equal(t, task2.ref, h.Init.Before[1].ref)
	}

	if assert.Equal(t, 2, len(h.Create.Before)) {
		assert.Equal(t, 0, len(h.Create.After))
		assert.Equal(t, task1.ref, h.Create.Before[0].ref)
		assert.Equal(t, task2.ref, h.Create.Before[1].ref)
	}

	if assert.Equal(t, 2, len(h.Install.Before)) {
		assert.Equal(t, 0, len(h.Install.After))
		assert.Equal(t, task1.ref, h.Install.Before[0].ref)
		assert.Equal(t, task2.ref, h.Install.Before[1].ref)
	}

	if assert.Equal(t, 2, len(h.Deploy.Before)) {
		assert.Equal(t, 0, len(h.Deploy.After))
		assert.Equal(t, task1.ref, h.Deploy.Before[0].ref)
		assert.Equal(t, task2.ref, h.Deploy.Before[1].ref)
	}

	if assert.Equal(t, 2, len(h.Destroy.Before)) {
		assert.Equal(t, 0, len(h.Destroy.After))
		assert.Equal(t, task1.ref, h.Destroy.Before[0].ref)
		assert.Equal(t, task2.ref, h.Destroy.Before[1].ref)
	}
}

func TestMergeEnvironmentHookAfter(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := EnvironmentHooks{}
	h.Init.After = append(h.Init.After, task1)
	h.Create.After = append(h.Create.After, task1)
	h.Install.After = append(h.Install.After, task1)
	h.Deploy.After = append(h.Deploy.After, task1)
	h.Destroy.After = append(h.Destroy.After, task1)
	o := EnvironmentHooks{}
	o.Init.After = append(o.Init.After, task2)
	o.Create.After = append(o.Create.After, task2)
	o.Install.After = append(o.Install.After, task2)
	o.Deploy.After = append(o.Deploy.After, task2)
	o.Destroy.After = append(o.Destroy.After, task2)

	h.merge(o)
	assert.True(t, h.HasTasks())

	if assert.Equal(t, 2, len(h.Init.After)) {
		assert.Equal(t, 0, len(h.Init.Before))
		assert.Equal(t, task1.ref, h.Init.After[0].ref)
		assert.Equal(t, task2.ref, h.Init.After[1].ref)
	}

	if assert.Equal(t, 2, len(h.Create.After)) {
		assert.Equal(t, 0, len(h.Create.Before))
		assert.Equal(t, task1.ref, h.Create.After[0].ref)
		assert.Equal(t, task2.ref, h.Create.After[1].ref)
	}

	if assert.Equal(t, 2, len(h.Install.After)) {
		assert.Equal(t, 0, len(h.Install.Before))
		assert.Equal(t, task1.ref, h.Install.After[0].ref)
		assert.Equal(t, task2.ref, h.Install.After[1].ref)
	}

	if assert.Equal(t, 2, len(h.Deploy.After)) {
		assert.Equal(t, 0, len(h.Deploy.Before))
		assert.Equal(t, task1.ref, h.Deploy.After[0].ref)
		assert.Equal(t, task2.ref, h.Deploy.After[1].ref)
	}

	if assert.Equal(t, 2, len(h.Destroy.After)) {
		assert.Equal(t, 0, len(h.Destroy.Before))
		assert.Equal(t, task1.ref, h.Destroy.After[0].ref)
		assert.Equal(t, task2.ref, h.Destroy.After[1].ref)
	}
}

func TestMergeEnvironmentHookItself(t *testing.T) {
	task1 := TaskRef{ref: "ref1"}
	h := EnvironmentHooks{}
	h.Init.After = append(h.Init.After, task1)
	h.Create.After = append(h.Create.After, task1)
	h.Install.After = append(h.Install.After, task1)
	h.Deploy.After = append(h.Deploy.After, task1)
	h.Destroy.After = append(h.Destroy.After, task1)

	h.merge(h)
	assert.True(t, h.HasTasks())
	assert.Equal(t, 0, len(h.Init.Before))
	assert.Equal(t, 0, len(h.Create.Before))
	assert.Equal(t, 0, len(h.Install.Before))
	assert.Equal(t, 0, len(h.Deploy.Before))
	assert.Equal(t, 0, len(h.Destroy.Before))
	assert.Equal(t, 1, len(h.Init.After))
	assert.Equal(t, 1, len(h.Create.After))
	assert.Equal(t, 1, len(h.Install.After))
	assert.Equal(t, 1, len(h.Deploy.After))
	assert.Equal(t, 1, len(h.Destroy.After))
	assert.Equal(t, task1.ref, h.Init.After[0].ref)
	assert.Equal(t, task1.ref, h.Create.After[0].ref)
	assert.Equal(t, task1.ref, h.Install.After[0].ref)
	assert.Equal(t, task1.ref, h.Deploy.After[0].ref)
	assert.Equal(t, task1.ref, h.Destroy.After[0].ref)
}

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var taskOtherT1, taskOtherT2, taskOriginT1, taskOriginT2 TaskRef

func TestTaskDescType(t *testing.T) {
	s := Task{}
	assert.Equal(t, s.DescType(), "Task")
}

func TestTaskDescName(t *testing.T) {
	s := Task{Name: "my_name"}
	assert.Equal(t, s.DescName(), "my_name")
}

func getTaskOrigin() Task {
	t1 := Task{
		Name:     "my_name",
		Playbook: "Playbook",
	}
	t1.envVars = make(map[string]string)
	t1.envVars["key1"] = "val1_target"
	t1.envVars["key2"] = "val2_target"
	t1.params = make(map[string]interface{})
	t1.params["key1"] = "val1_target"
	t1.params["key2"] = "val2_target"

	t1.cRef = componentRef{
		ref: "cOriginal",
	}

	taskOriginT1 = TaskRef{ref: "T1"}
	taskOriginT2 = TaskRef{ref: "T2"}
	hT1 := TaskHooks{}
	hT1.Execute.Before = append(hT1.Execute.Before, taskOriginT1)
	hT1.Execute.After = append(hT1.Execute.After, taskOriginT2)
	t1.Hooks = hT1

	return t1
}

func getTaskOther(name string) Task {
	other := Task{
		Name:     name,
		Playbook: "Playbook_overwritten",
	}
	other.envVars = make(map[string]string)
	other.envVars["key2"] = "val2_other"
	other.envVars["key3"] = "val3_other"
	other.params = make(map[string]interface{})
	other.params["key2"] = "val2_other"
	other.params["key3"] = "val3_other"

	other.cRef = componentRef{
		ref: "cOther",
	}

	taskOtherT1 = TaskRef{ref: "T3"}
	taskOtherT2 = TaskRef{ref: "T4"}
	oT1 := TaskHooks{}
	oT1.Execute.Before = append(oT1.Execute.Before, taskOtherT1)
	oT1.Execute.After = append(oT1.Execute.After, taskOtherT2)
	other.Hooks = oT1

	return other
}

func checkTaskMerge(t *testing.T, ta Task) {
	assert.Equal(t, "cOther", ta.cRef.ref)
	assert.Equal(t, "Playbook_overwritten", ta.Playbook)

	if assert.Len(t, ta.envVars, 3) {
		checkMap(t, ta.envVars, "key1", "val1_target")
		checkMap(t, ta.envVars, "key2", "val2_other")
		checkMap(t, ta.envVars, "key3", "val3_other")
	}

	if assert.Len(t, ta.params, 3) {
		checkMapInterface(t, ta.params, "key1", "val1_target")
		checkMapInterface(t, ta.params, "key2", "val2_other")
		checkMapInterface(t, ta.params, "key3", "val3_other")
	}

	if assert.Len(t, ta.Hooks.Execute.Before, 2) {
		assert.Contains(t, ta.Hooks.Execute.Before, taskOriginT1, taskOtherT1)
		assert.Equal(t, taskOriginT1, ta.Hooks.Execute.Before[0])
		assert.Equal(t, taskOtherT1, ta.Hooks.Execute.Before[1])
	}

	if assert.Len(t, ta.Hooks.Execute.After, 2) {
		assert.Contains(t, ta.Hooks.Execute.After, taskOriginT2, taskOtherT2)
		assert.Equal(t, taskOriginT2, ta.Hooks.Execute.After[0])
		assert.Equal(t, taskOtherT2, ta.Hooks.Execute.After[1])
	}

}

func TestTaskMerge(t *testing.T) {
	o := getTaskOrigin()
	o.merge(getTaskOther("my_name"))
	checkTaskMerge(t, o)
}

func TestMergeTaskItself(t *testing.T) {
	o := getTaskOrigin()
	oi := o
	o.merge(o)
	assert.Equal(t, oi, o)
}

func TestTasksMerge(t *testing.T) {
	origins := make(Tasks)
	origins["myS"] = getTaskOrigin()
	others := make(Tasks)
	others["myS"] = getTaskOther("my_name")

	origins.merge(others)
	if assert.Len(t, origins, 1) {
		o := origins["myS"]
		checkTaskMerge(t, o)
	}
}

func TestTasksMergeAddition(t *testing.T) {
	origins := make(Tasks)
	origins["myS"] = getTaskOrigin()
	others := make(Tasks)
	others["myS"] = getTaskOther("my_name")
	others["new"] = getTaskOther("new")

	origins.merge(others)
	assert.Len(t, origins, 2)
}

func TestTasksEmptyMerge(t *testing.T) {
	origins := make(Tasks)
	o := getTaskOrigin()
	origins["myS"] = o
	others := make(Tasks)

	origins.merge(others)
	if assert.Len(t, origins, 1) {
		oc := origins["myS"]
		assert.Equal(t, o, oc)
	}
}

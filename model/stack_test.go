package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var stackOtherT1, stackOtherT2, stackOriginT1, stackOriginT2 TaskRef

func TestStackDescType(t *testing.T) {
	s := Stack{}
	assert.Equal(t, s.DescType(), "Stack")
}

func TestStackDescName(t *testing.T) {
	s := Stack{Name: "my_name"}
	assert.Equal(t, s.DescName(), "my_name")
}

func getStackOrigin() Stack {
	s1 := Stack{Name: "my_name", envVars: make(map[string]string)}
	s1.EnvVars()["key1"] = "val1_target"
	s1.EnvVars()["key2"] = "val2_target"
	s1.params = make(map[string]interface{})
	s1.params["key1"] = "val1_target"
	s1.params["key2"] = "val2_target"

	s1.cRef = componentRef{
		ref: "cOriginal",
	}
	s1.Dependencies = []string{"dep1", "dep2"}

	stackOriginT1 = TaskRef{ref: "T1"}
	stackOriginT2 = TaskRef{ref: "T2"}
	s1.Hooks.Deploy.Before = append(s1.Hooks.Deploy.Before, stackOriginT1)
	s1.Hooks.Deploy.After = append(s1.Hooks.Deploy.After, stackOriginT2)

	s1.Copies = Copies{
		Content: make(map[string]Copy),
	}
	s1.Copies.Content["path1"] = Copy{
		Labels: map[string]string{
			"label_key1": "label_value1",
			"label_key2": "label_value2",
		},
		Sources: []string{"path1", "path2"},
	}
	return s1
}

func getStackOther(name string) Stack {
	other := Stack{Name: name, envVars: make(map[string]string)}
	other.EnvVars()["key2"] = "val2_other"
	other.EnvVars()["key3"] = "val3_other"
	other.params = make(map[string]interface{})
	other.params["key2"] = "val2_other"
	other.params["key3"] = "val3_other"
	other.cRef = componentRef{
		ref: "cOther",
	}
	other.Dependencies = []string{"dep2", "dep3"}

	stackOtherT1 = TaskRef{ref: "T3"}
	stackOtherT1 = TaskRef{ref: "T4"}
	other.Hooks.Deploy.Before = append(other.Hooks.Deploy.Before, stackOtherT1)
	other.Hooks.Deploy.After = append(other.Hooks.Deploy.After, stackOtherT2)

	other.Copies = Copies{
		Content: make(map[string]Copy),
	}
	other.Copies.Content["path2"] = Copy{
		Labels: map[string]string{
			"label_key3": "label_value3",
			"label_key4": "label_value4",
		},
		Sources: []string{"path3", "path4"},
	}
	return other
}

func checkStackMerge(t *testing.T, s Stack) {
	assert.Equal(t, "cOther", s.cRef.ref)
	assert.Contains(t, s.Dependencies, "dep1")
	assert.Contains(t, s.Dependencies, "dep2")
	assert.Contains(t, s.Dependencies, "dep3")

	if assert.Len(t, s.EnvVars(), 3) {
		checkMap(t, s.EnvVars(), "key1", "val1_target")
		checkMap(t, s.EnvVars(), "key2", "val2_other")
		checkMap(t, s.EnvVars(), "key3", "val3_other")
	}

	if assert.Len(t, s.params, 3) {
		checkMapInterface(t, s.params, "key1", "val1_target")
		checkMapInterface(t, s.params, "key2", "val2_other")
		checkMapInterface(t, s.params, "key3", "val3_other")
	}

	if assert.Len(t, s.Hooks.Deploy.Before, 2) {
		assert.Contains(t, s.Hooks.Deploy.Before, stackOriginT1, stackOtherT1)
		assert.Equal(t, s.Hooks.Deploy.Before[0], stackOriginT1)
		assert.Equal(t, s.Hooks.Deploy.Before[1], stackOtherT1)
	}

	if assert.Len(t, s.Hooks.Deploy.After, 2) {
		assert.Contains(t, s.Hooks.Deploy.After, stackOriginT2, stackOtherT2)
		assert.Equal(t, s.Hooks.Deploy.After[0], stackOriginT2)
		assert.Equal(t, s.Hooks.Deploy.After[1], stackOtherT2)
	}

	assert.Len(t, s.Copies.Content, 2)
}

func TestStackMerge(t *testing.T) {
	o := getStackOrigin()
	o.merge(getStackOther("my_name"))
	checkStackMerge(t, o)
}

func TestMergeStackItself(t *testing.T) {
	o := getStackOrigin()
	oi := o
	o.merge(o)
	assert.Equal(t, oi, o)
}

func TestStacksMerge(t *testing.T) {
	origins := make(Stacks)
	origins["myS"] = getStackOrigin()
	others := make(Stacks)
	others["myS"] = getStackOther("my_name")

	origins.merge(others)
	if assert.Len(t, origins, 1) {
		o := origins["myS"]
		checkStackMerge(t, o)
	}
}

func TestStacksMergeAddition(t *testing.T) {
	origins := make(Stacks)
	origins["myS"] = getStackOrigin()
	others := make(Stacks)
	others["myS"] = getStackOther("my_name")
	others["new"] = getStackOther("new")

	origins.merge(others)
	assert.Len(t, origins, 2)
}

func TestStacksEmptyMerge(t *testing.T) {
	origins := make(Stacks)
	o := getStackOrigin()
	origins["myS"] = o
	others := make(Stacks)

	origins.merge(others)
	if assert.Len(t, origins, 1) {
		oc := origins["myS"]
		assert.Equal(t, o, oc)
	}
}

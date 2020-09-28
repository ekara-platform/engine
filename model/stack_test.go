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

	s1.Copies = make(map[string]Copy)
	s1.Copies["path1"] = Copy{
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

	other.Copies = make(map[string]Copy)
	other.Copies["path2"] = Copy{
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

	assert.Len(t, s.Copies, 2)
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

func TestStacksDependencies(t *testing.T) {
	sts := Stacks{}

	// This test will check the dependencies on the following tree
	//
	//        0                    1
	//                           / | \
	//                          2  7  8
	//                        / |     | \
	//                       3  6     9  12
	//                      / |       | \
	//                     4  5       10  11

	sts["6"] = Stack{Name: "6", Dependencies: []string{"2"}}
	sts["1"] = Stack{Name: "1"}
	sts["3"] = Stack{Name: "3", Dependencies: []string{"2"}}
	sts["4"] = Stack{Name: "4", Dependencies: []string{"3"}}
	sts["9"] = Stack{Name: "9", Dependencies: []string{"8"}}
	sts["5"] = Stack{Name: "5", Dependencies: []string{"3"}}
	sts["7"] = Stack{Name: "7", Dependencies: []string{"1"}}
	sts["8"] = Stack{Name: "8", Dependencies: []string{"1"}}
	sts["10"] = Stack{Name: "10", Dependencies: []string{"9"}}
	sts["2"] = Stack{Name: "2", Dependencies: []string{"1"}}
	sts["12"] = Stack{Name: "12", Dependencies: []string{"8"}}
	sts["11"] = Stack{Name: "11", Dependencies: []string{"9"}}
	sts["0"] = Stack{Name: "0"}

	assert.Len(t, sts, 13)
	resolved, err := sts.sorted()
	assert.Equal(t, len(resolved), len(sts))
	assert.Nil(t, err)

	processed := make(map[string]Stack)
	for _, val := range resolved {
		if len(processed) == 0 {
			processed[val.Name] = val
			continue
		}
		//Check than all the dependencies has been already processd
		for _, d := range val.Dependencies {
			if _, ok := processed[d]; !ok {
				assert.Fail(t, "Dependency %s has not been yet processed", d)
			}

		}
		processed[val.Name] = val
	}
	//Check that the original Stacks has been untouched
	assert.Len(t, sts, 13)

}

func TestStacksMultiplesDependencies(t *testing.T) {
	sts := Stacks{}

	// This test will check the dependencies on the following tree
	//
	//        0                    1
	//                           / | \
	//                          2  7  8
	//                        / | /    | \
	//                       3  6     9  12
	//                      / |   \   | \
	//                     4  5    \-10  11

	sts["6"] = Stack{Name: "6", Dependencies: []string{"2", "7"}}
	sts["1"] = Stack{Name: "1"}
	sts["3"] = Stack{Name: "3", Dependencies: []string{"2"}}
	sts["4"] = Stack{Name: "4", Dependencies: []string{"3"}}
	sts["9"] = Stack{Name: "9", Dependencies: []string{"8"}}
	sts["5"] = Stack{Name: "5", Dependencies: []string{"3"}}
	sts["7"] = Stack{Name: "7", Dependencies: []string{"1"}}
	sts["8"] = Stack{Name: "8", Dependencies: []string{"1"}}
	sts["10"] = Stack{Name: "10", Dependencies: []string{"9", "6"}}
	sts["2"] = Stack{Name: "2", Dependencies: []string{"1"}}
	sts["12"] = Stack{Name: "12", Dependencies: []string{"8"}}
	sts["11"] = Stack{Name: "11", Dependencies: []string{"9"}}
	sts["0"] = Stack{Name: "0"}

	assert.Len(t, sts, 13)
	resolved, err := sts.sorted()
	assert.Equal(t, len(resolved), len(sts))
	assert.Nil(t, err)

	processed := make(map[string]Stack)
	for _, val := range resolved {
		if len(processed) == 0 {
			processed[val.Name] = val
			continue
		}
		//Check than all the dependencies has been already processd
		for _, d := range val.Dependencies {
			if _, ok := processed[d]; !ok {
				assert.Fail(t, "Dependency %s has not been yet processed", d)
			}

		}
		processed[val.Name] = val
	}
	//Check that the original Stacks has been untouched
	assert.Len(t, sts, 13)

}

func TestStacksCyclicDependencies(t *testing.T) {
	sts := Stacks{}

	sts["1"] = Stack{Name: "1", Dependencies: []string{"3"}}
	sts["2"] = Stack{Name: "2", Dependencies: []string{"1"}}
	sts["3"] = Stack{Name: "3", Dependencies: []string{"2"}}

	assert.Len(t, sts, 3)
	_, err := sts.sorted()
	assert.NotNil(t, err)

	//Check that the original Stacks has been untouched
	assert.Len(t, sts, 3)

}

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderDescType(t *testing.T) {
	p := Provider{}
	assert.Equal(t, p.DescType(), "Provider")
}

func TestProviderDescName(t *testing.T) {
	p := Provider{Name: "my_name"}
	assert.Equal(t, p.DescName(), "my_name")
}

func getProviderOrigin() Provider {
	p1 := Provider{Name: "my_name"}
	p1.envVars = make(map[string]string)
	p1.envVars["key1"] = "val1_target"
	p1.envVars["key2"] = "val2_target"
	p1.params = make(map[string]interface{})
	p1.params["key1"] = "val1_target"
	p1.params["key2"] = "val2_target"

	p1.cRef = componentRef{
		ref: "cOriginal",
	}

	return p1
}

func getProviderOther(name string) Provider {
	other := Provider{Name: name}
	other.envVars = make(map[string]string)
	other.envVars["key2"] = "val2_other"
	other.envVars["key3"] = "val3_other"
	other.params = make(map[string]interface{})
	other.params["key2"] = "val2_other"
	other.params["key3"] = "val3_other"

	other.cRef = componentRef{
		ref: "cOther",
	}

	return other
}

func checkProviderMerge(t *testing.T, p Provider) {
	assert.Equal(t, p.cRef.ref, "cOther")

	if assert.Equal(t, 3, len(p.envVars)) {
		checkMap(t, p.envVars, "key1", "val1_target")
		checkMap(t, p.envVars, "key2", "val2_other")
		checkMap(t, p.envVars, "key3", "val3_other")
	}

	if assert.Equal(t, 3, len(p.Parameters())) {
		checkMapInterface(t, p.Parameters(), "key1", "val1_target")
		checkMapInterface(t, p.Parameters(), "key2", "val2_other")
		checkMapInterface(t, p.Parameters(), "key3", "val3_other")
	}
}

func TestProviderMerge(t *testing.T) {
	o := getProviderOrigin()
	o.merge(getProviderOther("my_name"))
	checkProviderMerge(t, o)
}

func TestMergeProviderItself(t *testing.T) {
	o := getProviderOrigin()
	oi := o
	o.merge(o)
	assert.Equal(t, oi, o)
}

func TestProvidersMerge(t *testing.T) {
	origins := make(Providers)
	origins["myP"] = getProviderOrigin()
	others := make(Providers)
	others["myP"] = getProviderOther("my_name")

	origins.merge(others)
	if assert.Len(t, origins, 1) {
		o := origins["myP"]
		checkProviderMerge(t, o)
	}
}

func TestProvidersMergeAddition(t *testing.T) {
	origins := make(Providers)
	origins["myP"] = getProviderOrigin()
	others := make(Providers)
	others["myP"] = getProviderOther("my_name")
	others["new"] = getProviderOther("new")

	origins.merge(others)
	assert.Len(t, origins, 2)
}

func TestProvidersEmptyMerge(t *testing.T) {
	origins := make(Providers)
	o := getProviderOrigin()
	origins["myP"] = o
	others := make(Providers)

	origins.merge(others)
	if assert.Len(t, origins, 1) {
		oc := origins["myP"]
		assert.Equal(t, o, oc)
	}
}

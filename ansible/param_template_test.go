package ansible

import (
	"testing"

	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"
)

func TestParseParamSimple(t *testing.T) {
	path := "./testdata/params-template1.yaml"
	pt, err := ParseParams(path)
	assert.NotNil(t, pt)
	assert.Nil(t, err)
	checkContent(t, pt, "value1")
}

func TestParseParamLeadingSpaces(t *testing.T) {
	path := "./testdata/params-template2.yaml"
	pt, err := ParseParams(path)
	assert.NotNil(t, pt)
	assert.Nil(t, err)
	checkContent(t, pt, "  value1")
}

func TestParseParamTrailingSpaces(t *testing.T) {
	path := "./testdata/params-template3.yaml"
	pt, err := ParseParams(path)
	assert.NotNil(t, pt)
	assert.Nil(t, err)
	checkContent(t, pt, "value1  ")
}

func TestParseParamSpaces(t *testing.T) {
	path := "./testdata/params-template4.yaml"
	pt, err := ParseParams(path)
	assert.NotNil(t, pt)
	assert.Nil(t, err)
	checkContent(t, pt, "  value1  ")
}

func checkContent(t *testing.T, pt model.Parameters, wanted string) {
	assert.Len(t, pt, 1)

	val, ok := pt["key1"]
	if assert.True(t, ok) {
		_, ok = val.(map[interface{}]interface{})
		assert.True(t, ok)

		switch x := val.(type) {
		case map[interface{}]interface{}:
			assert.Len(t, x, 1)
			val2, ok := x["key2"]
			assert.True(t, ok)
			assert.Equal(t, wanted, val2)
		default:
			assert.FailNow(t, "Wrong type")
		}
	}
}

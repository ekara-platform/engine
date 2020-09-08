package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTemplateContext(t *testing.T) {
	p := CreateParameters(map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})
	assert.Equal(t, 2, len(p))

	pc := CreateTemplateContext(p)

	assert.NotNil(t, pc)
	assert.Equal(t, 2, len(pc.Vars))
	va, ok := pc.Vars["key1"]
	assert.True(t, ok)
	assert.Equal(t, va, "value1")

	va, ok = pc.Vars["key2"]
	assert.True(t, ok)
	assert.Equal(t, va, "value2")
}

func TestMergeTemplateContext(t *testing.T) {
	p := CreateParameters(map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})

	others := CreateParameters(map[string]interface{}{
		"key2": "value2_overwritten",
		"key3": "value3",
	})

	pc := CreateTemplateContext(p)
	pc.Vars = pc.Vars.Override(others)
	assert.NotNil(t, pc)
	assert.Equal(t, 3, len(pc.Vars))
	va, ok := pc.Vars["key1"]
	assert.True(t, ok)
	assert.Equal(t, va, "value1")

	// Value should be overwritten
	va, ok = pc.Vars["key2"]
	assert.True(t, ok)
	assert.Equal(t, va, "value2_overwritten")

	// Missing stuff should be added
	va, ok = pc.Vars["key3"]
	assert.True(t, ok)
	assert.Equal(t, va, "value3")
}

func TestYaml(t *testing.T) {
	p := CreateParameters(map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})

	tplC := CreateTemplateContext(p)
	res, err := tplC.Execute(`{{ .Vars | yaml }}`)
	assert.Nil(t, err)
	assert.Equal(t, "key1: value1\nkey2: value2\n", res)
}

func TestIndentedYaml(t *testing.T) {
	p := CreateParameters(map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})

	tplC := CreateTemplateContext(p)
	res, err := tplC.Execute(`{{ .Vars | yaml | indent 4 }}`)
	assert.Nil(t, err)
	assert.Equal(t, "    key1: value1\n    key2: value2\n    \n", res)
}

func TestJson(t *testing.T) {
	p := CreateParameters(map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})

	tplC := CreateTemplateContext(p)
	res, err := tplC.Execute(`{{ .Vars | json }}`)
	assert.Nil(t, err)
	assert.Equal(t, "{\"key1\":\"value1\",\"key2\":\"value2\"}", res)
}

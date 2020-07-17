package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEngineFromBadHttp(t *testing.T) {
	tplC := CreateTemplateContext(CreateEmptyParameters())
	e := parseYaml("https://github.com/ekara-platform/engine/tree/master/testdata/DUMMY.yaml", tplC, &yamlEnvironment{})
	// an error occurred
	assert.NotNil(t, e)
}

func TestCreateEngineFromLocal(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	tplC := CreateTemplateContext(CreateEmptyParameters())
	e := parseYaml("testdata/yaml/complete.yaml", tplC, &yamlEnv)
	assert.Nil(t, e) // no error occurred

	assert.Equal(t, "testEnvironment", yamlEnv.Name)                              // importing file have has precedence
	assert.Equal(t, "This is my awesome Ekara environment.", yamlEnv.Description) // imported files are merged
}

func TestCreateEngineFromLocalComplexParams(t *testing.T) {
	yamlEnv := yamlEnvironment{}
	tplC := CreateTemplateContext(CreateEmptyParameters())
	e := parseYaml("testdata/yaml/complex.yaml", tplC, &yamlEnv)
	assert.Nil(t, e) // no error occurred
	assert.NotNil(t, yamlEnv)
}

func TestCreateEngineFromLocalWithData(t *testing.T) {
	vars := CreateParameters(map[string]interface{}{
		"info": map[string]string{
			"name": "Name from data",
			"desc": "Description from data",
		},
	})

	yamlEnv := yamlEnvironment{}
	tplC := CreateTemplateContext(vars)
	e := parseYaml("testdata/yaml/data.yaml", tplC, &yamlEnv)
	assert.Nil(t, e) // no error occurred
	assert.NotNil(t, yamlEnv)
	assert.Equal(t, "Name from data", yamlEnv.Name)
	assert.Equal(t, "Description from data", yamlEnv.Description)
}

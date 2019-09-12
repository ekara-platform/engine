package engine

import (
	"testing"

	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateRuntimeContext(t *testing.T) {

	p, _ := model.CreateParameters(map[string]interface{}{
		"key1": "val1",
		"key2": "val2",
	})

	c := &MockLaunchContext{data: p}

	rC := CreateRuntimeContext(c)
	assert.NotNil(t, rC)
	//The runtime context must be initialized with the param passed into the LauchContect
	if assert.Equal(t, 2, len(rC.data.Vars)) {
		cp(t, rC.data.Vars, "key1", "val1")
		cp(t, rC.data.Vars, "key2", "val2")
	}
}

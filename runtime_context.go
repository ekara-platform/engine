package engine

import (
	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
)

type (
	runtimeContext struct {
		ekaraError error
		report     ReportFileContent
		buffer     map[string]ansible.Buffer
		data       *model.TemplateContext
	}
)

func CreateRuntimeContext(lC LaunchContext) *runtimeContext {
	// Initialization of the runtime context
	rC := &runtimeContext{
		data: &model.TemplateContext{
			Vars: lC.ParamsFile(),
		},
	}

	rC.buffer = make(map[string]ansible.Buffer)
	return rC
}

func (c *runtimeContext) getBuffer(p *util.FolderPath) ansible.Buffer {
	// We check if we have a buffer corresponding to the provided folder path
	if val, ok := c.buffer[p.Path()]; ok {
		return val
	}
	return ansible.CreateBuffer()
}

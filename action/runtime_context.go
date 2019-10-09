package action

import (
	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/engine/util"
)

type (
	runtimeContext struct {
		lC util.LaunchContext
		cM component.Manager
		aM ansible.Manager

		ekaraError error
		report     ReportFileContent
		buffer     map[string]ansible.Buffer
		pN         util.ProgressNotifier
	}
)

//CreateRuntimeContext creates a new context for the runtime
func CreateRuntimeContext(lC util.LaunchContext, cM component.Manager, aM ansible.Manager, pN util.ProgressNotifier) *runtimeContext {
	// Initialization of the runtime context
	rC := &runtimeContext{
		lC: lC,
		cM: cM,
		aM: aM,
		pN: pN,
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

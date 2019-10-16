package action

import (
	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
)

type (
	runtimeContext struct {
		lC util.LaunchContext
		cF component.Finder
		aM ansible.Manager
		pN util.ProgressNotifier

		tplC        *model.TemplateContext
		environment *model.Environment
		report      ReportFileContent
		buffer      map[string]ansible.Buffer
	}
)

//createRuntimeContext creates a new context for the runtime
func createRuntimeContext(lC util.LaunchContext, tplC model.TemplateContext, env model.Environment, cF component.Finder, aM ansible.Manager, pN util.ProgressNotifier) *runtimeContext {
	// Initialization of the runtime context
	rC := &runtimeContext{
		lC:          lC,
		cF:          cF,
		aM:          aM,
		pN:          pN,
		environment: &env,
		tplC:        &tplC,
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

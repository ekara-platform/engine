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

		tplC        *model.TemplateContext
		environment *model.Environment
		report      ReportFileContent
	}
)

//createRuntimeContext creates a new context for the runtime
func createRuntimeContext(lC util.LaunchContext, cF component.Finder, aM ansible.Manager, env *model.Environment, tplC *model.TemplateContext) *runtimeContext {
	// Initialization of the runtime context
	rC := &runtimeContext{
		lC:          lC,
		cF:          cF,
		aM:          aM,
		environment: env,
		tplC:        tplC,
	}
	return rC
}

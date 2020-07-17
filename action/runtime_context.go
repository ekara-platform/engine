package action

import (
	"github.com/GroupePSA/componentizer"
	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/model"
	"github.com/ekara-platform/engine/util"
)

type (
	runtimeContext struct {
		lC util.LaunchContext
		cM componentizer.ComponentManager
		aM ansible.Manager

		tplC        componentizer.TemplateContext
		environment model.Environment
		report      ReportFileContent
		result      Result
	}
)

//createRuntimeContext creates a new context for the runtime
func createRuntimeContext(lC util.LaunchContext, cM componentizer.ComponentManager, aM ansible.Manager, env model.Environment, tplC componentizer.TemplateContext) *runtimeContext {
	// Initialization of the runtime context
	rC := &runtimeContext{
		lC:          lC,
		cM:          cM,
		aM:          aM,
		environment: env,
		tplC:        tplC,
	}
	return rC
}

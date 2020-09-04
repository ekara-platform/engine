package action

import (
	"github.com/GroupePSA/componentizer"
	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/model"
	"github.com/ekara-platform/engine/util"
)

type (
	RuntimeContext struct {
		lC util.LaunchContext
		cM componentizer.ComponentManager
		aM ansible.Manager

		tplC        componentizer.TemplateContext
		environment model.Environment
		result      Result
	}
)

//createRuntimeContext creates a new context for the runtime
func CreateRuntimeContext(lC util.LaunchContext, cM componentizer.ComponentManager, aM ansible.Manager, env model.Environment, tplC componentizer.TemplateContext) *RuntimeContext {
	// Initialization of the runtime context
	rC := &RuntimeContext{
		lC:          lC,
		cM:          cM,
		aM:          aM,
		environment: env,
		tplC:        tplC,
	}
	return rC
}

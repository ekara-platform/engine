package engine

import (
	"fmt"
	"strings"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
)

type (

	// Hook represents tasks to be executed linked to an ekara life cycle event
	hookContext struct {
		action    string
		target    model.Describable
		hookOnwer string
		hookName  string
		baseParam ansible.BaseParam
		envVar    ansible.EnvVars
		buffer    ansible.Buffer
	}
)

//RunHookBefore Runs the hooks defined to be executed before a task.
func RunHookBefore(cm component.ComponentManager, lC LaunchContext, rC *runtimeContext, r *StepResults, h model.Hook, ctx hookContext, cl Cleanup) {
	runHooks(cm, h.Before, lC, rC, r, ctx, cl)
}

//RunHookAfter Runs the hooks defined to be executed after a task.
func RunHookAfter(cm component.ComponentManager, lC LaunchContext, rC *runtimeContext, r *StepResults, h model.Hook, ctx hookContext, cl Cleanup) {
	runHooks(cm, h.After, lC, rC, r, ctx, cl)
}

func runHooks(cm component.ComponentManager, hooks []model.TaskRef, lC LaunchContext, rC *runtimeContext, r *StepResults, ctx hookContext, cl Cleanup) {
	for i, hook := range hooks {
		repName := fmt.Sprintf("%s_%s_hook_%s_%s_%s_%d", ctx.action, ctx.target.DescName(), ctx.hookOnwer, ctx.hookName, hook.HookLocation, i)
		sc := InitHookStepResult(folderAsMessage(repName), ctx.target, cl)
		ef, ko := createChildExchangeFolder(lC.Ef().Input, repName, &sc, lC.Log())
		if ko {
			r.Add(sc)
			continue
		}

		t, err := hook.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the task"), nil)
			r.Add(sc)
			continue
		}

		bp := ctx.baseParam.Copy()
		bp.AddNamedMap("hook_param", t.Parameters)

		if ko := saveBaseParams(bp, lC, ef.Input, &sc); ko {
			r.Add(sc)
			continue
		}

		// Prepare components map
		if ko := saveComponentMap(lC, ef.Input, &sc); ko {
			r.Add(sc)
			continue
		}

		exv := ansible.BuildExtraVars("", *ef.Input, *ef.Output, ctx.buffer)

		runTask(cm, lC, rC, t, ctx.target, sc, r, ef, exv, ctx.envVar)
		// TODO Consume hook buffer here
	}
}

func runTask(cm component.ComponentManager, lC LaunchContext, rC *runtimeContext, task model.Task, target model.Describable, sc StepResult, r *StepResults, ef *util.ExchangeFolder, exv ansible.ExtraVars, env ansible.EnvVars) {
	usable, err := cm.Use(task)
	if err != nil {
		FailsOnCode(&sc, err, "An error occurred getting the usable task", nil)
	}
	defer usable.Release()
	code, err := lC.Ekara().AnsibleManager().Execute(usable, task.Playbook, exv, env)
	if err != nil {
		pfd := playBookFailureDetail{
			Playbook:  task.Playbook,
			Component: task.ComponentName(),
			Code:      code,
		}
		FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
		r.Add(sc)
		return
	}
	r.Add(sc)
}

func folderAsMessage(s string) string {
	r := strings.Replace(s, "_", " ", -1)
	return strings.Title(r[:1]) + r[1:]
}

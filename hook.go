package engine

import (
	"fmt"
	"strings"

	"github.com/ekara-platform/engine/ansible"
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
		inventory string
	}
)

func RunHookBefore(lC LaunchContext, rC *runtimeContext, r *StepResults, h model.Hook, ctx hookContext, cl Cleanup) {
	//previousBuffer := ansible.Buffer{}
	for i, hook := range h.Before {
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

		runTask(lC, rC, t, ctx.target, sc, r, ef, exv, ctx.envVar, ctx.inventory)
		// TODO Consume hook buffer here
	}
}

func RunHookAfter(lC LaunchContext, rC *runtimeContext, r *StepResults, h model.Hook, ctx hookContext, cl Cleanup) {
	//previousBuffer := ansible.Buffer{}
	for i, hook := range h.After {
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

		runTask(lC, rC, t, ctx.target, sc, r, ef, exv, ctx.envVar, ctx.inventory)
		// TODO Consume hook buffer here
	}
}

func runTask(lC LaunchContext, rC *runtimeContext, task model.Task, target model.Describable, sc StepResult, r *StepResults, ef *util.ExchangeFolder, exv ansible.ExtraVars, env ansible.EnvVars, inventory string) {

	comp, err := task.Component.Resolve()
	if err != nil {
		FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the task's component"), nil)
		r.Add(sc)
		return
	}

	err, code := lC.Ekara().AnsibleManager().Execute(comp, task.Playbook, exv, env, inventory)
	if err != nil {
		pfd := playBookFailureDetail{
			Playbook:  task.Playbook,
			Compoment: comp.Id,
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

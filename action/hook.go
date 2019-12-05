package action

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
		hookOwner string
		hookName  string
		baseParam ansible.BaseParam
	}
)

//RunHookBefore Runs the hooks defined to be executed before a task.
func runHookBefore(rC *runtimeContext, r *StepResults, h model.Hook, ctx hookContext, cl Cleanup) {
	runHooks(h.Before, rC, r, ctx, cl)
}

//RunHookAfter Runs the hooks defined to be executed after a task.
func runHookAfter(rC *runtimeContext, r *StepResults, h model.Hook, ctx hookContext, cl Cleanup) {
	runHooks(h.After, rC, r, ctx, cl)
}

func runHooks(hooks []model.TaskRef, rC *runtimeContext, r *StepResults, ctx hookContext, cl Cleanup) {
	for i, hook := range hooks {
		repName := fmt.Sprintf("%s_%s_hook_%s_%s_%s_%d", ctx.action, ctx.target.DescName(), ctx.hookOwner, ctx.hookName, hook.HookLocation, i)
		sc := InitHookStepResult(folderAsMessage(repName), ctx.target, cl)
		ef, ko := createChildExchangeFolder(rC.lC.Ef().Input, repName, &sc)
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

		if ko := saveBaseParams(bp, ef.Input, &sc); ko {
			r.Add(sc)
			continue
		}

		exv := ansible.CreateExtraVars(ef.Input, ef.Output)

		runTask(rC, t, sc, r, exv)
		r.Add(fconsumeHookResult(rC, ctx.target, repName, ef))
	}
}

func fconsumeHookResult(rC *runtimeContext, target model.Describable, repName string, ef util.ExchangeFolder) StepResult {
	sc := InitCodeStepResult("Consuming the hook result", target, NoCleanUpRequired)
	rC.lC.Log().Printf("Consuming the hook result %s", target.DescName())
	// TODO
	return sc
}

func runTask(rC *runtimeContext, task model.Task, sc StepResult, r *StepResults, exv ansible.ExtraVars) {
	usable, err := rC.cF.Use(task, rC.tplC)
	if err != nil {
		FailsOnCode(&sc, err, "An error occurred getting the usable task", nil)
	}
	defer usable.Release()

	code, err := rC.aM.Play(usable, rC.tplC, task.Playbook, exv)
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

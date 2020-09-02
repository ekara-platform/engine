package action

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/model"
	"github.com/ekara-platform/engine/util"
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
	runHooks(h.Before, rC, r, ctx, cl, "before")
}

//RunHookAfter Runs the hooks defined to be executed after a task.
func runHookAfter(rC *runtimeContext, r *StepResults, h model.Hook, ctx hookContext, cl Cleanup) {
	runHooks(h.After, rC, r, ctx, cl, "after")
}

func runHooks(hooks []model.TaskRef, rC *runtimeContext, r *StepResults, ctx hookContext, cl Cleanup, phase string) {
	for i, hook := range hooks {

		repName := fmt.Sprintf("%s_%s_hook_%s_%s_%s_%d", ctx.action, ctx.target.DescName(), ctx.hookOwner, ctx.hookName, phase, i)
		sc := InitHookStepResult(folderAsMessage(repName), ctx.target, cl)
		ef, ko := createChildExchangeFolder(rC.lC.Ef().Input, repName, &sc)
		if ko {
			r.Add(sc)
			continue
		}

		t, err := hook.Resolve(rC.environment)
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the task"), nil)
			r.Add(sc)
			continue
		}

		bp := ctx.baseParam.Copy()
		bp.AddNamedMap("params", t.Parameters())

		if ko := saveBaseParams(bp, ef.Input, &sc); ko {
			r.Add(sc)
			continue
		}

		exv := ansible.CreateExtraVars(ef.Input, ef.Output)
		runTask(rC, t, sc, r, exv)
		r.Add(fconsumeHookResult(rC, ctx.target, ctx, ef, hook.Prefix))
	}
}

func fconsumeHookResult(rC *runtimeContext, target model.Describable, ctx hookContext, ef util.ExchangeFolder, prefix string) StepResult {
	sc := InitCodeStepResult("Consuming the hook result", target, NoCleanUpRequired)

	if ef.Output.Contains(util.OutputYamlFileName) {
		filePath := util.JoinPaths(ef.Output.Path(), util.OutputYamlFileName)
		b, err := util.FileRead(filePath)

		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred reading the output.yaml"), nil)
			return sc
		}
		content := make(map[string]interface{})
		err = yaml.Unmarshal(b, content)
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred Unmarshalling the output.yaml"), nil)
			return sc
		}

		if prefix != "" {
			rC.tplC.(model.TemplateContext).Runtime[prefix] = content
		} else {
			if val, ok := rC.tplC.(model.TemplateContext).Runtime[ctx.action]; ok {
				vv := reflect.ValueOf(val)
				if vv.Kind() == reflect.Map {
					ma := val.(map[string]interface{})
					ma[ctx.target.DescName()] = content
				}
			} else {
				m := make(map[string]interface{})
				m[ctx.target.DescName()] = content
				rC.tplC.(model.TemplateContext).Runtime[ctx.action] = m
			}
		}
	}
	return sc
}

func runTask(rC *runtimeContext, task model.Task, sc StepResult, r *StepResults, exv ansible.ExtraVars) {
	usable, err := rC.cM.Use(task, rC.tplC)
	if err != nil {
		FailsOnCode(&sc, err, "An error occurred getting the usable task", nil)
	}
	defer usable.Release()

	code, err := rC.aM.Play(usable, rC.tplC, task.Playbook, exv)
	if err != nil {
		pfd := playBookFailureDetail{
			Playbook:  task.Playbook,
			Component: task.ComponentId(),
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

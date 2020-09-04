package action

import (
	"encoding/json"
	"fmt"
	"github.com/ekara-platform/engine/ansible"
)

const (
	destroyPlaybook = "destroy.yaml"
)

type (
	//DestroyResult contains the results of environment destruction
	DestroyResult struct {
		Success bool
	}
)

func (r DestroyResult) IsSuccess() bool {
	return r.Success
}

func (r *DestroyResult) FromJson(s string) error {
	err := json.Unmarshal([]byte(s), r)
	if err != nil {
		return err
	}
	return nil
}

func (r DestroyResult) AsJson() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

var (
	destroyAction = Action{
		DestroyActionID,
		CheckActionID,
		"Destroy",
		[]step{providerSetup, destroyHookBefore, providerDestroy, destroyHookAfter},
	}
)

func destroyHookBefore(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()

	if len(rC.environment.Hooks.Destroy.Before) == 0 {
		return *sCs
	}

	// Notify destruction progress
	rC.lC.Feedback().ProgressG("destroy.hook.before", 1, "Hook before destroying node sets")

	// Prepare parameters
	bp := buildBaseParam(rC, "")

	// Process hook : environment - destroy - before
	runHookBefore(
		rC,
		sCs,
		rC.environment.Hooks.Destroy,
		hookContext{"destroy", rC.environment, "environment", "destroy", bp},
		NoCleanUpRequired,
	)

	rC.lC.Feedback().Progress("destroy.hook.before", "All hooks executed")
	return *sCs
}

func providerDestroy(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()

	for _, n := range rC.environment.NodeSets {
		sc := InitPlaybookStepResult("Running the destroy phase", n, NoCleanUpRequired)

		// Resolve provider
		p, err := n.Provider.Resolve(rC.environment)
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the provider"), nil)
			sCs.Add(sc)
			return *sCs
		}

		// Notify destruction progress
		rC.lC.Feedback().ProgressG("provider.destroy", len(rC.environment.NodeSets), "Destroying node set '%s' with provider '%s'", n.Name, p.Name)

		// Prepare parameters
		bp := buildBaseParam(rC, n.Name)
		bp.AddInt("instances", n.Instances)
		bp.AddInterface("labels", n.Labels)
		bp.AddNamedMap("params", p.Parameters())
		bp.AddInterface("proxy", p.Proxy())

		// Process hook : nodeset - destroy - before
		runHookBefore(
			rC,
			sCs,
			n.Hooks.Destroy,
			hookContext{"destroy", n, "nodeset", "destroy", bp},
			NoCleanUpRequired,
		)

		// Node destruction exchange folder
		destroy, ko := createChildExchangeFolder(rC.lC.Ef().Input, "destroy_"+n.Name, &sc)
		if ko {
			sCs.Add(sc)
			return *sCs
		}

		if ko := saveBaseParams(bp, destroy.Input, &sc); ko {
			sCs.Add(sc)
			return *sCs
		}

		// Prepare extra vars
		exv := ansible.CreateExtraVars(destroy.Input, destroy.Output)

		// Make the component usable
		usable, err := rC.cM.Use(p, rC.tplC)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable provider", nil)
		}
		defer usable.Release()

		// Launch the playbook
		code, err := rC.aM.Play(usable, rC.tplC, destroyPlaybook, exv)

		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  destroyPlaybook,
				Component: p.ComponentId(),
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
			sCs.Add(sc)
			return *sCs
		}
		sCs.Add(sc)

		// Process hook : nodeset - destroy - after
		runHookAfter(
			rC,
			sCs,
			n.Hooks.Destroy,
			hookContext{"destroy", n, "nodeset", "destroy", bp},
			NoCleanUpRequired,
		)
	}

	// Notify destruction finish
	rC.lC.Feedback().Progress("provider.destroy", "All node sets destroyed")

	return *sCs
}

func destroyHookAfter(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()

	if len(rC.environment.Hooks.Destroy.After) == 0 {
		return *sCs
	}

	// Notify destruction progress
	rC.lC.Feedback().ProgressG("destroy.hook.after", 1, "Hook after destruction node sets")

	// Prepare parameters
	bp := buildBaseParam(rC, "")

	// Process hook : environment - destroy - after
	runHookAfter(
		rC,
		sCs,
		rC.environment.Hooks.Destroy,
		hookContext{"destroy", rC.environment, "environment", "destroy", bp},
		NoCleanUpRequired,
	)

	rC.lC.Feedback().Progress("provider.destroy.hook.after", "All hooks executed")
	return *sCs
}

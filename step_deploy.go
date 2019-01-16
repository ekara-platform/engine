package engine

import (
	"fmt"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/model"
)

var deploySteps = []step{fstack}

func fstack(lC LaunchContext, rC *runtimeContext) StepResults {
	sCs := InitStepResults()
	for _, s := range lC.Ekara().ComponentManager().Environment().Stacks {

		sc := InitCodeStepResult("Starting a stack setup phase", s, NoCleanUpRequired)
		lC.Log().Printf("Checking how to install %s", s.Name)
		// Check if the stacks holds an "install.yml" playbook
		r, err := s.Component.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occured resolving the stack"), nil)
			sCs.Add(sc)
			continue
		}
		sCs.Add(sc)
		if ok := lC.Ekara().AnsibleManager().Contains(r, "install.yml"); ok {
			fstackPlabook(lC, rC, s, sCs)
		} else {
			fstackCompose(lC, rC, s, sCs)
		}
	}
	return *sCs
}

func fstackPlabook(lC LaunchContext, rC *runtimeContext, s model.Stack, sCs *StepResults) {
	for _, n := range lC.Ekara().ComponentManager().Environment().NodeSets {

		sc := InitPlaybookStepResult("Running the stack playbook installation phase", model.ChainDescribable(s, n), NoCleanUpRequired)
		p, err := n.Provider.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occured resolving the nodeset"), nil)
			sCs.Add(sc)
			continue
		}

		// Provider setup exchange folder
		setupProviderEf := lC.Ef().Input.Children["setup_provider_"+p.Name]
		// We check if we have a buffer corresponding to the provider setup
		bufferPro := rC.getBuffer(setupProviderEf.Output)

		inventory := ""
		if len(bufferPro.Inventories) > 0 {
			inventory = bufferPro.Inventories["inventory_path"]
		}

		// We use an empty buffer because no one is coming from the previous step
		buffer := ansible.CreateBuffer()

		bp := BuildBaseParam(lC, n.Name, p.Name)

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.Add("http_proxy", lC.HttpProxy())
		env.Add("https_proxy", lC.HttpsProxy())
		env.Add("no_proxy", lC.NoProxy())
		env.AddBuffer(buffer)

		// ugly but .... TODO change this
		env.AddBuffer(bufferPro)

		// Process hook : environment - deploy - before
		RunHookBefore(lC,
			rC,
			sCs,
			lC.Ekara().ComponentManager().Environment().Hooks.Deploy,
			hookContext{"install", s, "environment", "deploy", bp, env, buffer, inventory},
			NoCleanUpRequired,
		)

		// Process hook : stack - deploy - before
		RunHookBefore(lC,
			rC,
			sCs,
			s.Hooks.Deploy,
			hookContext{"install", s, "stack", "deploy", bp, env, buffer, inventory},
			NoCleanUpRequired,
		)

		// Stack install exchange folder
		fName := fmt.Sprintf("install_stack_%s_on_%s", s.Name, n.Name)

		stackEf, ko := createChildExchangeFolder(lC.Ef().Input, fName, &sc, lC.Log())
		if ko {
			sCs.Add(sc)
			continue
		}

		// Prepare extra vars
		exv := ansible.BuildExtraVars("", *stackEf.Input, *stackEf.Output, buffer)

		// We launch the playbook
		r, err := s.Component.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occured resolving the stack"), nil)
			sCs.Add(sc)
			continue
		}
		err, code := lC.Ekara().AnsibleManager().Execute(r, "install.yml", exv, env, inventory)
		if err != nil {
			r, err := p.Component.Resolve()
			if err != nil {
				FailsOnCode(&sc, err, fmt.Sprintf("An error occured resolving the provider"), nil)
				sCs.Add(sc)
				continue
			}
			pfd := playBookFailureDetail{
				Playbook:  "install.yml",
				Compoment: r.Id,
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occured executing the playbook", pfd)
			sCs.Add(sc)
			continue
		}
		if ko {
			sCs.Add(sc)
			continue
		}

		sCs.Add(sc)

		// Process hook : stack - deploy - after
		RunHookAfter(lC,
			rC,
			sCs,
			s.Hooks.Deploy,
			hookContext{"install", s, "stack", "deploy", bp, env, buffer, inventory},
			NoCleanUpRequired,
		)

		// Process hook : environment - deploy - after
		RunHookAfter(lC,
			rC,
			sCs,
			lC.Ekara().ComponentManager().Environment().Hooks.Deploy,
			hookContext{"install", s, "environment", "deploy", bp, env, buffer, inventory},
			NoCleanUpRequired,
		)
	}
}

func fstackCompose(lC LaunchContext, rC *runtimeContext, s model.Stack, sCs *StepResults) {
	for _, n := range lC.Ekara().ComponentManager().Environment().Stacks {

		// TODO add hooks process here
		// Process hook : environment - deploy - before
		// Process hook : stack - deploy - before

		sc := InitPlaybookStepResult("Running the stack Docker Compose installation phase", model.ChainDescribable(s, n), NoCleanUpRequired)
		lC.Log().Printf(LOG_PROCESSING_STACK_COMPOSE, s.Name, n.Name)
		err := fmt.Errorf("The Ekara plaform is not ready yet to install stacks using Docker Compose")
		FailsOnNotImplemented(&sc, err, "", nil)
		sCs.Add(sc)

		// TODO add hooks process here
		// Process hook : stack - deploy - before
		// Process hook : environment - deploy - before
	}
}

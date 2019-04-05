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
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the stack"), nil)
			sCs.Add(sc)
			continue
		}
		sCs.Add(sc)
		if ok := lC.Ekara().AnsibleManager().Contains(r, "install.yml"); ok {
			fstackPlabook(lC, rC, s, sCs)
		} else {
			fstackCompose(lC, rC, lC.Ekara().ComponentManager().Environment().Ekara.Distribution, s, sCs)
		}
	}
	return *sCs
}

func fstackPlabook(lC LaunchContext, rC *runtimeContext, s model.Stack, sCs *StepResults) {
	// Map to keep trace on the processed providers
	ps := make(map[string]model.Provider, len(lC.Ekara().ComponentManager().Environment().NodeSets))
	for _, n := range lC.Ekara().ComponentManager().Environment().NodeSets {

		sc := InitPlaybookStepResult("Running the stack playbook installation phase", model.ChainDescribable(s, n), NoCleanUpRequired)

		p, err := n.Provider.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the provider"), nil)
			sCs.Add(sc)
			continue
		}

		if _, ok := ps[p.Name]; ok {
			continue // The provider has already been processed
		}
		ps[p.Name] = p

		// Provider setup exchange folder
		setupProviderEf := lC.Ef().Input.Children["setup_provider_"+p.Name]
		// We check if we have a buffer corresponding to the provider setup
		bufferPro := rC.getBuffer(setupProviderEf.Output)

		// Stack install exchange folder for the given provider
		fName := fmt.Sprintf("install_stack_%s_for_%s", s.Name, p.Name)

		stackEf, ko := createChildExchangeFolder(lC.Ef().Input, fName, &sc, lC.Log())
		if ko {
			sCs.Add(sc)
			continue
		}
		// We use the provider inventory
		inventory := ""
		if len(bufferPro.Inventories) > 0 {
			inventory = bufferPro.Inventories["inventory_path"]
		}

		// We use an empty buffer because no one is coming from the previous step
		buffer := ansible.CreateBuffer()

		bp := BuildBaseParam(lC, "", p.Name)
		bp.AddNamedMap("params", s.Parameters)

		if ko := saveBaseParams(bp, lC, stackEf.Input, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare components map
		if ko := saveComponentMap(lC, stackEf.Input, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.Add("http_proxy", lC.HttpProxy())
		env.Add("https_proxy", lC.HttpsProxy())
		env.Add("no_proxy", lC.NoProxy())

		// Adding the environment variables from the stack
		for envK, envV := range s.EnvVars {
			env.Add(envK, envV)
		}

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

		// Prepare extra vars
		exv := ansible.BuildExtraVars("", *stackEf.Input, *stackEf.Output, buffer)

		// We launch the playbook
		r, err := s.Component.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the stack"), nil)
			sCs.Add(sc)
			continue
		}
		err, code := lC.Ekara().AnsibleManager().Execute(r, "install.yml", exv, env, inventory)
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  "install.yml",
				Compoment: r.Id,
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
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

func fstackCompose(lC LaunchContext, rC *runtimeContext, coreStack model.Component, s model.Stack, sCs *StepResults) {

	// Map to keep trace on the processed providers
	ps := make(map[string]model.Provider, len(lC.Ekara().ComponentManager().Environment().NodeSets))
	for _, n := range lC.Ekara().ComponentManager().Environment().NodeSets {

		sc := InitPlaybookStepResult("Running the stack compose installation phase", model.ChainDescribable(s, n), NoCleanUpRequired)

		p, err := n.Provider.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the provider"), nil)
			sCs.Add(sc)
			continue
		}

		if _, ok := ps[p.Name]; ok {
			continue // The provider has already been processed
		}
		ps[p.Name] = p

		// Provider setup exchange folder
		setupProviderEf := lC.Ef().Input.Children["setup_provider_"+p.Name]
		// We check if we have a buffer corresponding to the provider setup
		bufferPro := rC.getBuffer(setupProviderEf.Output)

		// Stack install exchange folder for the given provider
		fName := fmt.Sprintf("install_stack_%s_for_%s", s.Name, p.Name)

		stackEf, ko := createChildExchangeFolder(lC.Ef().Input, fName, &sc, lC.Log())
		if ko {
			sCs.Add(sc)
			continue
		}
		// We use the provider inventory
		inventory := ""
		if len(bufferPro.Inventories) > 0 {
			inventory = bufferPro.Inventories["inventory_path"]
		}

		// We use an empty buffer because no one is coming from the previous step
		buffer := ansible.CreateBuffer()

		bp := BuildBaseParam(lC, "", p.Name)
		bp.AddNamedMap("params", s.Parameters)

		if ko := saveBaseParams(bp, lC, stackEf.Input, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare components map
		if ko := saveComponentMap(lC, stackEf.Input, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.Add("http_proxy", lC.HttpProxy())
		env.Add("https_proxy", lC.HttpsProxy())
		env.Add("no_proxy", lC.NoProxy())

		// Adding the environment variables from the stack
		for envK, envV := range s.EnvVars {
			env.Add(envK, envV)
		}

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

		// We launch the playbook
		r, err := s.Component.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the stack"), nil)
			sCs.Add(sc)
			continue
		}

		// Prepare extra vars
		exv := ansible.BuildExtraVars("compose_path="+lC.Ekara().ComponentManager().ComponentPath(r.Id), *stackEf.Input, *stackEf.Output, buffer)

		err, code := lC.Ekara().AnsibleManager().Execute(coreStack, "install.yml", exv, env, inventory)
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  "install.yml",
				Compoment: r.Id,
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
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

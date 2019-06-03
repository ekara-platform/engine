package engine

import (
	"fmt"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/model"
)

var deploySteps = []step{ftemplate, fstack}

func ftemplate(lC LaunchContext, rC *runtimeContext) StepResults {
	sCs := InitStepResults()
	cpMan := lC.Ekara().ComponentManager()
	for _, s := range cpMan.Environment().Stacks {
		sc := InitCodeStepResult("Starting a stack template phase", s, NoCleanUpRequired)
		if len(s.Templates.Content) > 0 {
			lC.Log().Printf("The stack %s has templates to process ", s.Name)
			runTemplate(lC, sCs, &sc, s.Templates, s)
		} else {
			lC.Log().Printf("The stack %s has no templates to process ", s.Name)
		}
		sCs.Add(sc)
	}
	return *sCs
}

func fstack(lC LaunchContext, rC *runtimeContext) StepResults {
	sCs := InitStepResults()
	for _, s := range lC.Ekara().ComponentManager().Environment().Stacks {
		sc := InitCodeStepResult("Starting a stack setup phase", s, NoCleanUpRequired)
		lC.Log().Printf("Checking how to install %s", s.Name)
		sCs.Add(sc)
		if ok := lC.Ekara().AnsibleManager().Contains(s, "install.yml"); ok {
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
		env.Add("http_proxy", lC.HTTPProxy())
		env.Add("https_proxy", lC.HTTPSProxy())
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

		code, err := lC.Ekara().AnsibleManager().Execute(s, "install.yml", exv, env, inventory)
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  "install.yml",
				Component: s.ComponentName(),
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

// TODO currently only the distribution is able to deploy a compose...
func fstackCompose(lC LaunchContext, rC *runtimeContext, distribution model.Distribution, s model.Stack, sCs *StepResults) {

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
		env.Add("http_proxy", lC.HTTPProxy())
		env.Add("https_proxy", lC.HTTPSProxy())
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
		exv := ansible.BuildExtraVars("compose_path="+lC.Ekara().ComponentManager().ComponentPath(s.ComponentName()), *stackEf.Input, *stackEf.Output, buffer)

		code, err := lC.Ekara().AnsibleManager().Execute(distribution, "deploy_compose.yaml", exv, env, inventory)
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  "install.yml",
				Component: distribution.ComponentName(),
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

package engine

import (
	"fmt"
	"path/filepath"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/model"
)

var deploySteps = []step{fstack}

func fstack(lC LaunchContext, rC *runtimeContext) StepResults {
	cm := lC.Ekara().ComponentManager()
	sCs := InitStepResults()
	for _, st := range cm.Environment().Stacks {
		sc := InitCodeStepResult("Starting a stack setup phase", st, NoCleanUpRequired)
		lC.Log().Printf("Checking how to install %s", st.Name)
		sCs.Add(sc)

		ust, err := cm.Use(st)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable stack", nil)
		}
		defer ust.Release()
		if ok, _ := ust.ContainsFile("install.yml"); ok {
			fstackPlabook(lC, rC, st, ust, sCs)
		} else {
			fstackCompose(lC, rC, lC.Ekara().ComponentManager().Environment().Ekara.Distribution, st, sCs)
		}
	}
	return *sCs
}

func fstackPlabook(lC LaunchContext, rC *runtimeContext, st model.Stack, ust component.UsableComponent, sCs *StepResults) {
	// Map to keep trace on the processed providers
	cm := lC.Ekara().ComponentManager()
	ps := make(map[string]model.Provider, len(cm.Environment().NodeSets))
	for _, n := range cm.Environment().NodeSets {

		sc := InitPlaybookStepResult("Running the stack playbook installation phase", model.ChainDescribable(st, n), NoCleanUpRequired)

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

		// Stack install exchange folder for the given provider
		fName := fmt.Sprintf("install_stack_%s_for_%s", st.Name, p.Name)

		stackEf, ko := createChildExchangeFolder(lC.Ef().Input, fName, &sc, lC.Log())
		if ko {
			sCs.Add(sc)
			continue
		}

		// We use an empty buffer because no one is coming from the previous step
		buffer := ansible.CreateBuffer()

		bp := BuildBaseParam(lC, "", p.Name)
		bp.AddNamedMap("params", st.Parameters)

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
		env.AddDefaultOsVars()
		env.AddProxy(p.Proxy)

		// Adding the environment variables from the stack
		for envK, envV := range st.EnvVars {
			env.Add(envK, envV)
		}

		// Process hook : environment - deploy - before
		RunHookBefore(cm,
			lC,
			rC,
			sCs,
			cm.Environment().Hooks.Deploy,
			hookContext{"install", st, "environment", "deploy", bp, env, buffer},
			NoCleanUpRequired,
		)

		// Process hook : stack - deploy - before
		RunHookBefore(cm,
			lC,
			rC,
			sCs,
			st.Hooks.Deploy,
			hookContext{"install", st, "stack", "deploy", bp, env, buffer},
			NoCleanUpRequired,
		)

		// Prepare extra vars
		exv := ansible.BuildExtraVars("", *stackEf.Input, *stackEf.Output, buffer)

		code, err := lC.Ekara().AnsibleManager().Execute(ust, "install.yml", exv, env)
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  "install.yml",
				Component: st.ComponentName(),
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
		RunHookAfter(cm,
			lC,
			rC,
			sCs,
			st.Hooks.Deploy,
			hookContext{"install", st, "stack", "deploy", bp, env, buffer},
			NoCleanUpRequired,
		)

		// Process hook : environment - deploy - after
		RunHookAfter(cm,
			lC,
			rC,
			sCs,
			cm.Environment().Hooks.Deploy,
			hookContext{"install", st, "environment", "deploy", bp, env, buffer},
			NoCleanUpRequired,
		)
	}
}

// TODO currently only the distribution is able to deploy a compose...
// TODO refactor this method in order to run the docker stack deploy only once for the whole cluster,
// regardless of the number of providers...
func fstackCompose(lC LaunchContext, rC *runtimeContext, distribution model.Distribution, s model.Stack, sCs *StepResults) {
	cm := lC.Ekara().ComponentManager()
	// Map to keep trace on the processed providers
	ps := make(map[string]model.Provider, len(cm.Environment().NodeSets))
	for _, n := range cm.Environment().NodeSets {

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

		// Stack install exchange folder for the given provider
		fName := fmt.Sprintf("install_stack_%s_for_%s", s.Name, p.Name)

		stackEf, ko := createChildExchangeFolder(lC.Ef().Input, fName, &sc, lC.Log())
		if ko {
			sCs.Add(sc)
			continue
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
		env.AddDefaultOsVars()
		env.AddProxy(p.Proxy)

		// Adding the environment variables from the stack
		for envK, envV := range s.EnvVars {
			env.Add(envK, envV)
		}

		// Process hook : environment - deploy - before
		RunHookBefore(cm,
			lC,
			rC,
			sCs,
			cm.Environment().Hooks.Deploy,
			hookContext{"install", s, "environment", "deploy", bp, env, buffer},
			NoCleanUpRequired,
		)

		// Process hook : stack - deploy - before
		RunHookBefore(cm,
			lC,
			rC,
			sCs,
			s.Hooks.Deploy,
			hookContext{"install", s, "stack", "deploy", bp, env, buffer},
			NoCleanUpRequired,
		)

		// Prepare extra vars
		var exv ansible.ExtraVars

		d, err := cm.Use(distribution)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable distribution", nil)
		}
		defer d.Release()
		su, err := cm.Use(s)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable stack", nil)
		}
		defer su.Release()
		if ok, match := su.ContainsFile("docker-compose.yml"); ok {
			exv = ansible.BuildExtraVars("compose_path="+filepath.Join(match.Component().RootPath(), match.RelativePath()), *stackEf.Input, *stackEf.Output, buffer)
		} else {
			exv = ansible.BuildExtraVars("", *stackEf.Input, *stackEf.Output, buffer)
		}

		code, err := lC.Ekara().AnsibleManager().Execute(d, "deploy_compose.yaml", exv, env)
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  "deploy_compose.yaml",
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
		RunHookAfter(cm,
			lC,
			rC,
			sCs,
			s.Hooks.Deploy,
			hookContext{"install", s, "stack", "deploy", bp, env, buffer},
			NoCleanUpRequired,
		)

		// Process hook : environment - deploy - after
		RunHookAfter(cm,
			lC,
			rC,
			sCs,
			cm.Environment().Hooks.Deploy,
			hookContext{"install", s, "environment", "deploy", bp, env, buffer},
			NoCleanUpRequired,
		)
	}
}

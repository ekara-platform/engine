package engine

import (
	"fmt"
	"path/filepath"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/model"
)

var deploySteps = []step{fstack}

// TODO: refactor this as most of playbook/compose deployment is the same
func fstack(lC LaunchContext, rC *runtimeContext) StepResults {
	cm := lC.Ekara().ComponentManager()
	sCs := InitStepResults()
	for _, st := range cm.Environment().Stacks {
		sc := InitCodeStepResult("Deploying stack", st, NoCleanUpRequired)
		lC.Log().Printf("Checking how to deploy stack %s, for component %s", st.Name, st.ComponentName())
		sCs.Add(sc)

		// Make the stack usable
		ust, err := cm.Use(st, rC.data)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable stack", nil)
		}
		defer ust.Release()

		if ok, _ := ust.ContainsFile("install.yml"); ok {
			fstackPlaybook(lC, rC, st, ust, sCs)
		} else {
			fstackCompose(lC, rC, st, ust, sCs)
		}
	}
	return *sCs
}

func fstackPlaybook(lC LaunchContext, rC *runtimeContext, st model.Stack, ust component.UsableComponent, sCs *StepResults) {
	cm := lC.Ekara().ComponentManager()
	sc := InitPlaybookStepResult("Deploying stack with a custom playbook", st, NoCleanUpRequired)

	// Stack install exchange folder for the given provider
	fName := fmt.Sprintf("install_stack_%s", st.Name)

	stackEf, ko := createChildExchangeFolder(lC.Ef().Input, fName, &sc, lC.Log())
	if ko {
		sCs.Add(sc)
		return
	}

	// Create a new buffer
	buffer := ansible.CreateBuffer()

	// Prepare parameters
	bp := BuildBaseParam(lC, "")
	bp.AddNamedMap("params", st.Parameters)
	if ko := saveBaseParams(bp, lC, stackEf.Input, &sc); ko {
		sCs.Add(sc)
		return
	}

	// Prepare environment variables
	env := ansible.BuildEnvVars()
	env.AddDefaultOsVars()
	env.AddProxy(lC.Proxy())

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

	code, err := lC.Ekara().AnsibleManager().Execute(ust, "install.yml", exv, env, rC.data)
	if err != nil {
		pfd := playBookFailureDetail{
			Playbook:  "install.yml",
			Component: st.ComponentName(),
			Code:      code,
		}
		FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
		sCs.Add(sc)
		return
	}

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

	sCs.Add(sc)
}

// TODO currently only the parent is able to deploy a compose...
func fstackCompose(lC LaunchContext, rC *runtimeContext, st model.Stack, ust component.UsableComponent, sCs *StepResults) {
	cm := lC.Ekara().ComponentManager()
	rm := lC.Ekara().ReferenceManager()
	sc := InitPlaybookStepResult("aaaaa Deploying stack with a Docker compose", st, NoCleanUpRequired)

	// Stack install exchange folder for the given provider
	fName := fmt.Sprintf("install_stack_%s", st.Name)

	stackEf, ko := createChildExchangeFolder(lC.Ef().Input, fName, &sc, lC.Log())
	if ko {
		sCs.Add(sc)
		return
	}

	// Create a new buffer
	buffer := ansible.CreateBuffer()

	// Prepare parameters
	bp := BuildBaseParam(lC, "")
	bp.AddNamedMap("params", st.Parameters)
	if ko := saveBaseParams(bp, lC, stackEf.Input, &sc); ko {
		sCs.Add(sc)
		return
	}

	// Prepare environment variables
	env := ansible.BuildEnvVars()
	env.AddDefaultOsVars()
	env.AddProxy(lC.Proxy())

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
	var exv ansible.ExtraVars

	if ok, match := ust.ContainsFile("docker-compose.yml"); ok {
		exv = ansible.BuildExtraVars(
			fmt.Sprintf("component_path=%s compose_path=%s stack_name=%s", match.Component().RootPath(), match.RelativePath(), st.Name),
			*stackEf.Input,
			*stackEf.Output,
			buffer)
	} else {
		FailsOnCode(&sc, fmt.Errorf("missing docker-compose.yml"), fmt.Sprintf("Stack %s neither have an install.yml playbook nor a docker-compose.yml file", st.Name), nil)
	}

	// Looking for the availability of a docker-compose.yml
	referencers := make([]model.ComponentReferencer, 0, 0)
	referencers = append(referencers, st)
	for _, p := range rm.Parents {
		referencers = append(referencers, p)
	}

	mPaths := cm.ContainsFile("deploy_compose.yaml", rC.data, referencers...)
	if mPaths.Count() == 0 {
		FailsOnCode(&sc, fmt.Errorf("missing deploy_compose.yaml"), fmt.Sprintf("aaaaaa Unable to located the file %s in order to deploy %s", "deploy_compose.yaml", st.Name), nil)
		sCs.Add(sc)
		return
	}

	match := mPaths.Paths[0]
	for _, i := range mPaths.Paths {
		defer i.Component().Release()
	}

	composeLocated := fmt.Sprintf("compose_path=%s stack_name=%s", filepath.Join(match.Component().RootPath(), match.RelativePath()), st.Name)
	lC.Log().Printf("docker-compose.yml has been located  %s", composeLocated)

	code, err := lC.Ekara().AnsibleManager().Execute(match.Component(), "deploy_compose.yaml", exv, env, rC.data)
	if err != nil {
		pfd := playBookFailureDetail{
			Playbook:  "deploy_compose.yaml",
			Component: match.Component().Name(),
			Code:      code,
		}
		FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
		sCs.Add(sc)
		return
	}

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

	sCs.Add(sc)
}

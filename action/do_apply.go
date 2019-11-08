package action

import (
	"encoding/json"
	"fmt"
	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/engine/util"
)

const (
	setupPlaybook   = "setup.yaml"
	createPlaybook  = "create.yaml"
	installPlaybook = "install.yaml"
	deployPlaybook  = "deploy.yaml"
)

type (
	//ApplyResult contains the results of environment application
	ApplyResult struct {
		Success   bool
		Inventory ansible.Inventory
	}
)

func (r ApplyResult) IsSuccess() bool {
	return r.Success
}

func (r *ApplyResult) FromJson(s string) error {
	err := json.Unmarshal([]byte(s), r)
	if err != nil {
		return err
	}
	return nil
}

func (r ApplyResult) AsJson() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

var (
	applyAction = Action{
		ApplyActionID,
		CheckActionID,
		"Apply",
		[]step{providerSetup, providerCreate, orchestratorSetup, orchestratorInstall, stackDeploy, ansibleInventory},
	}
)

func providerSetup(rC *runtimeContext) (StepResults, Result) {
	sCs := InitStepResults()
	for _, p := range rC.environment.Providers {
		sc := InitPlaybookStepResult("Running the setup phase", p, NoCleanUpRequired)

		// Notify setup progress
		rC.lC.Feedback().ProgressG("apply.provider.setup", len(rC.environment.Providers), "Preparing provider '%s'", p.Name)

		// Provider setup exchange folder
		setupProviderEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, "setup_provider_"+p.Name, &sc)
		if ko {
			sCs.Add(sc)
			return *sCs, nil
		}

		setupProviderEfIn := setupProviderEf.Input
		setupProviderEfOut := setupProviderEf.Output

		// Create a new buffer
		buffer := ansible.CreateBuffer()

		// Prepare parameters
		bp := buildBaseParam(rC, "")
		bp.AddNamedMap("params", p.Parameters)
		if ko := saveBaseParams(bp, setupProviderEfIn, &sc); ko {
			sCs.Add(sc)
			return *sCs, nil
		}

		// Prepare extra vars
		exv := ansible.BuildExtraVars("", setupProviderEfIn, setupProviderEfOut, buffer)

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.AddDefaultOsVars()
		env.AddProxy(rC.lC.Proxy())

		// Adding the environment variables from the provider
		for envK, envV := range p.EnvVars {
			env.Add(envK, envV)
		}

		// We launch the playbook
		usable, err := rC.cF.Use(p, *rC.tplC)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable provider", nil)
		}
		defer usable.Release()

		code, err := rC.aM.Play(usable, *rC.tplC, setupPlaybook, exv, env, rC.lC.Feedback())
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  setupPlaybook,
				Component: p.ComponentName(),
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
			sCs.Add(sc)
			return *sCs, nil
		}
		sCs.Add(sc)
	}

	// Notify setup finish
	rC.lC.Feedback().Progress("apply.provider.setup", "All providers prepared")

	return *sCs, nil
}

func providerCreate(rC *runtimeContext) (StepResults, Result) {
	sCs := InitStepResults()

	if rC.lC.Skipping() > 0 {
		// Notify creation skipping
		rC.lC.Feedback().Progress("apply.provider.create", "Nodeset creation skipped by user request")
		return *sCs, nil
	}

	for _, n := range rC.environment.NodeSets {
		sc := InitPlaybookStepResult("Running the create phase", n, NoCleanUpRequired)

		// Resolve provider
		p, err := n.Provider.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the provider"), nil)
			sCs.Add(sc)
			return *sCs, nil
		}

		// Notify creation progress
		rC.lC.Feedback().ProgressG("apply.provider.create", len(rC.environment.NodeSets), "Creating node set '%s' with provider '%s'", n.Name, p.Name)

		// Create a new buffer
		buffer := ansible.CreateBuffer()

		// Prepare parameters
		bp := buildBaseParam(rC, n.Name)
		bp.AddInt("instances", n.Instances)
		bp.AddInterface("labels", n.Labels)
		bp.AddNamedMap("params", p.Parameters)
		bp.AddInterface("proxy", p.Proxy)

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.AddDefaultOsVars()
		env.AddProxy(rC.lC.Proxy())

		// Adding the environment variables from the nodeset provider
		for envK, envV := range p.EnvVars {
			env.Add(envK, envV)
		}

		// Process hook : environment - provision - before
		runHookBefore(
			rC,
			sCs,
			rC.environment.Hooks.Provision,
			hookContext{"create", n, "environment", "provision", bp, env, buffer},
			NoCleanUpRequired,
		)

		// Process hook : nodeset - provision - before
		runHookBefore(
			rC,
			sCs,
			n.Hooks.Provision,
			hookContext{"create", n, "nodeset", "provision", bp, env, buffer},
			NoCleanUpRequired,
		)

		// Node creation exchange folder
		nodeCreateEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, "create_"+n.Name, &sc)
		if ko {
			sCs.Add(sc)
			return *sCs, nil
		}

		if ko := saveBaseParams(bp, nodeCreateEf.Input, &sc); ko {
			sCs.Add(sc)
			return *sCs, nil
		}

		// Prepare extra vars
		exv := ansible.BuildExtraVars("", nodeCreateEf.Input, nodeCreateEf.Output, buffer)

		// Make the component usable
		usable, err := rC.cF.Use(p, *rC.tplC)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable provider", nil)
		}
		defer usable.Release()

		// Launch the playbook
		code, err := rC.aM.Play(usable, *rC.tplC, createPlaybook, exv, env, rC.lC.Feedback())

		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  createPlaybook,
				Component: p.ComponentName(),
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
			sCs.Add(sc)
			return *sCs, nil
		}
		sCs.Add(sc)

		// Process hook : nodeset - provision - after
		runHookAfter(
			rC,
			sCs,
			n.Hooks.Provision,
			hookContext{"create", n, "nodeset", "provision", bp, env, buffer},
			NoCleanUpRequired,
		)

		// Process hook : environment - provision - after
		runHookAfter(
			rC,
			sCs,
			rC.environment.Hooks.Provision,
			hookContext{"create", n, "environment", "provision", bp, env, buffer},
			NoCleanUpRequired,
		)
	}

	// Notify creation finish
	rC.lC.Feedback().Progress("apply.provider.create", "All node sets created")

	return *sCs, nil
}

func orchestratorSetup(rC *runtimeContext) (StepResults, Result) {
	o := rC.environment.Orchestrator
	sCs := InitStepResults()
	sc := InitPlaybookStepResult("Running the orchestrator setup phase", o, NoCleanUpRequired)

	// Notify setup progress
	rC.lC.Feedback().ProgressG("apply.orchestrator.setup", 1, "Preparing orchestrator")

	// Create a new buffer
	buffer := ansible.CreateBuffer()

	// Orchestrator setup exchange folder
	setupOrchestratorEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, "setup_orchestrator", &sc)
	if ko {
		sCs.Add(sc)
		return *sCs, nil
	}

	// Prepare parameters
	bp := buildBaseParam(rC, "")
	bp.AddNamedMap("params", o.Parameters)
	if ko := saveBaseParams(bp, setupOrchestratorEf.Input, &sc); ko {
		sCs.Add(sc)
		return *sCs, nil
	}

	// Prepare environment variables
	env := ansible.BuildEnvVars()
	env.AddDefaultOsVars()
	env.AddProxy(rC.lC.Proxy())

	// Adding the environment variables from the nodeset orchestrator
	for envK, envV := range o.EnvVars {
		env.Add(envK, envV)
	}

	// Prepare extra vars
	exv := ansible.BuildExtraVars("", setupOrchestratorEf.Input, setupOrchestratorEf.Output, buffer)

	// Make the component usable
	usable, err := rC.cF.Use(o, *rC.tplC)
	if err != nil {
		FailsOnCode(&sc, err, "An error occurred getting the usable orchestrator", nil)
	}
	defer usable.Release()

	// We launch the playbook
	code, err := rC.aM.Play(usable, *rC.tplC, setupPlaybook, exv, env, rC.lC.Feedback())
	if err != nil {
		pfd := playBookFailureDetail{
			Playbook:  setupPlaybook,
			Component: o.ComponentName(),
			Code:      code,
		}
		FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
		sCs.Add(sc)
		return *sCs, nil
	}

	// Notify setup progress
	rC.lC.Feedback().Progress("apply.orchestrator.setup", "Orchestrator prepared")

	sCs.Add(sc)
	return *sCs, nil
}

func orchestratorInstall(rC *runtimeContext) (StepResults, Result) {
	sCs := InitStepResults()

	if rC.lC.Skipping() > 1 {
		// Notify installation skipping
		rC.lC.Feedback().Progress("apply.orchestrator.install", "Orchestrator installation skipped by user request")
		return *sCs, nil
	}

	for _, n := range rC.environment.NodeSets {
		sc := InitPlaybookStepResult("Running the orchestrator installation phase", n, NoCleanUpRequired)

		// Resolve the provider
		p, err := n.Provider.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the provider reference"), nil)
			sCs.Add(sc)
			return *sCs, nil
		}

		// Resolve the orchestrator
		o, err := n.Orchestrator.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the orchestrator reference"), nil)
			sCs.Add(sc)
			return *sCs, nil
		}

		// Notify setup progress
		rC.lC.Feedback().ProgressG("apply.orchestrator.install", len(rC.environment.NodeSets), "Installing orchestrator on node set '%s'", n.Name)

		// Create a new buffer
		buffer := ansible.CreateBuffer()

		// Orchestrator install exchange folder
		installOrchestratorEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, "install_orchestrator_"+n.Name, &sc)
		if ko {
			sCs.Add(sc)
			return *sCs, nil
		}

		// Prepare parameters
		bp := buildBaseParam(rC, n.Name)
		bp.AddInterface("labels", n.Labels)
		bp.AddNamedMap("params", o.Parameters)
		bp.AddInterface("proxy", p.Proxy)
		if ko := saveBaseParams(bp, installOrchestratorEf.Input, &sc); ko {
			sCs.Add(sc)
			return *sCs, nil
		}

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.AddDefaultOsVars()
		env.AddProxy(rC.lC.Proxy())

		// Adding the environment variables from the nodeset orchestrator
		for envK, envV := range o.EnvVars {
			env.Add(envK, envV)
		}

		// Prepare extra vars
		exv := ansible.BuildExtraVars("", installOrchestratorEf.Input, installOrchestratorEf.Output, buffer)

		// Make the component usable
		usable, err := rC.cF.Use(o, *rC.tplC)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable orchestrator", nil)
		}
		defer usable.Release()

		// Launch the playbook
		code, err := rC.aM.Play(usable, *rC.tplC, installPlaybook, exv, env, rC.lC.Feedback())
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  installPlaybook,
				Component: o.ComponentName(),
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
			sCs.Add(sc)
			return *sCs, nil
		}
		sCs.Add(sc)
	}

	// Notify setup finish
	rC.lC.Feedback().ProgressG("apply.orchestrator.install", len(rC.environment.NodeSets), "Orchestrator installed on all node sets")

	return *sCs, nil
}

func stackDeploy(rC *runtimeContext) (StepResults, Result) {
	sCs := InitStepResults()

	if rC.lC.Skipping() > 2 {
		// Notify installation skipping
		rC.lC.Feedback().Progress("apply.stack.deploy", "Stack deployment installation skipped by user request")
		return *sCs, nil
	}

	for _, st := range rC.environment.Stacks {
		sc := InitPlaybookStepResult("Deploying stack", st, NoCleanUpRequired)
		sCs.Add(sc)

		// Notify stack deploy
		rC.lC.Feedback().ProgressG("apply.stack.deploy", len(rC.environment.Stacks), "Deploying stack '%s'", st.Name)

		// Stack deploy exchange folder for the given provider
		fName := fmt.Sprintf("deploy_stack_%s", st.Name)

		stackEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, fName, &sc)
		if ko {
			sCs.Add(sc)
			return *sCs, nil
		}

		// Create a new buffer
		buffer := ansible.CreateBuffer()

		// Prepare parameters
		bp := buildBaseParam(rC, "")
		bp.AddNamedMap("params", st.Parameters)
		if ko := saveBaseParams(bp, stackEf.Input, &sc); ko {
			sCs.Add(sc)
			return *sCs, nil
		}

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.AddDefaultOsVars()
		env.AddProxy(rC.lC.Proxy())

		// Adding the environment variables from the stack
		for envK, envV := range st.EnvVars {
			env.Add(envK, envV)
		}

		// Process hook : environment - deploy - before
		runHookBefore(
			rC,
			sCs,
			rC.environment.Hooks.Deploy,
			hookContext{"deploy", st, "environment", "deploy", bp, env, buffer},
			NoCleanUpRequired,
		)

		// Process hook : stack - deploy - before
		runHookBefore(
			rC,
			sCs,
			st.Hooks.Deploy,
			hookContext{"deploy", st, "stack", "deploy", bp, env, buffer},
			NoCleanUpRequired,
		)

		// Make the stack usable
		ust, err := rC.cF.Use(st, *rC.tplC)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable stack", nil)
		}
		defer ust.Release()

		// If the stack is not self deployable, use the orchestrator deploy playbook
		var target component.UsableComponent
		var deployExtraVars string
		if ok, _ := ust.ContainsFile(deployPlaybook); !ok {
			o, err := rC.cF.Use(rC.environment.Orchestrator, *rC.tplC)
			if err != nil {
				FailsOnCode(&sc, err, "An error occurred getting the usable orchestrator", nil)
			}
			defer o.Release()
			target = o
			deployExtraVars = fmt.Sprintf("stack_path=%s stack_name=%s", ust.RootPath(), st.Name)
		} else {
			target = ust
		}

		// Prepare the extra vars
		exv := ansible.BuildExtraVars(
			deployExtraVars,
			stackEf.Input,
			stackEf.Output,
			buffer)

		// Execute the playbook
		code, err := rC.aM.Play(target, *rC.tplC, deployPlaybook, exv, env, rC.lC.Feedback())
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  deployPlaybook,
				Component: target.Name(),
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
			sCs.Add(sc)
			return *sCs, nil
		}

		// Process hook : stack - deploy - after
		runHookAfter(
			rC,
			sCs,
			st.Hooks.Deploy,
			hookContext{"deploy", st, "stack", "deploy", bp, env, buffer},
			NoCleanUpRequired,
		)

		// Process hook : environment - deploy - after
		runHookAfter(
			rC,
			sCs,
			rC.environment.Hooks.Deploy,
			hookContext{"deploy", st, "environment", "deploy", bp, env, buffer},
			NoCleanUpRequired,
		)

		sCs.Add(sc)
	}

	// Notify stack deploy finish
	rC.lC.Feedback().Progress("apply.stack.deploy", "All stacks deployed")

	return *sCs, nil
}

func ansibleInventory(rC *runtimeContext) (StepResults, Result) {
	sCs := InitStepResults()
	sr := InitPlaybookStepResult("Building inventory", nil, NoCleanUpRequired)

	inv, err := rC.aM.Inventory(*rC.tplC)
	if err != nil {
		FailsOnCode(&sr, err, "An error occurred during inventory", nil)
	}

	sCs.Add(sr)
	return *sCs, &ApplyResult{
		Success:   true,
		Inventory: inv,
	}
}

func buildBaseParam(rC *runtimeContext, nodeSetName string) ansible.BaseParam {
	return ansible.BuildBaseParam(rC.environment, rC.lC.SSHPublicKey(), rC.lC.SSHPrivateKey(), nodeSetName)
}

func createChildExchangeFolder(parent util.FolderPath, name string, sr *StepResult) (util.ExchangeFolder, bool) {
	ef, e := parent.AddChildExchangeFolder(name)
	if e != nil {
		err := fmt.Errorf(ErrorAddingExchangeFolder, name, e.Error())
		FailsOnCode(sr, e, err.Error(), nil)
		return ef, true
	}
	e = ef.Create()
	if e != nil {
		err := fmt.Errorf(ErrorCreatingExchangeFolder, name, e.Error())
		FailsOnCode(sr, e, err.Error(), nil)
		return ef, true
	}
	return ef, false
}

func saveBaseParams(bp ansible.BaseParam, dest util.FolderPath, sr *StepResult) bool {
	b, e := bp.Content()
	if e != nil {
		FailsOnCode(sr, e, fmt.Sprintf("An error occurred creating the base parameters"), nil)
		return true
	}
	_, e = util.SaveFile(dest, util.ParamYamlFileName, b)
	if e != nil {
		FailsOnCode(sr, e, fmt.Sprintf("An error occurred saving the parameter file into :%v", dest.Path()), nil)
		return true
	}
	return false
}

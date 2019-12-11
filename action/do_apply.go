package action

import (
	"encoding/json"
	"fmt"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
)

const (
	setupPlaybook   = "setup.yaml"
	createPlaybook  = "create.yaml"
	installPlaybook = "install.yaml"
	deployPlaybook  = "deploy.yaml"
	copyPlaybook    = "copy.yaml"
)

type (
	//ApplyResult contains the results of environment application
	ApplyResult struct {
		Success   bool
		Inventory ansible.Inventory
	}
)

//IsSuccess returns true id the action execution was successful
func (r ApplyResult) IsSuccess() bool {
	return r.Success
}

//FromJson fills an action returned content from a JSON content
func (r *ApplyResult) FromJson(s string) error {
	err := json.Unmarshal([]byte(s), r)
	if err != nil {
		return err
	}
	return nil
}

//AsJson returns the action returned content as JSON
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
		[]step{providerSetup, providerCreate, orchestratorSetup, orchestratorInstall, copy, stackDeploy, ansibleInventory},
	}
)

func providerSetup(rC *runtimeContext) (StepResults, Result) {
	sCs := InitStepResults()
	for _, p := range rC.environment.Providers {
		if !rC.cF.IsAvailable(p) {
			continue
		}

		sc := InitPlaybookStepResult("Running the provider setup phase", p, NoCleanUpRequired)

		// Notify setup progress
		rC.lC.Feedback().ProgressG("provider.setup", len(rC.environment.Providers), "Preparing provider '%s'", p.Name)

		// Provider setup exchange folder
		setupProviderEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, "setup_provider_"+p.Name, &sc)
		if ko {
			sCs.Add(sc)
			return *sCs, nil
		}

		setupProviderEfIn := setupProviderEf.Input
		setupProviderEfOut := setupProviderEf.Output

		// Prepare parameters
		bp := buildBaseParam(rC, "")
		bp.AddNamedMap("params", p.Parameters)
		if ko := saveBaseParams(bp, setupProviderEfIn, &sc); ko {
			sCs.Add(sc)
			return *sCs, nil
		}

		// Prepare extra vars
		exv := ansible.CreateExtraVars(setupProviderEfIn, setupProviderEfOut)

		// We launch the playbook
		usable, err := rC.cF.Use(p, rC.tplC)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable provider", nil)
			sCs.Add(sc)
			return *sCs, nil
		}
		defer usable.Release()

		code, err := rC.aM.Play(usable, rC.tplC, setupPlaybook, exv)
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
	rC.lC.Feedback().Progress("provider.setup", "All providers prepared")

	return *sCs, nil
}

func providerCreate(rC *runtimeContext) (StepResults, Result) {
	sCs := InitStepResults()

	if rC.lC.Skipping() > 0 {
		// Notify creation skipping
		rC.lC.Feedback().Progress("provider.create", "Nodeset creation skipped by user request")
		return *sCs, nil
	}

	for _, n := range rC.environment.NodeSets {
		sc := InitPlaybookStepResult("Running the provider create phase", n, NoCleanUpRequired)

		// Resolve provider
		p, err := n.Provider.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the provider"), nil)
			sCs.Add(sc)
			return *sCs, nil
		}

		// Notify creation progress
		rC.lC.Feedback().ProgressG("provider.create", len(rC.environment.NodeSets), "Creating node set '%s' with provider '%s'", n.Name, p.Name)

		// Prepare parameters
		bp := buildBaseParam(rC, n.Name)
		bp.AddInt("instances", n.Instances)
		bp.AddInterface("labels", n.Labels)
		bp.AddNamedMap("params", p.Parameters)
		bp.AddInterface("proxy", p.Proxy)

		// Process hook : environment - provision - before
		runHookBefore(
			rC,
			sCs,
			rC.environment.Hooks.Provision,
			hookContext{"create", n, "environment", "provision", bp},
			NoCleanUpRequired,
		)

		// Process hook : nodeset - provision - before
		runHookBefore(
			rC,
			sCs,
			n.Hooks.Provision,
			hookContext{"create", n, "nodeset", "provision", bp},
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
		exv := ansible.CreateExtraVars(nodeCreateEf.Input, nodeCreateEf.Output)

		// Make the component usable
		usable, err := rC.cF.Use(p, rC.tplC)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable provider", nil)
			sCs.Add(sc)
			return *sCs, nil
		}
		defer usable.Release()

		// Launch the playbook
		code, err := rC.aM.Play(usable, rC.tplC, createPlaybook, exv)

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

		// Process hook : nodeset - provision - after
		runHookAfter(
			rC,
			sCs,
			n.Hooks.Provision,
			hookContext{"create", n, "nodeset", "provision", bp},
			NoCleanUpRequired,
		)

		// Process hook : environment - provision - after
		runHookAfter(
			rC,
			sCs,
			rC.environment.Hooks.Provision,
			hookContext{"create", n, "environment", "provision", bp},
			NoCleanUpRequired,
		)
		sCs.Add(sc)
	}

	// Notify creation finish
	rC.lC.Feedback().Progress("provider.create", "All node sets created")

	return *sCs, nil
}

func orchestratorSetup(rC *runtimeContext) (StepResults, Result) {
	o := rC.environment.Orchestrator
	sCs := InitStepResults()
	sc := InitPlaybookStepResult("Running the orchestrator setup phase", o, NoCleanUpRequired)

	// Notify setup progress
	rC.lC.Feedback().ProgressG("orchestrator.setup", 1, "Preparing orchestrator")

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

	// Prepare extra vars
	exv := ansible.CreateExtraVars(setupOrchestratorEf.Input, setupOrchestratorEf.Output)

	// Make the component usable
	usable, err := rC.cF.Use(o, rC.tplC)
	if err != nil {
		FailsOnCode(&sc, err, "An error occurred getting the usable orchestrator", nil)
		sCs.Add(sc)
		return *sCs, nil
	}
	defer usable.Release()

	// We launch the playbook
	code, err := rC.aM.Play(usable, rC.tplC, setupPlaybook, exv)
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
	rC.lC.Feedback().Progress("orchestrator.setup", "Orchestrator prepared")

	sCs.Add(sc)
	return *sCs, nil
}

func orchestratorInstall(rC *runtimeContext) (StepResults, Result) {
	o := rC.environment.Orchestrator
	sCs := InitStepResults()

	if rC.lC.Skipping() > 1 {
		// Notify installation skipping
		rC.lC.Feedback().Progress("orchestrator.install", "Orchestrator installation skipped by user request")
		return *sCs, nil
	}

	sc := InitPlaybookStepResult("Running the orchestrator install phase", nil, NoCleanUpRequired)

	// Notify setup progress
	rC.lC.Feedback().Progress("orchestrator.install", "Installing orchestrator")

	// Orchestrator install exchange folder
	installOrchestratorEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, "install_orchestrator", &sc)
	if ko {
		sCs.Add(sc)
		return *sCs, nil
	}

	// Prepare parameters
	bp := buildBaseParam(rC, "")
	bp.AddNamedMap("params", o.Parameters)
	if ko := saveBaseParams(bp, installOrchestratorEf.Input, &sc); ko {
		sCs.Add(sc)
		return *sCs, nil
	}

	// Prepare extra vars
	exv := ansible.CreateExtraVars(installOrchestratorEf.Input, installOrchestratorEf.Output)

	// Make the component usable
	usable, err := rC.cF.Use(o, rC.tplC)
	if err != nil {
		FailsOnCode(&sc, err, "An error occurred getting the usable orchestrator", nil)
		sCs.Add(sc)
		return *sCs, nil
	}
	defer usable.Release()

	// Launch the playbook
	code, err := rC.aM.Play(usable, rC.tplC, installPlaybook, exv)
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

	// Notify install finish
	rC.lC.Feedback().Progress("orchestrator.install", "Orchestrator installed")

	return *sCs, nil
}

func copy(rC *runtimeContext) (StepResults, Result) {
	sCs := InitStepResults()

	if rC.lC.Skipping() > 2 {
		// Notify installation skipping
		rC.lC.Feedback().Progress("stack.deploy", "Stack deployment copy skipped by user request")
		return *sCs, nil
	}

	withCopies := make([]model.Stack, 0, 0)
	for _, st := range rC.environment.Stacks {
		if len(st.Copies.Content) > 0 {
			withCopies = append(withCopies, st)
		}
	}

	for _, st := range withCopies {
		// Notify stack copy
		rC.lC.Feedback().ProgressG("stack.copy", len(withCopies), "Copying stack '%s'", st.Name)
		sc := InitPlaybookStepResult("Copying stack", st, NoCleanUpRequired)

		// Make the stack usable
		ust, err := rC.cF.Use(st, rC.tplC)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable stack", nil)
		}
		defer ust.Release()

		for path, cp := range st.Copies.Content {

			// Stack copy exchange folder for the given stack
			fName := fmt.Sprintf("copy_stack_%s", st.Name)

			stackEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, fName, &sc)
			if ko {
				sCs.Add(sc)
				return *sCs, nil
			}

			// Prepare the extra vars
			exv := ansible.CreateExtraVars(stackEf.Input, stackEf.Output)

			// If the stack is not self copyable, use the orchestrator copy playbook
			var target component.UsableComponent
			if ok, _ := ust.ContainsFile(copyPlaybook); !ok {
				o, err := rC.cF.Use(rC.environment.Orchestrator, rC.tplC)
				if err != nil {
					FailsOnCode(&sc, err, "An error occurred getting the usable orchestrator", nil)
				}
				defer o.Release()
				target = o
			} else {
				target = ust
			}

			exv.Add("stack_path", ust.RootPath())
			if cp.Once {
				exv.Add("copy_once", "true")
			} else {
				exv.Add("copy_once", "false")
			}

			exv.Add("copy_path", path)
			exv.AddArray("copy_sources", cp.Sources.Content)
			exv.AddMap("copy_labels", cp.Labels)

			// Execute the playbook
			code, err := rC.aM.Play(target, rC.tplC, copyPlaybook, exv)
			if err != nil {
				pfd := playBookFailureDetail{
					Playbook:  copyPlaybook,
					Component: target.Name(),
					Code:      code,
				}
				FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
				sCs.Add(sc)
				return *sCs, nil
			}
		}

		sCs.Add(sc)

		// Notify stack copy finish
		rC.lC.Feedback().Progress("stack.copy", "The stack has been copied")
	}
	return *sCs, nil
}

func stackDeploy(rC *runtimeContext) (StepResults, Result) {
	sCs := InitStepResults()

	if rC.lC.Skipping() > 2 {
		// Notify installation skipping
		rC.lC.Feedback().Progress("stack.deploy", "Stack deployment installation skipped by user request")
		return *sCs, nil
	}

	for _, st := range rC.environment.Stacks {
		sc := InitPlaybookStepResult("Deploying stack", st, NoCleanUpRequired)

		// Notify stack deploy
		rC.lC.Feedback().ProgressG("stack.deploy", len(rC.environment.Stacks), "Deploying stack '%s'", st.Name)

		// Stack deploy exchange folder for the given stack
		fName := fmt.Sprintf("deploy_stack_%s", st.Name)

		stackEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, fName, &sc)
		if ko {
			sCs.Add(sc)
			return *sCs, nil
		}

		// Prepare parameters
		bp := buildBaseParam(rC, "")
		bp.AddNamedMap("params", st.Parameters)
		if ko := saveBaseParams(bp, stackEf.Input, &sc); ko {
			sCs.Add(sc)
			return *sCs, nil
		}

		// Process hook : environment - deploy - before
		runHookBefore(
			rC,
			sCs,
			rC.environment.Hooks.Deploy,
			hookContext{"deploy", st, "environment", "deploy", bp},
			NoCleanUpRequired,
		)

		// Process hook : stack - deploy - before
		runHookBefore(
			rC,
			sCs,
			st.Hooks.Deploy,
			hookContext{"deploy", st, "stack", "deploy", bp},
			NoCleanUpRequired,
		)

		// Make the stack usable
		ust, err := rC.cF.Use(st, rC.tplC)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable stack", nil)
			sCs.Add(sc)
			return *sCs, nil
		}
		defer ust.Release()

		// Prepare the extra vars
		exv := ansible.CreateExtraVars(stackEf.Input, stackEf.Output)

		// If the stack is not self deployable, use the orchestrator deploy playbook
		var target component.UsableComponent
		if ok, _ := ust.ContainsFile(deployPlaybook); st.Playbook == "" && !ok {
			o, err := rC.cF.Use(rC.environment.Orchestrator, rC.tplC)
			if err != nil {
				FailsOnCode(&sc, err, "An error occurred getting the usable orchestrator", nil)
				sCs.Add(sc)
				return *sCs, nil
			}
			defer o.Release()
			target = o
			exv.Add("stack_path", ust.RootPath())
			exv.Add("stack_name", st.Name)
		} else {
			target = ust
		}

		// Execute the playbook
		var effectivePlaybook string
		if st.Playbook != "" {
			effectivePlaybook = st.Playbook
		} else {
			effectivePlaybook = deployPlaybook
		}
		code, err := rC.aM.Play(target, rC.tplC, effectivePlaybook, exv)
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  effectivePlaybook,
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
			hookContext{"deploy", st, "stack", "deploy", bp},
			NoCleanUpRequired,
		)

		// Process hook : environment - deploy - after
		runHookAfter(
			rC,
			sCs,
			rC.environment.Hooks.Deploy,
			hookContext{"deploy", st, "environment", "deploy", bp},
			NoCleanUpRequired,
		)

		sCs.Add(sc)
	}

	// Notify stack deploy finish
	rC.lC.Feedback().Progress("stack.deploy", "All stacks deployed")

	return *sCs, nil
}

func ansibleInventory(rC *runtimeContext) (StepResults, Result) {
	sCs := InitStepResults()
	sc := InitPlaybookStepResult("Building inventory", nil, NoCleanUpRequired)

	inv, err := rC.aM.Inventory(rC.tplC)
	if err != nil {
		FailsOnCode(&sc, err, "An error occurred during inventory", nil)
		sCs.Add(sc)
		return *sCs, nil
	}

	sCs.Add(sc)
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

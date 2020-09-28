package action

import (
	"encoding/json"
	"fmt"
	"github.com/GroupePSA/componentizer"
	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/model"
	"github.com/ekara-platform/engine/util"
)

const (
	setupPlaybook   = "setup.yaml"
	createPlaybook  = "create.yaml"
	installPlaybook = "install.yaml"
	deployPlaybook  = "deploy.yaml"
	copyPlaybook    = "copy.yaml"
	checkPlaybook   = "check.yaml"
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
		[]step{
			providerSetup, //TODO to be moved soon
			createHookBefore,
			providerCreate,
			createHookAfter,

			orchestratorSetup, //TODO to be moved soon
			installHookBefore,
			orchestratorInstall,
			installHookAfter,

			deployHookBefore,
			stackCopy, // TODO merge into one deploy step with sub-function calls (so we can sort stacks only once)
			stackCheck,
			stackDeploy,
			deployHookAfter,

			ansibleInventory,
			//initApi,
		},
	}
)

func providerSetup(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()
	for _, p := range rC.environment.Providers {
		if !rC.cM.IsAvailable(p) {
			continue
		}

		sc := InitPlaybookStepResult("Running the provider setup phase", p, NoCleanUpRequired)

		// Notify setup progress
		rC.lC.Feedback().ProgressG("provider.setup", len(rC.environment.Providers), "Preparing provider '%s'", p.Name)

		// Provider setup exchange folder
		setupProviderEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, "setup_provider_"+p.Name, &sc)
		if ko {
			sCs.Add(sc)
			return *sCs
		}

		setupProviderEfIn := setupProviderEf.Input
		setupProviderEfOut := setupProviderEf.Output

		// Prepare parameters
		bp := buildBaseParam(rC, "")
		bp.AddNamedMap("params", p.Parameters())
		if ko := saveBaseParams(bp, setupProviderEfIn, &sc); ko {
			sCs.Add(sc)
			return *sCs
		}

		// Prepare extra vars
		exv := ansible.CreateExtraVars(setupProviderEfIn, setupProviderEfOut)

		// We make the provider usable
		usable, err := rC.cM.Use(p, rC.tplC)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable provider", nil)
			sCs.Add(sc)
			return *sCs
		}
		defer usable.Release()

		// We launch the playbook
		code, err := rC.aM.Play(usable, rC.tplC, setupPlaybook, exv)
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  setupPlaybook,
				Component: p.ComponentId(),
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
			sCs.Add(sc)
			return *sCs
		}
		sCs.Add(sc)
	}

	// Notify setup finish
	rC.lC.Feedback().Progress("provider.setup", "All providers prepared")
	return *sCs
}

func createHookBefore(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()

	if len(rC.environment.Hooks.Create.Before) == 0 {
		return *sCs
	}

	if rC.lC.Skipping() > 0 {
		rC.lC.Feedback().Progress("create.hook.before", "Creation skipped by user request")
		return *sCs
	}

	// Notify creation progress
	rC.lC.Feedback().ProgressG("create.hook.before", 1, "Hook before creating node sets ")

	// Prepare parameters
	bp := buildBaseParam(rC, "")

	// Process hook : environment - create - before
	runHookBefore(
		rC,
		sCs,
		rC.environment.Hooks.Create,
		hookContext{"create", rC.environment, "environment", "create", bp},
		NoCleanUpRequired,
	)

	rC.lC.Feedback().Progress("create.hook.before", "All hooks executed")
	return *sCs
}

func providerCreate(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()

	if rC.lC.Skipping() > 0 {
		rC.lC.Feedback().Progress("provider.create", "Nodeset creation skipped by user request")
		return *sCs
	}

	for _, n := range rC.environment.NodeSets {
		sc := InitPlaybookStepResult("Running the provider create phase", n, NoCleanUpRequired)

		// Resolve provider
		p, err := n.Provider.Resolve(rC.environment)
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the provider"), nil)
			sCs.Add(sc)
			return *sCs
		}

		// Notify creation progress
		rC.lC.Feedback().ProgressG("provider.create", len(rC.environment.NodeSets), "Creating node set '%s' with provider '%s'", n.Name, p.Name)

		// Prepare parameters
		bp := buildBaseParam(rC, n.Name)
		bp.AddInt("instances", n.Instances)
		bp.AddInterface("labels", n.Labels)
		bp.AddNamedMap("params", p.Parameters())
		bp.AddInterface("proxy", p.Proxy())

		// Process hook : nodeset - create - before
		runHookBefore(
			rC,
			sCs,
			n.Hooks.Create,
			hookContext{"create", n, "nodeset", "create", bp},
			NoCleanUpRequired,
		)

		// Node creation exchange folder
		nodeCreateEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, "create_"+n.Name, &sc)
		if ko {
			sCs.Add(sc)
			return *sCs
		}

		if ko := saveBaseParams(bp, nodeCreateEf.Input, &sc); ko {
			sCs.Add(sc)
			return *sCs
		}

		// Prepare extra vars
		exv := ansible.CreateExtraVars(nodeCreateEf.Input, nodeCreateEf.Output)

		// Make the provider usable
		usable, err := rC.cM.Use(p, rC.tplC)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable provider", nil)
			sCs.Add(sc)
			return *sCs
		}
		defer usable.Release()

		// Launch the playbook
		code, err := rC.aM.Play(usable, rC.tplC, createPlaybook, exv)

		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  createPlaybook,
				Component: p.ComponentId(),
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
			sCs.Add(sc)
			return *sCs
		}

		// Process hook : nodeset - create - after
		runHookAfter(
			rC,
			sCs,
			n.Hooks.Create,
			hookContext{"create", n, "nodeset", "create", bp},
			NoCleanUpRequired,
		)
		sCs.Add(sc)
	}

	// Notify creation finish
	rC.lC.Feedback().Progress("provider.create", "All node sets created")
	return *sCs
}

func createHookAfter(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()

	if len(rC.environment.Hooks.Create.After) == 0 {
		return *sCs
	}

	if rC.lC.Skipping() > 0 {
		rC.lC.Feedback().Progress("create.hook.after", "Creation skipped by user request")
		return *sCs
	}

	// Notify creation progress
	rC.lC.Feedback().ProgressG("create.hook.after", 1, "Hook after creating node sets")

	// Prepare parameters
	bp := buildBaseParam(rC, "")

	// Process hook : environment - create - after
	runHookAfter(
		rC,
		sCs,
		rC.environment.Hooks.Create,
		hookContext{"create", rC.environment, "environment", "create", bp},
		NoCleanUpRequired,
	)

	rC.lC.Feedback().Progress("provider.create.hook.after", "All hooks executed")
	return *sCs
}

func orchestratorSetup(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()

	if rC.lC.Skipping() > 1 {
		rC.lC.Feedback().Progress("orchestrator.setup", "Orchestrator installation skipped by user request")
		return *sCs
	}

	o := rC.environment.Orchestrator
	sc := InitPlaybookStepResult("Running the orchestrator setup phase", o, NoCleanUpRequired)

	// Notify setup progress
	rC.lC.Feedback().ProgressG("orchestrator.setup", 1, "Preparing orchestrator")

	// Orchestrator setup exchange folder
	setupOrchestratorEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, "setup_orchestrator", &sc)
	if ko {
		sCs.Add(sc)
		return *sCs
	}

	// Prepare parameters
	bp := buildBaseParam(rC, "")
	bp.AddNamedMap("params", o.Parameters())
	if ko := saveBaseParams(bp, setupOrchestratorEf.Input, &sc); ko {
		sCs.Add(sc)
		return *sCs
	}

	// Prepare extra vars
	exv := ansible.CreateExtraVars(setupOrchestratorEf.Input, setupOrchestratorEf.Output)

	// Make the orchestrator usable
	usable, err := rC.cM.Use(o, rC.tplC)
	if err != nil {
		FailsOnCode(&sc, err, "An error occurred getting the usable orchestrator", nil)
		sCs.Add(sc)
		return *sCs
	}
	defer usable.Release()

	// We launch the playbook
	code, err := rC.aM.Play(usable, rC.tplC, setupPlaybook, exv)
	if err != nil {
		pfd := playBookFailureDetail{
			Playbook:  setupPlaybook,
			Component: o.ComponentId(),
			Code:      code,
		}
		FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
		sCs.Add(sc)
		return *sCs
	}

	// Notify setup progress
	rC.lC.Feedback().Progress("orchestrator.setup", "Orchestrator prepared")

	sCs.Add(sc)
	return *sCs
}

func installHookBefore(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()

	if len(rC.environment.Hooks.Install.Before) == 0 {
		return *sCs
	}

	if rC.lC.Skipping() > 1 {
		rC.lC.Feedback().Progress("install.hook.before", "Installation skipped by user request")
		return *sCs
	}

	// Notify creation progress
	rC.lC.Feedback().ProgressG("install.hook.before", 1, "Hook before installing the orchestrator ")

	// Prepare parameters
	bp := buildBaseParam(rC, "")

	// Process hook : environment - install - before
	runHookBefore(
		rC,
		sCs,
		rC.environment.Hooks.Install,
		hookContext{"install", rC.environment, "environment", "install", bp},
		NoCleanUpRequired,
	)

	rC.lC.Feedback().Progress("install.hook.before", "All hooks executed")
	return *sCs
}

func orchestratorInstall(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()

	if rC.lC.Skipping() > 1 {
		rC.lC.Feedback().Progress("orchestrator.install", "Orchestrator installation skipped by user request")
		return *sCs
	}

	o := rC.environment.Orchestrator
	sc := InitPlaybookStepResult("Running the orchestrator install phase", nil, NoCleanUpRequired)

	// Notify setup progress
	rC.lC.Feedback().Progress("orchestrator.install", "Installing orchestrator")

	// Orchestrator install exchange folder
	installOrchestratorEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, "install_orchestrator", &sc)
	if ko {
		sCs.Add(sc)
		return *sCs
	}

	// Prepare parameters
	bp := buildBaseParam(rC, "")
	bp.AddNamedMap("params", o.Parameters())
	if ko := saveBaseParams(bp, installOrchestratorEf.Input, &sc); ko {
		sCs.Add(sc)
		return *sCs
	}

	// Prepare extra vars
	exv := ansible.CreateExtraVars(installOrchestratorEf.Input, installOrchestratorEf.Output)

	// Make the orchestrator usable
	usable, err := rC.cM.Use(o, rC.tplC)
	if err != nil {
		FailsOnCode(&sc, err, "An error occurred getting the usable orchestrator", nil)
		sCs.Add(sc)
		return *sCs
	}
	defer usable.Release()

	// Launch the playbook
	code, err := rC.aM.Play(usable, rC.tplC, installPlaybook, exv)
	if err != nil {
		pfd := playBookFailureDetail{
			Playbook:  installPlaybook,
			Component: o.ComponentId(),
			Code:      code,
		}
		FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
		sCs.Add(sc)
		return *sCs
	}
	sCs.Add(sc)

	// Notify install finish
	rC.lC.Feedback().Progress("orchestrator.install", "Orchestrator installed")
	return *sCs
}

func installHookAfter(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()

	if len(rC.environment.Hooks.Install.After) == 0 {
		return *sCs
	}

	if rC.lC.Skipping() > 1 {
		rC.lC.Feedback().Progress("install.hook.after", "Installation skipped by user request")
		return *sCs
	}

	// Notify creation progress
	rC.lC.Feedback().ProgressG("install.hook.after", 1, "Hook after installing the orchestrator ")

	// Prepare parameters
	bp := buildBaseParam(rC, "")

	// Process hook : environment - install - after
	runHookAfter(
		rC,
		sCs,
		rC.environment.Hooks.Install,
		hookContext{"install", rC.environment, "environment", "install", bp},
		NoCleanUpRequired,
	)

	rC.lC.Feedback().Progress("install.hook.after", "All hooks executed")
	return *sCs
}

func deployHookBefore(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()

	if len(rC.environment.Hooks.Deploy.Before) == 0 {
		return *sCs
	}

	if rC.lC.Skipping() > 2 {
		rC.lC.Feedback().Progress("deploy.hook.before", "Stack deploy skipped by user request")
		return *sCs
	}

	// Notify creation progress
	rC.lC.Feedback().ProgressG("deploy.hook.before", 1, "Hook before deploying the stacks ")

	// Prepare parameters
	bp := buildBaseParam(rC, "")

	// Process hook : environment - deploy - before
	runHookBefore(
		rC,
		sCs,
		rC.environment.Hooks.Deploy,
		hookContext{"deploy", rC.environment, "environment", "deploy", bp},
		NoCleanUpRequired,
	)

	rC.lC.Feedback().Progress("deploy.hook.before", "All hooks executed")
	return *sCs
}

func stackCopy(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()

	if rC.lC.Skipping() > 2 {
		rC.lC.Feedback().Progress("stack.copy", "Stack copy skipped by user request")
		return *sCs
	}

	stackWithCopies := make([]model.Stack, 0, 0)
	for _, st := range rC.environment.Stacks.Sorted() {
		if len(st.Copies) > 0 {
			stackWithCopies = append(stackWithCopies, st)
		}
	}

	for _, st := range stackWithCopies {
		sc := InitPlaybookStepResult("Copying stack", st, NoCleanUpRequired)

		// Notify stack copy
		rC.lC.Feedback().ProgressG("stack.copy", len(stackWithCopies), "Copying files for stack '%s'", st.Name)

		for cpName, cp := range st.Copies {
			if cp.Path != "" {
				rC.lC.Feedback().Detail("Copying '%s' files to '%s'", cpName, cp.Path)

				// Create exchange folder
				stackEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, fmt.Sprintf("copy_stack_%s_%s", st.Name, cpName), &sc)
				if ko {
					sCs.Add(sc)
					return *sCs
				}

				// Make the stack usable
				ust, err := rC.cM.Use(st, rC.tplC)
				if err != nil {
					FailsOnCode(&sc, err, "An error occurred getting the usable stack", nil)
					sCs.Add(sc)
					return *sCs
				}
				defer ust.Release()

				// If the stack is not self copyable, use the orchestrator copy playbook
				var target componentizer.UsableComponent
				if ok, _ := ust.ContainsFile(copyPlaybook); !ok {
					o, err := rC.cM.Use(rC.environment.Orchestrator, rC.tplC)
					if err != nil {
						FailsOnCode(&sc, err, "An error occurred getting the usable orchestrator", nil)
						sCs.Add(sc)
						return *sCs
					}
					defer o.Release()
					target = o
				} else {
					target = ust
				}

				// Prepare the extra vars
				exv := ansible.CreateExtraVars(stackEf.Input, stackEf.Output)
				exv.Add("stack_path", ust.RootPath())
				exv.Add("copy_path", cp.Path)
				if cp.Once {
					exv.Add("copy_once", "true")
				} else {
					exv.Add("copy_once", "false")
				}
				if len(cp.Sources) > 0 {
					exv.AddArray("copy_sources", cp.Sources)
				}
				if len(cp.Labels) > 0 {
					exv.AddMap("copy_labels", cp.Labels)
				}

				// Execute the playbook
				code, err := rC.aM.Play(target, rC.tplC, copyPlaybook, exv)
				if err != nil {
					pfd := playBookFailureDetail{
						Playbook:  copyPlaybook,
						Component: target.Id(),
						Code:      code,
					}
					FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
					sCs.Add(sc)
					return *sCs
				}
			} else {
				rC.lC.Feedback().Detail("No destination path for '%s' files, skipping copy", cpName)
			}
		}

		sCs.Add(sc)

		// Notify stack copy finish
		rC.lC.Feedback().Progress("stack.copy", "All stack files have been copied")
	}
	return *sCs
}

func stackCheck(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()

	if rC.lC.Skipping() > 2 {
		rC.lC.Feedback().Progress("stack.check", "Stack deployment ( check ) skipped by user request")
		return *sCs
	}

	for _, st := range rC.environment.Stacks.Sorted() {
		sc := InitPlaybookStepResult("Checking stacks", st, NoCleanUpRequired)

		// Notify stack check
		rC.lC.Feedback().ProgressG("stack.check", len(rC.environment.Stacks), "Checking stacks '%s'", st.Name)

		// Make the stack usable
		ust, err := rC.cM.Use(st, rC.tplC)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable stack", nil)
			sCs.Add(sc)
			return *sCs
		}
		defer ust.Release()

		//Verify that the stack contains the check playbook
		if ok, _ := ust.ContainsFile(checkPlaybook); !ok {
			// Notify stack deploy finish
			rC.lC.Feedback().Progress("stack.deploy", "No check playbook available for the stack")
			continue
		}

		// Stack deploy exchange folder for the given stack
		fName := fmt.Sprintf("check_stack_%s", st.Name)

		stackEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, fName, &sc)
		if ko {
			sCs.Add(sc)
			return *sCs
		}

		// Prepare parameters
		bp := buildBaseParam(rC, "")
		bp.AddNamedMap("params", st.Parameters())
		if ko := saveBaseParams(bp, stackEf.Input, &sc); ko {
			sCs.Add(sc)
			return *sCs
		}

		// Prepare the extra vars
		exv := ansible.CreateExtraVars(stackEf.Input, stackEf.Output)

		code, err := rC.aM.Play(ust, rC.tplC, checkPlaybook, exv)
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  checkPlaybook,
				Component: ust.Id(),
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
			sCs.Add(sc)
			return *sCs
		}

		sCs.Add(sc)
	}

	// Notify stack deploy finish
	rC.lC.Feedback().Progress("stack.deploy", "All stacks checked")
	return *sCs
}

func stackDeploy(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()

	if rC.lC.Skipping() > 2 {
		// Notify installation skipping
		rC.lC.Feedback().Progress("stack.deploy", "Stack deployment installation skipped by user request")
		return *sCs
	}

	for _, st := range rC.environment.Stacks.Sorted() {
		sc := InitPlaybookStepResult("Deploying stack", st, NoCleanUpRequired)

		// Notify stack deploy
		rC.lC.Feedback().ProgressG("stack.deploy", len(rC.environment.Stacks), "Deploying stack '%s'", st.Name)

		// Stack deploy exchange folder for the given stack
		fName := fmt.Sprintf("deploy_stack_%s", st.Name)

		stackEf, ko := createChildExchangeFolder(rC.lC.Ef().Input, fName, &sc)
		if ko {
			sCs.Add(sc)
			return *sCs
		}

		// Prepare parameters
		bp := buildBaseParam(rC, "")
		bp.AddNamedMap("params", st.Parameters())
		if ko := saveBaseParams(bp, stackEf.Input, &sc); ko {
			sCs.Add(sc)
			return *sCs
		}

		// Process hook : stack - deploy - before
		runHookBefore(
			rC,
			sCs,
			st.Hooks.Deploy,
			hookContext{"deploy", st, "stack", "deploy", bp},
			NoCleanUpRequired,
		)

		// Make the stack usable
		ust, err := rC.cM.Use(st, rC.tplC)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable stack", nil)
			sCs.Add(sc)
			return *sCs
		}
		defer ust.Release()

		// Prepare the extra vars
		exv := ansible.CreateExtraVars(stackEf.Input, stackEf.Output)

		// If the stack is not self deployable, use the orchestrator deploy playbook
		var target componentizer.UsableComponent
		if ok, _ := ust.ContainsFile(deployPlaybook); !ok {
			o, err := rC.cM.Use(rC.environment.Orchestrator, rC.tplC)
			if err != nil {
				FailsOnCode(&sc, err, "An error occurred getting the usable orchestrator", nil)
				sCs.Add(sc)
				return *sCs
			}
			defer o.Release()
			target = o
			exv.Add("stack_path", ust.RootPath())
			exv.Add("stack_name", st.Name)
		} else {
			target = ust
		}

		// Execute the playbook
		code, err := rC.aM.Play(target, rC.tplC, deployPlaybook, exv)
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  deployPlaybook,
				Component: target.Id(),
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
			sCs.Add(sc)
			return *sCs
		}

		// Process hook : stack - deploy - after
		runHookAfter(
			rC,
			sCs,
			st.Hooks.Deploy,
			hookContext{"deploy", st, "stack", "deploy", bp},
			NoCleanUpRequired,
		)

		sCs.Add(sc)
	}

	// Notify stack deploy finish
	rC.lC.Feedback().Progress("stack.deploy", "All stacks deployed")
	return *sCs
}

func deployHookAfter(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()

	if len(rC.environment.Hooks.Deploy.After) == 0 {
		return *sCs
	}

	if rC.lC.Skipping() > 2 {
		rC.lC.Feedback().Progress("deploy.hook.after", "Stack deploy skipped by user request")
		return *sCs
	}

	// Notify creation progress
	rC.lC.Feedback().ProgressG("deploy.hook.after", 1, "Hook after deploying the stacks ")

	// Prepare parameters
	bp := buildBaseParam(rC, "")

	// Process hook : environment - deploy - after
	runHookAfter(
		rC,
		sCs,
		rC.environment.Hooks.Deploy,
		hookContext{"deploy", rC.environment, "environment", "deploy", bp},
		NoCleanUpRequired,
	)

	rC.lC.Feedback().Progress("deploy.hook.after", "All hooks executed")
	return *sCs
}

func ansibleInventory(rC *RuntimeContext) StepResults {
	sCs := InitStepResults()
	sc := InitPlaybookStepResult("Building inventory", nil, NoCleanUpRequired)

	rC.lC.Feedback().ProgressG("inventory", 1, "Generating inventory")

	inv, err := rC.aM.Inventory(rC.tplC)
	if err != nil {
		FailsOnCode(&sc, err, "An error occurred during inventory", nil)
		sCs.Add(sc)
		return *sCs
	}

	sCs.Add(sc)
	rC.result = ApplyResult{
		Success:   true,
		Inventory: inv,
	}
	rC.lC.Feedback().Progress("inventory", "Inventory generated")
	return *sCs
}

//func initApi(rC *RuntimeContext) StepResults {
//	sCs := InitStepResults()
//	sc := InitPlaybookStepResult("Initializing the Ekara API", nil, NoCleanUpRequired)
//	rC.lC.Feedback().ProgressG("api", 1, "Filling Ekara API ")
//
//	r, ok := rC.result.(ApplyResult)
//	if ok {
//		if !r.Success {
//			FailsOnCode(&sc, fmt.Errorf("The Ekara API cannot be filled because the previous result was not successful"), "", nil)
//			sCs.Add(sc)
//			return *sCs
//		}
//		inv := r.Inventory
//		for k := range inv.Hosts {
//			err := post(k, "desc_name", rC.lC.DescriptorName())
//			if err != nil {
//				FailsOnCode(&sc, err, "Error saving the descriptor name", nil)
//				sCs.Add(sc)
//				return *sCs
//			}
//
//			err = post(k, "desc_url", rC.lC.Location())
//			if err != nil {
//				FailsOnCode(&sc, err, "Error saving the descriptor location", nil)
//				sCs.Add(sc)
//				return *sCs
//			}
//
//			err = post(k, "user", rC.lC.User())
//			if err != nil {
//				FailsOnCode(&sc, err, "Error saving the ekara user", nil)
//				sCs.Add(sc)
//				return *sCs
//			}
//
//			err = post(k, "password", rC.lC.Password())
//			if err != nil {
//				FailsOnCode(&sc, err, "Error saving the ekara password", nil)
//				sCs.Add(sc)
//				return *sCs
//			}
//
//			f, err := util.FileRead(rC.lC.SSHPublicKey())
//			if err != nil {
//				FailsOnCode(&sc, err, "Error reading the ekara public SSH key", nil)
//				sCs.Add(sc)
//				return *sCs
//			}
//
//			err = post(k, "ssh_public", base64.StdEncoding.EncodeToString(f))
//			if err != nil {
//				FailsOnCode(&sc, err, "Error saving the ekara public SSH key", nil)
//				sCs.Add(sc)
//				return *sCs
//			}
//
//			f, err = util.FileRead(rC.lC.SSHPrivateKey())
//			if err != nil {
//				FailsOnCode(&sc, err, "Error reading the ekara private SSH key", nil)
//				sCs.Add(sc)
//				return *sCs
//			}
//
//			err = post(k, "ssh_private", base64.StdEncoding.EncodeToString(f))
//			if err != nil {
//				FailsOnCode(&sc, err, "Error saving the ekara private SSH key", nil)
//				sCs.Add(sc)
//				return *sCs
//			}
//
//			f, err = rC.lC.ExternalVars().ToYAML()
//			if err != nil {
//				FailsOnCode(&sc, err, "Error reading the ekara external vars", nil)
//				sCs.Add(sc)
//				return *sCs
//			}
//
//			err = post(k, "params", base64.StdEncoding.EncodeToString(f))
//			if err != nil {
//				FailsOnCode(&sc, err, "Error saving the ekara external parameters", nil)
//				sCs.Add(sc)
//				return *sCs
//			}
//			//We break because the first host is enough to reach the api
//			break
//		}
//	} else {
//		FailsOnCode(&sc, fmt.Errorf("Wrong result type, ApplyResult was expected"), "", nil)
//		sCs.Add(sc)
//		return *sCs
//	}
//	rC.lC.Feedback().Progress("api", "Ekara API Filled")
//	return *sCs
//}
//
//func post(url string, key string, value interface{}) error {
//	//find a way to put the port outside of the engine
//	u := fmt.Sprintf("http://%s:8090/storage/ekara/", url)
//	reqP := `{"key":"%s", "value":"%v"}`
//	jsonStr := []byte(fmt.Sprintf(reqP, key, value))
//
//	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonStr))
//	req.Header.Set("Content-Type", "application/json")
//
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//	fmt.Println("response Status:", resp.Status)
//	fmt.Println("response Headers:", resp.Header)
//	return nil
//}

func buildBaseParam(rC *RuntimeContext, nodeSetName string) ansible.BaseParam {
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

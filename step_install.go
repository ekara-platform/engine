package engine

import (
	"fmt"

	"github.com/ekara-platform/engine/ansible"
)

var installSteps = []step{fsetuporchestrator, fconsumesetuporchestrator, forchestrator}

func fsetuporchestrator(lC LaunchContext, rC *runtimeContext) StepResults {
	sCs := InitStepResults()
	for _, n := range lC.Ekara().ComponentManager().Environment().NodeSets {
		sc := InitPlaybookStepResult("Running the orchestrator setup phase", n, NoCleanUpRequired)
		lC.Log().Printf(LOG_PROCESSING_NODE, n.Name)

		p, err := n.Provider.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the provider"), nil)
			sCs.Add(sc)
			continue
		}

		// Provider setup exchange folder
		setupProviderEf := lC.Ef().Input.Children["setup_provider_"+p.Name]
		// We check if we have a buffer corresponding to the provider setup
		bufferPro := rC.getBuffer(setupProviderEf.Output)

		// Node exchange folder
		nodeCreationEf := lC.Ef().Input.Children["create_"+n.Name]
		// We check if we have a buffer corresponding to the node output
		buffer := rC.getBuffer(nodeCreationEf.Output)

		// Orchestrator setup exchange folder
		setupOrcherstratorEf, ko := createChildExchangeFolder(lC.Ef().Input, "setup_orchestrator_"+n.Name, &sc, lC.Log())
		if ko {
			sCs.Add(sc)
			continue
		}

		// Prepare parameters
		bp := BuildBaseParam(lC, n.Name, p.Name)
		op, err := n.Orchestrator.OrchestratorParams()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred getting the orchestrator parameters"), nil)
			sCs.Add(sc)
			continue
		}
		bp.AddNamedMap("orchestrator", op)
		bp.AddBuffer(buffer)

		if ko := saveBaseParams(bp, lC, setupOrcherstratorEf.Input, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare components map
		if ko := saveComponentMap(lC, setupOrcherstratorEf.Input, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.Add("http_proxy", lC.HttpProxy())
		env.Add("https_proxy", lC.HttpsProxy())
		env.Add("no_proxy", lC.NoProxy())
		env.AddBuffer(buffer)

		// ugly but .... TODO change this
		env.AddBuffer(bufferPro)

		// Prepare extra vars
		exv := ansible.BuildExtraVars("", *setupOrcherstratorEf.Input, *setupOrcherstratorEf.Output, buffer)

		inventory := ""
		if len(bufferPro.Inventories) > 0 {
			inventory = bufferPro.Inventories["inventory_path"]
		}

		// We launch the playbook
		//TODO move this resolve outside of the for loop
		r, err := lC.Ekara().ComponentManager().Environment().Orchestrator.Component.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the orchestrator"), nil)
			sCs.Add(sc)
			continue
		}
		err, code := lC.Ekara().AnsibleManager().Execute(r, "setup.yml", exv, env, inventory)
		if err != nil {
			r, err := p.Component.Resolve()
			if err != nil {
				FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the provider"), nil)
				sCs.Add(sc)
				continue
			}
			pfd := playBookFailureDetail{
				Playbook:  "setup.yml",
				Compoment: r.Id,
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
			sCs.Add(sc)
			continue
		}
		sCs.Add(sc)
	}

	return *sCs
}

func fconsumesetuporchestrator(lC LaunchContext, rC *runtimeContext) StepResults {
	sCs := InitStepResults()
	for _, n := range lC.Ekara().ComponentManager().Environment().NodeSets {
		sc := InitCodeStepResult("Consuming the orchestrator setup phase", n, NoCleanUpRequired)
		lC.Log().Printf("Consume orchestrator setup for node %s", n.Name)
		setupOrcherstratorEf := lC.Ef().Input.Children["setup_orchestrator_"+n.Name].Output
		err, buffer := ansible.GetBuffer(setupOrcherstratorEf, lC.Log(), "node:"+n.Name)
		// Keep a reference on the buffer based on the output folder
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred getting the buffer"), nil)
			sCs.Add(sc)
			continue
		}
		rC.buffer[setupOrcherstratorEf.Path()] = buffer
		sCs.Add(sc)
	}
	return *sCs
}

func forchestrator(lC LaunchContext, rC *runtimeContext) StepResults {
	sCs := InitStepResults()
	for _, n := range lC.Ekara().ComponentManager().Environment().NodeSets {
		sc := InitPlaybookStepResult("Running the orchestrator installation phase", n, NoCleanUpRequired)
		lC.Log().Printf(LOG_PROCESSING_NODE, n.Name)

		p, err := n.Provider.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the nodeset"), nil)
			sCs.Add(sc)
			continue
		}

		// Provider setup exchange folder
		setupProviderEf := lC.Ef().Input.Children["setup_provider_"+p.Name]
		// We check if we have a buffer corresponding to the provider setup
		bufferPro := rC.getBuffer(setupProviderEf.Output)

		// Orchestrator setup exchange folder
		setupOrcherstratorEf := lC.Ef().Input.Children["setup_orchestrator_"+n.Name]
		// We check if we have a buffer corresponding to the orchestrator setup
		buffer := rC.getBuffer(setupOrcherstratorEf.Output)

		installOrcherstratorEf, ko := createChildExchangeFolder(lC.Ef().Input, "install_orchestrator_"+n.Name, &sc, lC.Log())
		if ko {
			sCs.Add(sc)
			continue
		}

		// Prepare parameters
		bp := BuildBaseParam(lC, n.Name, p.Name)
		bp.AddInterface("labels", n.Labels)
		op, err := n.Orchestrator.OrchestratorParams()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred getting the orchestrator parameters"), nil)
			sCs.Add(sc)
			continue
		}
		bp.AddNamedMap("orchestrator", op)

		// TODO check how to clean all proxies
		pr := lC.Ekara().ComponentManager().Environment().Providers[p.Name].Proxy
		bp.AddInterface("proxy", pr)
		bp.AddBuffer(buffer)

		if ko := saveBaseParams(bp, lC, installOrcherstratorEf.Input, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare components map
		if ko := saveComponentMap(lC, installOrcherstratorEf.Input, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.Add("http_proxy", lC.HttpProxy())
		env.Add("https_proxy", lC.HttpsProxy())
		env.Add("no_proxy", lC.NoProxy())
		env.AddBuffer(buffer)

		// ugly but .... TODO change this
		env.AddBuffer(bufferPro)

		// Prepare extra vars
		exv := ansible.BuildExtraVars("", *installOrcherstratorEf.Input, *installOrcherstratorEf.Output, buffer)

		inventory := ""
		if len(bufferPro.Inventories) > 0 {
			inventory = bufferPro.Inventories["inventory_path"]
		}

		// We launch the playbook
		//TODO move this resolve outside of the for loop
		r, err := lC.Ekara().ComponentManager().Environment().Orchestrator.Component.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the orchestrator"), nil)
			sCs.Add(sc)
			continue
		}
		err, code := lC.Ekara().AnsibleManager().Execute(r, "install.yml", exv, env, inventory)
		if err != nil {
			r, err := p.Component.Resolve()
			if err != nil {
				FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the provider"), nil)
				sCs.Add(sc)
				continue
			}
			pfd := playBookFailureDetail{
				Playbook:  "install.yml",
				Compoment: r.Id,
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
			sCs.Add(sc)
			continue
		}
		sCs.Add(sc)
	}
	return *sCs
}

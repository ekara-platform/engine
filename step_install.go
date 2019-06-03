package engine

import (
	"fmt"

	"github.com/ekara-platform/model"

	"github.com/ekara-platform/engine/ansible"
)

var installSteps = []step{fsetuporchestrator, fconsumesetuporchestrator, forchestrator}

func fsetuporchestrator(lC LaunchContext, rC *runtimeContext) StepResults {
	sCs := InitStepResults()
	for _, n := range lC.Ekara().ComponentManager().Environment().NodeSets {
		sc := InitPlaybookStepResult("Running the orchestrator setup phase", n, NoCleanUpRequired)
		lC.Log().Printf(LogProcessingNode, n.Name)

		// Resolve provider
		p, err := n.Provider.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the provider"), nil)
			sCs.Add(sc)
			continue
		}

		// Resolve orchestrator
		o, err := n.Orchestrator.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred getting the orchestrator parameters"), nil)
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
		setupOrchestratorEf, ko := createChildExchangeFolder(lC.Ef().Input, "setup_orchestrator_"+n.Name, &sc, lC.Log())
		if ko {
			sCs.Add(sc)
			continue
		}

		// Prepare parameters
		bp := BuildBaseParam(lC, n.Name, p.Name) // TODO : review if provider name is ok
		bp.AddNamedMap("orchestrator", buildOrchestratorParams(o))
		bp.AddBuffer(buffer)

		if ko := saveBaseParams(bp, lC, setupOrchestratorEf.Input, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare components map
		if ko := saveComponentMap(lC, setupOrchestratorEf.Input, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.Add("http_proxy", lC.HTTPProxy())
		env.Add("https_proxy", lC.HTTPSProxy())
		env.Add("no_proxy", lC.NoProxy())
		env.AddBuffer(buffer)

		// ugly but .... TODO change this
		env.AddBuffer(bufferPro)

		// Prepare extra vars
		exv := ansible.BuildExtraVars("", *setupOrchestratorEf.Input, *setupOrchestratorEf.Output, buffer)

		inventory := ""
		if len(bufferPro.Inventories) > 0 {
			inventory = bufferPro.Inventories["inventory_path"]
		}

		// We launch the playbook
		code, err := lC.Ekara().AnsibleManager().Execute(o, "setup.yml", exv, env, inventory)
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  "setup.yml",
				Component: o.ComponentName(),
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
		buffer, err := ansible.GetBuffer(setupOrcherstratorEf, lC.Log(), "node:"+n.Name)
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
		lC.Log().Printf(LogProcessingNode, n.Name)

		// Resolve the provider
		p, err := n.Provider.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the provider reference"), nil)
			sCs.Add(sc)
			continue
		}

		// Resolve the orchestrator
		o, err := n.Orchestrator.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the orchestrator reference"), nil)
			sCs.Add(sc)
			continue
		}

		// Provider setup exchange folder
		setupProviderEf := lC.Ef().Input.Children["setup_provider_"+p.Name]
		// We check if we have a buffer corresponding to the provider setup
		bufferPro := rC.getBuffer(setupProviderEf.Output)

		// Orchestrator setup exchange folder
		setupOrchestratorEf := lC.Ef().Input.Children["setup_orchestrator_"+n.Name]
		// We check if we have a buffer corresponding to the orchestrator setup
		buffer := rC.getBuffer(setupOrchestratorEf.Output)

		// Orchestrator install exchange folder
		installOrchestratorEf, ko := createChildExchangeFolder(lC.Ef().Input, "install_orchestrator_"+n.Name, &sc, lC.Log())
		if ko {
			sCs.Add(sc)
			continue
		}

		// Prepare parameters
		bp := BuildBaseParam(lC, n.Name, p.Name) // TODO : review if provider name is ok
		bp.AddInterface("labels", n.Labels)
		bp.AddNamedMap("orchestrator", buildOrchestratorParams(o))

		// TODO check how to clean all proxies
		pr := lC.Ekara().ComponentManager().Environment().Providers[p.Name].Proxy
		bp.AddInterface("proxy", pr)
		bp.AddBuffer(buffer)

		if ko := saveBaseParams(bp, lC, installOrchestratorEf.Input, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare components map
		if ko := saveComponentMap(lC, installOrchestratorEf.Input, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.Add("http_proxy", lC.HTTPProxy())
		env.Add("https_proxy", lC.HTTPSProxy())
		env.Add("no_proxy", lC.NoProxy())
		env.AddBuffer(buffer)

		// ugly but .... TODO change this
		env.AddBuffer(bufferPro)

		// Prepare extra vars
		exv := ansible.BuildExtraVars("", *installOrchestratorEf.Input, *installOrchestratorEf.Output, buffer)

		inventory := ""
		if len(bufferPro.Inventories) > 0 {
			inventory = bufferPro.Inventories["inventory_path"]
		}

		// We launch the playbook
		code, err := lC.Ekara().AnsibleManager().Execute(o, "install.yml", exv, env, inventory)
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  "install.yml",
				Component: o.ComponentName(),
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

func buildOrchestratorParams(o model.Orchestrator) map[string]interface{} {
	op := make(map[string]interface{})
	op["docker"] = o.Docker
	op["params"] = o.Parameters
	return op
}

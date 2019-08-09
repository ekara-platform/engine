package engine

import (
	"fmt"

	"github.com/ekara-platform/engine/ansible"
)

var installSteps = []step{fsetuporchestrator /*fconsumesetuporchestrator,*/, forchestrator}

func fsetuporchestrator(lC LaunchContext, rC *runtimeContext) StepResults {
	cm := lC.Ekara().ComponentManager()
	o := cm.Environment().Orchestrator
	sCs := InitStepResults()
	sc := InitPlaybookStepResult("Running the orchestrator setup phase", o, NoCleanUpRequired)

	// Create a new buffer
	buffer := ansible.CreateBuffer()

	// Orchestrator setup exchange folder
	setupOrchestratorEf, ko := createChildExchangeFolder(lC.Ef().Input, "setup_orchestrator", &sc, lC.Log())
	if ko {
		sCs.Add(sc)
		return *sCs
	}

	// Prepare parameters
	bp := BuildBaseParam(lC, "")
	bp.AddNamedMap("docker", o.Docker)
	bp.AddNamedMap("params", o.Parameters)
	if ko := saveBaseParams(bp, lC, setupOrchestratorEf.Input, &sc); ko {
		sCs.Add(sc)
		return *sCs
	}

	// Prepare environment variables
	env := ansible.BuildEnvVars()
	env.AddDefaultOsVars()
	env.AddProxy(lC.Proxy())

	// Adding the environment variables from the nodeset orchestrator
	for envK, envV := range o.EnvVars {
		env.Add(envK, envV)
	}

	// Prepare extra vars
	exv := ansible.BuildExtraVars("", *setupOrchestratorEf.Input, *setupOrchestratorEf.Output, buffer)

	// We launch the playbook
	usable, err := cm.Use(o)
	if err != nil {
		FailsOnCode(&sc, err, "An error occurred getting the usable orchestrator", nil)
	}
	defer usable.Release()
	code, err := lC.Ekara().AnsibleManager().Execute(usable, "setup.yml", exv, env)
	if err != nil {
		pfd := playBookFailureDetail{
			Playbook:  "setup.yml",
			Component: o.ComponentName(),
			Code:      code,
		}
		FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
		sCs.Add(sc)
		return *sCs
	}

	sCs.Add(sc)
	return *sCs
}

// func fconsumesetuporchestrator(lC LaunchContext, rC *runtimeContext) StepResults {
// 	sCs := InitStepResults()
// 	for _, n := range lC.Ekara().ComponentManager().Environment().NodeSets {
// 		sc := InitCodeStepResult("Consuming the orchestrator setup phase", n, NoCleanUpRequired)
// 		lC.Log().Printf("Consume orchestrator setup for node %s", n.Name)
// 		setupOrcherstratorEf := lC.Ef().Input.Children["setup_orchestrator_"+n.Name].Output
// 		buffer, err := ansible.GetBuffer(setupOrcherstratorEf, lC.Log(), "node:"+n.Name)
// 		// Keep a reference on the buffer based on the output folder
// 		if err != nil {
// 			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred getting the buffer"), nil)
// 			sCs.Add(sc)
// 			continue
// 		}
// 		rC.buffer[setupOrcherstratorEf.Path()] = buffer
// 		sCs.Add(sc)
// 	}
// 	return *sCs
// }

func forchestrator(lC LaunchContext, rC *runtimeContext) StepResults {
	sCs := InitStepResults()
	cm := lC.Ekara().ComponentManager()
	for _, n := range cm.Environment().NodeSets {
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

		// Create a new buffer
		buffer := ansible.CreateBuffer()

		// Orchestrator install exchange folder
		installOrchestratorEf, ko := createChildExchangeFolder(lC.Ef().Input, "install_orchestrator_"+n.Name, &sc, lC.Log())
		if ko {
			sCs.Add(sc)
			continue
		}

		// Prepare parameters
		bp := BuildBaseParam(lC, n.Name)
		bp.AddInterface("labels", n.Labels)
		bp.AddNamedMap("docker", o.Docker)
		bp.AddNamedMap("params", o.Parameters)
		bp.AddInterface("proxy", p.Proxy)
		if ko := saveBaseParams(bp, lC, installOrchestratorEf.Input, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.AddDefaultOsVars()
		env.AddProxy(lC.Proxy())

		// Adding the environment variables from the nodeset orchestrator
		for envK, envV := range o.EnvVars {
			env.Add(envK, envV)
		}

		// Prepare extra vars
		exv := ansible.BuildExtraVars("", *installOrchestratorEf.Input, *installOrchestratorEf.Output, buffer)

		// We launch the playbook
		usable, err := cm.Use(o)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable orchestrator", nil)
		}
		defer usable.Release()
		code, err := lC.Ekara().AnsibleManager().Execute(usable, "install.yml", exv, env)
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

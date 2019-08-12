package engine

import (
	"fmt"

	"github.com/ekara-platform/engine/ansible"
)

var createSteps = []step{fsetup /*fconsumesetup,*/, fcreate /*, fconsumecreate*/}

func fsetup(lC LaunchContext, rC *runtimeContext) StepResults {
	cm := lC.Ekara().ComponentManager()
	sCs := InitStepResults()
	for _, p := range cm.Environment().Providers {
		sc := InitPlaybookStepResult("Running the setup phase", p, NoCleanUpRequired)
		lC.Log().Printf(LogRunningSetupFor, p.Name)

		// TEST FAILURE FOR THE --limit addition
		//hf, _ := c.report.hasFailure()

		// Provider setup exchange folder
		setupProviderEf, ko := createChildExchangeFolder(lC.Ef().Input, "setup_provider_"+p.Name, &sc, lC.Log())
		if ko {
			sCs.Add(sc)
			continue
		}

		setupProviderEfIn := setupProviderEf.Input
		setupProviderEfOut := setupProviderEf.Output

		// Create a new buffer
		buffer := ansible.CreateBuffer()

		// Prepare parameters
		bp := BuildBaseParam(lC, "")
		bp.AddNamedMap("params", p.Parameters)
		if ko := saveBaseParams(bp, lC, setupProviderEfIn, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare extra vars
		exv := ansible.BuildExtraVars("", *setupProviderEfIn, *setupProviderEfOut, buffer)

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.AddDefaultOsVars()
		env.AddProxy(lC.Proxy())

		// Adding the environment variables from the provider
		for envK, envV := range p.EnvVars {
			env.Add(envK, envV)
		}

		// We launch the playbook
		usable, err := cm.Use(p)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable provider", nil)
		}
		defer usable.Release()
		code, err := lC.Ekara().AnsibleManager().Execute(usable, "setup.yml", exv, env)
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  "setup.yml",
				Component: p.ComponentName(),
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

// func fconsumesetup(lC LaunchContext, rC *runtimeContext) StepResults {
// 	sCs := InitStepResults()
// 	for _, p := range lC.Ekara().ComponentManager().Environment().Providers {
// 		sc := InitCodeStepResult("Consuming the setup phase", p, NoCleanUpRequired)
// 		lC.Log().Printf("Consume setup for provider %s", p.Name)
// 		setupProviderEfOut := lC.Ef().Input.Children["setup_provider_"+p.Name].Output
// 		buffer, err := ansible.GetBuffer(setupProviderEfOut, lC.Log(), "provider:"+p.Name)
// 		if err != nil {
// 			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred getting the buffer"), nil)
// 			sCs.Add(sc)
// 			continue
// 		}
// 		// Keep a reference on the buffer based on the output folder
// 		rC.buffer[setupProviderEfOut.Path()] = buffer
// 		sCs.Add(sc)
// 	}
// 	return *sCs
// }

func fcreate(lC LaunchContext, rC *runtimeContext) StepResults {
	cm := lC.Ekara().ComponentManager()
	sCs := InitStepResults()
	for _, n := range cm.Environment().NodeSets {
		sc := InitPlaybookStepResult("Running the create phase", n, NoCleanUpRequired)
		lC.Log().Printf(LogProcessingNode, n.Name)

		// Resolve provider
		p, err := n.Provider.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred resolving the provider"), nil)
			sCs.Add(sc)
			continue
		}

		// Create a new buffer
		buffer := ansible.CreateBuffer()

		// Prepare parameters
		bp := BuildBaseParam(lC, n.Name)
		bp.AddInt("instances", n.Instances)
		bp.AddInterface("labels", n.Labels)
		bp.AddNamedMap("params", p.Parameters)
		bp.AddInterface("proxy", p.Proxy)

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.AddDefaultOsVars()
		env.AddProxy(lC.Proxy())

		// Adding the environment variables from the nodeset provider
		for envK, envV := range p.EnvVars {
			env.Add(envK, envV)
		}

		// Process hook : environment - provision - before
		RunHookBefore(cm,
			lC,
			rC,
			sCs,
			cm.Environment().Hooks.Provision,
			hookContext{"create", n, "environment", "provision", bp, env, buffer},
			NoCleanUpRequired,
		)

		// Process hook : nodeset - provision - before
		RunHookBefore(cm,
			lC,
			rC,
			sCs,
			n.Hooks.Provision,
			hookContext{"create", n, "nodeset", "provision", bp, env, buffer},
			NoCleanUpRequired,
		)

		// Node creation exchange folder
		nodeCreateEf, ko := createChildExchangeFolder(lC.Ef().Input, "create_"+n.Name, &sc, lC.Log())
		if ko {
			sCs.Add(sc)
			continue
		}

		if ko := saveBaseParams(bp, lC, nodeCreateEf.Input, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare extra vars
		exv := ansible.BuildExtraVars("", *nodeCreateEf.Input, *nodeCreateEf.Output, buffer)

		// We launch the playbook
		usable, err := cm.Use(p)
		if err != nil {
			FailsOnCode(&sc, err, "An error occurred getting the usable provider", nil)
		}
		defer usable.Release()
		code, err := lC.Ekara().AnsibleManager().Execute(usable, "create.yml", exv, env)
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  "create.yml",
				Component: p.ComponentName(),
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occurred executing the playbook", pfd)
			sCs.Add(sc)
			continue
		}
		sCs.Add(sc)

		// Process hook : nodeset - provision - after
		RunHookAfter(cm,
			lC,
			rC,
			sCs,
			n.Hooks.Provision,
			hookContext{"create", n, "nodeset", "provision", bp, env, buffer},
			NoCleanUpRequired,
		)

		// Process hook : environment - provision - after
		RunHookAfter(cm,
			lC,
			rC,
			sCs,
			cm.Environment().Hooks.Provision,
			hookContext{"create", n, "environment", "provision", bp, env, buffer},
			NoCleanUpRequired,
		)
	}
	return *sCs
}

// func fconsumecreate(lC LaunchContext, rC *runtimeContext) StepResults {
// 	sCs := InitStepResults()
// 	for _, n := range lC.Ekara().ComponentManager().Environment().NodeSets {
// 		sc := InitCodeStepResult("Consuming the create phase", n, NoCleanUpRequired)
// 		lC.Log().Printf("Consume create for node %s", n.Name)
// 		nodeCreateEf := lC.Ef().Input.Children["create_"+n.Name].Output
// 		buffer, err := ansible.GetBuffer(nodeCreateEf, lC.Log(), "node:"+n.Name)
// 		// Keep a reference on the buffer based on the output folder
// 		if err != nil {
// 			FailsOnCode(&sc, err, fmt.Sprintf("An error occurred getting the buffer"), nil)
// 			sCs.Add(sc)
// 			continue
// 		}
// 		rC.buffer[nodeCreateEf.Path()] = buffer
// 		sCs.Add(sc)
// 	}
// 	return *sCs
// }

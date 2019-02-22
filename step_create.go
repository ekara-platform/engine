package engine

import (
	"fmt"

	"github.com/ekara-platform/engine/ansible"
)

var createSteps = []step{fsetup, fconsumesetup, fcreate, fconsumecreate}

func fsetup(lC LaunchContext, rC *runtimeContext) StepResults {
	sCs := InitStepResults()
	for _, p := range lC.Ekara().ComponentManager().Environment().Providers {
		sc := InitPlaybookStepResult("Running the setup phase", p, NoCleanUpRequired)
		lC.Log().Printf(LOG_RUNNING_SETUP_FOR, p.Name)

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

		// Prepare parameters
		bp := BuildBaseParam(lC, "", p.Name)
		bp.AddNamedMap("params", p.Parameters)

		if ko := saveBaseParams(bp, lC, setupProviderEfIn, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// This is the first "real" step of the process so the used buffer is empty
		emptyBuff := ansible.CreateBuffer()

		// Prepare components map
		if ko := saveComponentMap(lC, setupProviderEfIn, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare extra vars
		exv := ansible.BuildExtraVars("", *setupProviderEfIn, *setupProviderEfOut, emptyBuff)

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.Add("http_proxy", lC.HttpProxy())
		env.Add("https_proxy", lC.HttpsProxy())
		env.Add("no_proxy", lC.NoProxy())

		// Adding the environment variables from the provider
		for envK, envV := range p.EnvVars {
			env.Add(envK, envV)
		}

		// We launch the playbook
		r, err := p.Component.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occured resolving the provider"), nil)
			sCs.Add(sc)
			continue
		}
		err, code := lC.Ekara().AnsibleManager().Execute(r, "setup.yml", exv, env, "")

		if err != nil {
			r, err := p.Component.Resolve()
			if err != nil {
				FailsOnCode(&sc, err, fmt.Sprintf("An error occured resolving the provider"), nil)
				sCs.Add(sc)
				continue
			}
			pfd := playBookFailureDetail{
				Playbook:  "setup.yml",
				Compoment: r.Id,
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occured executing the playbook", pfd)
			sCs.Add(sc)
			continue
		}
		sCs.Add(sc)
	}
	return *sCs
}

func fconsumesetup(lC LaunchContext, rC *runtimeContext) StepResults {
	sCs := InitStepResults()
	for _, p := range lC.Ekara().ComponentManager().Environment().Providers {
		sc := InitCodeStepResult("Consuming the setup phase", p, NoCleanUpRequired)
		lC.Log().Printf("Consume setup for provider %s", p.Name)
		setupProviderEfOut := lC.Ef().Input.Children["setup_provider_"+p.Name].Output
		err, buffer := ansible.GetBuffer(setupProviderEfOut, lC.Log(), "provider:"+p.Name)
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occured getting the buffer"), nil)
			sCs.Add(sc)
			continue
		}
		// Keep a reference on the buffer based on the output folder
		rC.buffer[setupProviderEfOut.Path()] = buffer
		sCs.Add(sc)
	}
	return *sCs
}

func fcreate(lC LaunchContext, rC *runtimeContext) StepResults {
	sCs := InitStepResults()
	for _, n := range lC.Ekara().ComponentManager().Environment().NodeSets {

		sc := InitPlaybookStepResult("Running the create phase", n, NoCleanUpRequired)
		lC.Log().Printf(LOG_PROCESSING_NODE, n.Name)

		p, err := n.Provider.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occured resolving the provider"), nil)
			sCs.Add(sc)
			continue
		}

		// Provider setup exchange folder
		setupProviderEf := lC.Ef().Input.Children["setup_provider_"+p.Name]
		// We check if we have a buffer corresponding to the provider setup output
		buffer := rC.getBuffer(setupProviderEf.Output)

		// Prepare parameters
		bp := BuildBaseParam(lC, n.Name, p.Name)
		bp.AddInt("instances", n.Instances)
		bp.AddInterface("labels", n.Labels)
		bp.AddNamedMap("params", p.Parameters)
		bp.AddInterface("volumes", n.Volumes.AsArray())
		bp.AddBuffer(buffer)

		// Prepare environment variables
		env := ansible.BuildEnvVars()
		env.Add("http_proxy", lC.HttpProxy())
		env.Add("https_proxy", lC.HttpsProxy())
		env.Add("no_proxy", lC.NoProxy())
		env.AddBuffer(buffer)

		inventory := ""
		if len(buffer.Inventories) > 0 {
			inventory = buffer.Inventories["inventory_path"]
		}

		// Process hook : environment - provision - before
		RunHookBefore(lC,
			rC,
			sCs,
			lC.Ekara().ComponentManager().Environment().Hooks.Provision,
			hookContext{"create", n, "environment", "provision", bp, env, buffer, inventory},
			NoCleanUpRequired,
		)

		// Process hook : nodeset - provision - before
		RunHookBefore(lC,
			rC,
			sCs,
			n.Hooks.Provision,
			hookContext{"create", n, "nodeset", "provision", bp, env, buffer, inventory},
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

		// Prepare components map
		if ko := saveComponentMap(lC, nodeCreateEf.Input, &sc); ko {
			sCs.Add(sc)
			continue
		}

		// Prepare extra vars
		exv := ansible.BuildExtraVars("", *nodeCreateEf.Input, *nodeCreateEf.Output, buffer)

		// We launch the playbook
		r, err := p.Component.Resolve()
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occured resolving the provider"), nil)
			sCs.Add(sc)
			continue
		}
		err, code := lC.Ekara().AnsibleManager().Execute(r, "create.yml", exv, env, inventory)
		if err != nil {
			pfd := playBookFailureDetail{
				Playbook:  "create.yml",
				Compoment: r.Id,
				Code:      code,
			}
			FailsOnPlaybook(&sc, err, "An error occured executing the playbook", pfd)
			sCs.Add(sc)
			continue
		}
		sCs.Add(sc)

		// Process hook : nodeset - provision - after
		RunHookAfter(lC,
			rC,
			sCs,
			n.Hooks.Provision,
			hookContext{"create", n, "nodeset", "provision", bp, env, buffer, inventory},
			NoCleanUpRequired,
		)

		// Process hook : environment - provision - after
		RunHookAfter(lC,
			rC,
			sCs,
			lC.Ekara().ComponentManager().Environment().Hooks.Provision,
			hookContext{"create", n, "environment", "provision", bp, env, buffer, inventory},
			NoCleanUpRequired,
		)
	}
	return *sCs
}

func fconsumecreate(lC LaunchContext, rC *runtimeContext) StepResults {
	sCs := InitStepResults()
	for _, n := range lC.Ekara().ComponentManager().Environment().NodeSets {
		sc := InitCodeStepResult("Consuming the create phase", n, NoCleanUpRequired)
		lC.Log().Printf("Consume create for node %s", n.Name)
		nodeCreateEf := lC.Ef().Input.Children["create_"+n.Name].Output
		err, buffer := ansible.GetBuffer(nodeCreateEf, lC.Log(), "node:"+n.Name)
		// Keep a reference on the buffer based on the output folder
		if err != nil {
			FailsOnCode(&sc, err, fmt.Sprintf("An error occured getting the buffer"), nil)
			sCs.Add(sc)
			continue
		}
		rC.buffer[nodeCreateEf.Path()] = buffer
		sCs.Add(sc)
	}
	return *sCs
}

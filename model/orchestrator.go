package model

import (
	"github.com/GroupePSA/componentizer"
)

type (
	//Orchestrator specifies the orchestrator used to manage the environment
	Orchestrator struct {
		// The component containing the orchestrator
		cRef componentRef
		// The orchestrator parameters
		params Parameters
		// The orchestrator environment variables
		envVars EnvVars
	}
)

func createOrchestrator(yamlEnv yamlEnvironment) Orchestrator {
	yamlO := yamlEnv.Orchestrator
	o := Orchestrator{
		cRef:    componentRef{ref: yamlO.Component},
		params:  CreateParameters(yamlO.Params),
		envVars: CreateEnvVars(yamlO.Env),
	}
	return o
}

func (o Orchestrator) EnvVars() EnvVars {
	return o.envVars
}

func (o Orchestrator) Parameters() Parameters {
	return o.params
}

func (o *Orchestrator) merge(with Orchestrator) {
	if with.cRef.ref != "" {
		o.cRef = with.cRef
	}
	o.params = o.params.Override(with.params)
	o.envVars = o.envVars.Override(with.envVars)
}

func (o Orchestrator) ComponentId() string {
	return o.cRef.ComponentId()
}

func (o Orchestrator) Component(model interface{}) (componentizer.Component, error) {
	return o.cRef.Component(model)
}

func (o Orchestrator) DescType() string {
	return "Orchestrator"
}

func (o Orchestrator) DescName() string {
	return o.cRef.ref
}

func (o Orchestrator) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	return validate(e, loc, o.cRef)
}

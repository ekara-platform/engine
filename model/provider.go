package model

import (
	"errors"
	"github.com/GroupePSA/componentizer"
)

type (
	// Provider contains the whole specification of a cloud provider where to
	// create an environment
	Provider struct {
		// The Name of the provider
		Name string
		// The component containing the provider
		cRef componentRef
		// The provider parameters
		params Parameters
		// The provider environment variables
		envVars EnvVars
		// The provider proxy
		proxy Proxy
	}

	//Providers lists all the providers required to build the environemt
	Providers map[string]Provider
)

// createProviders creates all the providers declared into the provided environment
func createProviders(yamlEnv yamlEnvironment) Providers {
	res := Providers{}
	for name, yamlProvider := range yamlEnv.Providers {
		res[name] = Provider{
			Name:    name,
			cRef:    componentRef{ref: yamlProvider.Component},
			params:  CreateParameters(yamlProvider.Params),
			envVars: CreateEnvVars(yamlProvider.Env),
			proxy:   createProxy(yamlProvider.Proxy),
		}
	}
	return res
}

func (p Provider) DescType() string {
	return "Provider"
}

func (p Provider) DescName() string {
	return p.Name
}

func (p Provider) EnvVars() EnvVars {
	return p.envVars
}

func (p Provider) Parameters() Parameters {
	return p.params
}

func (p Provider) Proxy() Proxy {
	return p.proxy
}

func (p Provider) ComponentId() string {
	return p.cRef.ComponentId()
}

func (p Provider) Component(model interface{}) (componentizer.Component, error) {
	return p.cRef.Component(model)
}

func (p *Provider) merge(with Provider) {
	if with.cRef.ref != "" {
		p.cRef = with.cRef
	}
	p.params = p.params.Override(with.params)
	p.envVars = p.envVars.Override(with.envVars)
	p.proxy = p.proxy.override(with.proxy)
}

func (r *Providers) merge(with Providers) {
	for id, p := range with {
		if provider, ok := (*r)[id]; ok {
			provider.merge(p)
			(*r)[id] = provider
		} else {
			(*r)[id] = p
		}
	}
}

func (r Providers) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	vErrs := ValidationErrors{}
	if len(r) == 0 {
		vErrs.addError(errors.New("no provider specified"), loc)
	}
	return vErrs
}

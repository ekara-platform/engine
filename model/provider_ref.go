package model

import (
	"errors"
	"fmt"
	"github.com/GroupePSA/componentizer"
)

type (
	// ProviderRef represents a reference to a provider
	ProviderRef struct {
		ref     string
		params  Parameters
		envVars EnvVars
		proxy   Proxy
	}
)

func createProviderRef(yamlRef yamlProviderRef) ProviderRef {
	return ProviderRef{
		ref:     yamlRef.Name,
		params:  CreateParameters(yamlRef.Params),
		proxy:   createProxy(yamlRef.Proxy),
		envVars: CreateEnvVars(yamlRef.Env),
	}
}

func (r *ProviderRef) merge(with ProviderRef) {
	if r.ref == "" {
		r.ref = with.ref
	}
	r.params = r.params.Override(with.params)
	r.envVars = r.envVars.Override(with.envVars)
	r.proxy = r.proxy.override(with.proxy)
}

func (r ProviderRef) Resolve(model interface{}) (Provider, error) {
	provider := model.(Environment).Providers[r.ref]
	return Provider{
		Name:    provider.Name,
		cRef:    provider.cRef,
		params:  provider.params.Override(r.params),
		envVars: provider.envVars.Override(r.envVars),
		proxy:   provider.proxy.override(r.proxy),
	}, nil
}

func (r ProviderRef) ComponentId() string {
	return ""
}

func (r ProviderRef) Component(model interface{}) (componentizer.Component, error) {
	if provider, ok := model.(Environment).Providers[r.ref]; ok {
		return provider.cRef.Component(model)
	} else {
		return component{}, fmt.Errorf("unable to resolve provider %s", r.ref)
	}
}

func (r ProviderRef) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	vErrs := ValidationErrors{}
	if len(r.ref) == 0 {
		vErrs.addError(errors.New("empty provider reference"), loc)
	} else if _, ok := e.Providers[r.ref]; !ok {
		vErrs.addError(errors.New("no such provider: "+r.ref), loc)
	}
	return vErrs
}

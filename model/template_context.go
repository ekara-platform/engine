package model

import (
	"bytes"
	"fmt"
	"github.com/GroupePSA/componentizer"
	"text/template"
)

type (
	// templateContext the context passed to all ekara templates
	TemplateContext struct {
		// Vars represents accessible descriptor variables,
		Vars Parameters
		// Model represents the environment meta-model (read-only)
		// Model TEnvironment TODO fixme
		// Component represents information about the current component
		Component struct {
			// Type of the component
			Type string
			// Name of the component
			Name string
			// Parameters of the component
			Params Parameters
			// Proxy info of the component if any
			Proxy Proxy
			// Environment variables of the component
			EnvVars EnvVars
		}
		Runtime Parameters
	}
)

// CreateTemplateContext Returns a template context
func CreateTemplateContext(params Parameters) *TemplateContext {
	return &TemplateContext{
		Vars: params,
		// TODO fixme Model:   CreateTEnvironmentForEnvironment(Environment{}),
		Runtime: make(map[string]interface{}),
	}
}

func (tplC TemplateContext) Clone(ref componentizer.ComponentRef) componentizer.TemplateContext {
	newTplC := TemplateContext{
		Vars: CloneParameters(tplC.Vars),
		// TODO fixme Model:   CreateTEnvironmentForEnvironment(Environment{}),
		Runtime: CloneParameters(tplC.Runtime),
	}
	if o, ok := ref.(Describable); ok {
		newTplC.Component.Type = o.DescType()
		newTplC.Component.Name = o.DescName()
	}
	if o, ok := ref.(Parameterized); ok {
		newTplC.Component.Params = CloneParameters(o.Parameters())
	}
	if o, ok := ref.(ProxyAware); ok {
		newTplC.Component.Proxy = o.Proxy()
	}
	if o, ok := ref.(EnvVarsAware); ok {
		newTplC.Component.EnvVars = o.EnvVars()
	}
	return newTplC
}

func (tplC TemplateContext) Execute(content string) (string, error) {
	t, err := template.New(fmt.Sprintf("%s:%s", tplC.Component.Type, tplC.Component.Name)).Parse(content)
	if err != nil {
		return "", err
	}
	var result bytes.Buffer
	err = t.Execute(&result, tplC)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

func (tplC *TemplateContext) addVars(vars Parameters) {
	tplC.Vars = tplC.Vars.Override(vars)
}

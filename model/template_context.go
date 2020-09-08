package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/GroupePSA/componentizer"
	"gopkg.in/yaml.v2"
	"strings"
	"text/template"
)

type (
	// templateContext the context passed to all ekara templates
	TemplateContext struct {
		// Vars represents accessible descriptor variables,
		Vars Parameters
		// Model is a copy of the environment model
		Model Environment
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
		Vars:    params,
		Runtime: make(map[string]interface{}),
	}
}

func (tplC *TemplateContext) Clone(ref componentizer.ComponentRef) componentizer.TemplateContext {
	newTplC := TemplateContext{
		Vars:    CloneParameters(tplC.Vars),
		Runtime: CloneParameters(tplC.Runtime),
		Model:   tplC.Model,
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
	return &newTplC
}

func (tplC TemplateContext) Execute(content string) (string, error) {
	t, err := template.New(fmt.Sprintf("%s:%s", tplC.Component.Type, tplC.Component.Name)).Funcs(template.FuncMap{
		"yaml":   toYaml,
		"json":   toJson,
		"indent": indent,
	}).Parse(content)
	if err != nil {
		return "", err
	}
	var result bytes.Buffer
	err = t.Execute(&result, tplC)
	if err != nil {
		return "", fmt.Errorf("templating error: %s", err.Error())
	}
	return result.String(), nil
}

func (tplC *TemplateContext) addVars(vars Parameters) {
	tplC.Vars = tplC.Vars.Override(vars)
}

func toJson(v interface{}) string {
	strB, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(strB)
}

func toYaml(v interface{}) string {
	strB, err := yaml.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(strB)
}

func indent(spaceCount int, v string) string {
	builder := strings.Builder{}
	spaces := ""
	for i := 0; i < spaceCount; i++ {
		spaces = spaces + " "
	}
	for _, line := range strings.Split(v, "\n") {
		builder.WriteString(spaces)
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	return builder.String()
}

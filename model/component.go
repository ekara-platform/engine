package model

import (
	"errors"
	"fmt"
	"github.com/GroupePSA/componentizer"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type (
	//Component represents an element composing an ekara environment
	component struct {
		Id         string
		Repository componentizer.Repository
		// Templates defines the content to template for the component
		Templates []string
		// Playbooks define the playbooks paths for the component
		Playbooks map[string]string
	}

	componentRef struct {
		ref string
	}
)

func CreateComponent(id string, repo componentizer.Repository) componentizer.Component {
	return component{
		Id:         id,
		Repository: repo,
		Templates:  make([]string, 0, 0),
		Playbooks:  make(map[string]string, 0),
	}
}

func (c *component) merge(with component) {
	c.Id = with.Id
	c.Repository.Merge(with.Repository)
	c.Templates = union(c.Templates, with.Templates)
	for k, v := range with.Playbooks {
		c.Playbooks[k] = v
	}
}

func (c component) String() string {
	return c.Id
}

func (c component) GetRepository() componentizer.Repository {
	return c.Repository
}

func (c component) ComponentId() string {
	return c.Id
}

func (c component) Component(model interface{}) (componentizer.Component, error) {
	resolved, ok := model.(Environment).Platform.Components[c.Id]
	if !ok {
		return nil, fmt.Errorf("component %s cannot be found", c.Id)
	}
	return resolved, nil
}

func (c component) ParseComponents(path string, tplC componentizer.TemplateContext) (componentizer.Component, []componentizer.Component, error) {
	var parent componentizer.Component
	var comps []componentizer.Component

	descPath := filepath.Join(path, DefaultDescriptorName)
	if _, err := os.Stat(descPath); err != nil {
		return nil, nil, nil
	}

	yamlRefs := yamlRefs{}
	err := parseYaml(descPath, tplC.(*TemplateContext), &yamlRefs)
	if err != nil {
		return nil, nil, err
	}

	// Find parent
	if yamlRefs.Ekara.Parent.Repository != "" {
		parent, err = c.buildComponent(yamlRefs.Ekara.Base, c.Id+ParentComponentSuffix, yamlRefs.Ekara.Parent)
		if err != nil {
			return nil, nil, err
		}
	}

	// Gather other components
	for cName, yComp := range yamlRefs.Ekara.Components {
		comp, err := c.buildComponent(yamlRefs.Ekara.Base, cName, yComp)
		if err != nil {
			return nil, nil, err
		}
		comps = append(comps, comp)
	}

	// Sort components by name to ensure reproducible executions
	sort.Slice(comps, func(i, j int) bool {
		return comps[i].ComponentId() < comps[j].ComponentId()
	})

	return parent, comps, nil
}

func (c component) ParseModel(path string, tplC componentizer.TemplateContext) (componentizer.Model, error) {
	descPath := filepath.Join(path, DefaultDescriptorName)
	if _, err := os.Stat(descPath); err != nil {
		return nil, nil
	}

	yamlEnv := yamlEnvironment{}
	err := parseYaml(filepath.Join(path, DefaultDescriptorName), tplC.(*TemplateContext), &yamlEnv)
	if err != nil {
		return nil, err
	}
	return CreateEnvironment(c, yamlEnv)
}

func (c component) GetTemplates() (bool, []string) {
	return len(c.Templates) > 0, c.Templates
}

func (c component) Descriptor() string {
	return DefaultDescriptorName
}

func (r componentRef) ComponentId() string {
	return r.ref
}

func (r componentRef) Component(model interface{}) (componentizer.Component, error) {
	c, ok := model.(Environment).Platform.Components[r.ref]
	if !ok {
		return nil, fmt.Errorf("unable to resolve component %s", r.ref)
	}
	return c, nil
}

func (c component) buildComponent(base string, id string, yC yamlComponent) (component, error) {
	u, err := url.Parse(yC.Repository)
	if err != nil {
		return component{}, err
	}
	if base != "" {
		if !strings.HasSuffix(base, "/") {
			base = base + "/"
		}
		b, err := url.Parse(base)
		if err != nil {
			return component{}, err
		}
		u = b.ResolveReference(u)
	}
	repository, err := c.Repository.CreateChildRepository(u, yC.Ref, yC.Auth)
	if err != nil {
		return component{}, err
	}
	return CreateComponent(id, repository).(component), nil
}

func (r componentRef) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	vErrs := ValidationErrors{}
	if len(r.ref) == 0 {
		vErrs.addError(errors.New("empty component reference"), loc.appendPath("component"))
	} else if _, ok := e.Platform.Components[r.ref]; !ok {
		vErrs.addError(errors.New("no such component: "+r.ref), loc.appendPath("component"))
	}
	return vErrs
}

func (r componentRef) String() string {
	return fmt.Sprintf("&%s", r.ref)
}

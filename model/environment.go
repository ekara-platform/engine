package model

import (
	"github.com/GroupePSA/componentizer"
)

type (
	//Environment represents an environment build based on a descriptor
	Environment struct {
		// Ekara platform settings
		Platform Platform
		// The environment name
		QName QualifiedName
		// The environment description
		Description string
		// The orchestrator used to manage the environment
		Orchestrator Orchestrator
		// The providers where to create the environment node sets
		Providers Providers
		// The node sets to create
		NodeSets NodeSets
		// The software stacks to install on the created node sets
		Stacks Stacks
		// The tasks which can be ran against the environment
		Tasks Tasks
		// The hooks linked to the environment lifecycle events
		Hooks EnvironmentHooks
		// The location of the environment root
		loc DescriptorLocation
	}
)

//CreateEnvironment creates a new environment based on the provided yaml
func CreateEnvironment(from component, yamlEnv yamlEnvironment) (Environment, error) {
	env := Environment{}

	// Base information
	env.QName = QualifiedName{
		Name:      yamlEnv.Name,
		Qualifier: yamlEnv.Qualifier,
	}
	env.Description = yamlEnv.Description

	// Create platform first
	var err error
	env.Platform, err = createPlatform(from, yamlEnv.Ekara)
	if err != nil {
		return env, err
	}

	// Create all other items
	env.loc = DescriptorLocation{Descriptor: from.Repository.String()}
	env.Providers = createProviders(yamlEnv)
	env.Orchestrator = createOrchestrator(yamlEnv)
	env.Tasks = createTasks(yamlEnv)
	env.NodeSets = createNodeSets(yamlEnv)
	env.Stacks = createStacks(from, yamlEnv)
	env.Hooks = createEnvHooks(yamlEnv)

	return env, nil
}

func (r Environment) Merge(with componentizer.Model) (componentizer.Model, error) {
	env := with.(Environment)

	r.QName.merge(env.QName)
	r.Description = env.Description
	r.Platform.merge(env.Platform)
	r.Orchestrator.merge(env.Orchestrator)
	r.Providers.merge(env.Providers)
	r.NodeSets.merge(env.NodeSets)
	r.Stacks.merge(env.Stacks)
	r.Tasks.merge(env.Tasks)
	r.Hooks.merge(env.Hooks)

	return r, nil
}

func (r Environment) IsReferenced(c1 componentizer.Component) bool {
	// Check if the component is self or a parent
	if c1.ComponentId() == r.Platform.Self.ComponentId() {
		return true
	}
	for _, c2 := range r.Platform.Parents {
		if c1.ComponentId() == c2.ComponentId() {
			return true
		}
	}

	// Check providers through node sets
	for _, nodeSet := range r.NodeSets {
		c2, err := nodeSet.Provider.Component(r)
		if err == nil && c1.ComponentId() == c2.ComponentId() {
			return true
		}
	}

	// Check orchestrator
	c2, err := r.Orchestrator.Component(r)
	if err == nil && c1.ComponentId() == c2.ComponentId() {
		return true
	}

	// Check stacks
	for _, stack := range r.Stacks {
		c2, err := stack.Component(r)
		if err == nil && c1.ComponentId() == c2.ComponentId() {
			return true
		}
	}

	// Check tasks
	for _, task := range r.Tasks {
		c2, err := task.Component(r)
		if err == nil && c1.ComponentId() == c2.ComponentId() {
			return true
		}
	}

	return false
}

func (r Environment) DescType() string {
	return "Environment"
}

func (r Environment) DescName() string {
	return r.QName.String()
}

func (r Environment) Validate() ValidationErrors {
	vErrs := ValidationErrors{}
	vErrs.merge(validate(r, r.loc, r.QName))
	vErrs.merge(validate(r, r.loc.appendPath("orchestrator"), r.Orchestrator))
	vErrs.merge(validate(r, r.loc.appendPath("providers"), r.Providers))
	vErrs.merge(validate(r, r.loc.appendPath("nodes"), r.NodeSets))
	vErrs.merge(validate(r, r.loc.appendPath("stacks"), r.Stacks))
	vErrs.merge(validate(r, r.loc.appendPath("tasks"), r.Tasks))
	vErrs.merge(validate(r, r.loc.appendPath("hooks"), r.Hooks))
	return vErrs
}

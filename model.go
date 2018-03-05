package descriptor

// Reusable types

type Hook struct {
	Before []string
	After  []string
}

type Tagged struct {
	Tags []string
}

type Parameterized struct {
	params map[string]string
}

// Main sections

type Platform struct {
	Version  string
	Registry string
	Proxy struct {
		Http    string
		Https   string
		NoProxy string `yaml:"noProxy"`
	}
}

type Provider struct {
	Tagged
	Parameterized
}

type NodeSet struct {
	Tagged

	Provider struct {
		Parameterized

		Name string
	}

	Instances int

	Hooks struct {
		Provision Hook
		Destroy   Hook
	}
}

type Stack struct {
	Tagged

	Repository string
	Version    string
	DeployOn   []string `yaml:"deployOn"`

	Hooks struct {
		Deploy   Hook
		Undeploy Hook
	}
}

type Task struct {
	Tagged

	Playbook string
	Cron     string
	RunOn    []string `yaml:"runOn"`

	Hooks struct {
		Execute Hook
	}
}

// Descriptor

type Descriptor struct {
	Tagged

	// Global attributes
	Name         string
	Description  string
	Version      string
	BaseLocation string

	// Imports
	Imports []string

	// Platform attributes
	Lagoon Platform

	// Providers
	Providers Provider

	// Node sets
	Nodes map[string]NodeSet

	// Software stacks
	Stacks map[string]Stack

	// Custom tasks
	Tasks map[string]Task

	// Global hooks
	Hooks struct {
		Init      Hook
		Provision Hook
		Deploy    Hook
		Undeploy  Hook
		Destroy   Hook
	}
}

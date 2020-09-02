package model

import "errors"

const (
	//GenericNodeSetName is the name of the generic node set
	//
	//The generic node set is intended to be used for sharing common
	// content, example: parameter, environment variables..., with all
	// others node sets within the whole descriptor.
	GenericNodeSetName = "*"
)

type (
	// NodeSet contains the whole specification of a nodeset to create on a specific
	// cloud provider
	NodeSet struct {
		// The name of the machines
		Name string
		// The number of machines to create
		Instances int
		// The ref to the provider where to create the machines
		Provider ProviderRef
		// The hooks linked to the node set lifecycle events
		Hooks NodeHooks
		// The labels associated with the nodeset
		Labels Labels
	}

	NodeHooks struct {
		//Create specifies the hook tasks to run when a node set is created
		Create Hook
		//Destroy specifies the hook tasks to run when a node set is destroyed
		Destroy Hook
	}

	//NodeSets represents all the node sets of the environment
	NodeSets map[string]NodeSet
)

func (r NodeSet) DescType() string {
	return "NodeSet"
}

func (r NodeSet) DescName() string {
	return r.Name
}

func (r *NodeSet) merge(with NodeSet) {
	r.Instances = with.Instances
	r.Labels = r.Labels.override(with.Labels)
	r.Provider.merge(with.Provider)
	r.Hooks.merge(with.Hooks)
}

func (r *NodeHooks) merge(with NodeHooks) {
	r.Create.merge(with.Create)
	r.Destroy.merge(with.Destroy)
}

func (r NodeHooks) HasTasks() bool {
	return r.Create.HasTasks() || r.Destroy.HasTasks()
}

func createNodeSets(yamlEnv yamlEnvironment) NodeSets {
	var gNs NodeSet
	res := NodeSets{}
	for name, yamlNodeSet := range yamlEnv.Nodes {
		if name == GenericNodeSetName {
			// The generic node set has been located
			gNs = buildNode(name, yamlNodeSet)
		} else {
			// Standard node set
			res[name] = buildNode(name, yamlNodeSet)
		}
	}

	if gNs.Name == GenericNodeSetName {
		// The generic node set will be merged into all others
		// in order to propagate the common stuff.
		for name, ns := range res {
			mNs := gNs // duplicate generic node set
			mNs.merge(ns)
			res[name] = mNs
		}
	}
	return res
}

func buildNode(name string, yN yamlNode) NodeSet {
	return NodeSet{
		Name:      name,
		Instances: yN.Instances,
		Provider:  createProviderRef(yN.Provider),
		Hooks: struct {
			Create  Hook
			Destroy Hook
		}{
			Create:  createHook("create", yN.Hooks.Create),
			Destroy: createHook("destroy", yN.Hooks.Destroy),
		},
		Labels: yN.Labels,
	}
}

func (r *NodeSets) merge(with NodeSets) {
	for id, n := range with {
		if nodeSet, ok := (*r)[id]; ok {
			nodeSet.merge(n)
			(*r)[id] = nodeSet
		} else {
			(*r)[id] = n
		}
	}
}

func (r NodeSet) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	vErrs := ValidationErrors{}
	if r.Instances <= 0 {
		vErrs.addError(errors.New("instances must be a positive number"), loc.appendPath("instances"))
	}
	vErrs.merge(validate(e, loc.appendPath("provider"), r.Provider))
	vErrs.merge(validate(e, loc.appendPath("hooks"), r.Hooks))
	return vErrs
}

func (r NodeHooks) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	return validate(e, loc, r.Create, r.Destroy)
}

func (r NodeSets) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	vErrs := ValidationErrors{}
	if len(r) == 0 {
		vErrs.addError(errors.New("no node specified"), loc)
	}
	return vErrs
}

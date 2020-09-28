package model

import (
	"errors"
	"github.com/GroupePSA/componentizer"
)

type (
	//Stack represent an Stack installable on the built environment
	Stack struct {
		// The component containing the stack
		cRef componentRef
		// The component the stack was created from
		selfRef componentRef
		// The name of the stack
		Name string
		//Dependencies specifies the stacks on which this one depends
		Dependencies []string
		// The stack parameters
		params Parameters
		// The stack environment variables
		envVars EnvVars
		// The stack content to be copied on volumes
		Copies Copies
		// The hooks linked to the stack lifecycle events
		Hooks StackHooks
	}

	StackHooks struct {
		//Deploy specifies the hook tasks to run when a stack is deployed
		Deploy Hook
	}

	//Stacks represent all the stacks of an environment
	Stacks map[string]Stack

	// Copies represents a list of content to be copied
	// The key of the map is the path where the content should be copied
	// The map content is an array of path patterns to locate the content to be copied
	Copies map[string]Copy

	// Copy represents a content to be copied
	Copy struct {
		//Once indicates if the copy should be done only on one node matching the targeted labels
		Once bool
		// Labels identifies the nodesets where to copy
		Labels Labels
		// Path identifies the destination path of the copy
		Path string
		//Sources identifies the content to copy
		Sources []string
	}
)

func (s Stack) EnvVars() EnvVars {
	return s.envVars
}

func (s Stack) Parameters() Parameters {
	return s.params
}

func (s Stack) DescType() string {
	return "Stack"
}

func (s Stack) DescName() string {
	return s.Name
}

func (s Stack) ComponentId() string {
	if s.isSelfComponent() {
		return s.selfRef.ComponentId()
	} else {
		return s.cRef.ComponentId()
	}
}

func (s Stack) Component(model interface{}) (componentizer.Component, error) {
	if s.isSelfComponent() {
		return s.selfRef.Component(model)
	} else {
		return s.cRef.Component(model)
	}
}

func (s Stack) isSelfComponent() bool {
	return s.cRef.ref == "" || s.cRef.ref == "_"
}

func (r *Stacks) merge(with Stacks) {
	for id, s := range with {
		if stack, ok := (*r)[id]; ok {
			stack.merge(s)
			(*r)[id] = stack
		} else {
			(*r)[id] = s
		}
	}
}

func createStacks(from component, yamlEnv yamlEnvironment) Stacks {
	res := Stacks{}
	for name, yamlStack := range yamlEnv.Stacks {
		s := Stack{
			cRef:         componentRef{ref: yamlStack.Component},
			selfRef:      componentRef{ref: from.Id},
			Name:         name,
			Dependencies: yamlStack.Dependencies,
			params:       CreateParameters(yamlStack.Params),
			envVars:      CreateEnvVars(yamlStack.Env),
			Copies:       createCopies(yamlStack.Copies),
		}

		// Build hooks
		s.Hooks.Deploy = createHook("deploy", yamlStack.Hooks.Deploy)

		res[name] = s
	}
	return res
}

func (s *Stack) merge(with Stack) {
	if with.cRef.ref != "" {
		s.cRef = with.cRef
	}
	s.Dependencies = union(s.Dependencies, with.Dependencies)
	s.params = s.params.Override(with.params)
	s.envVars = s.envVars.Override(with.envVars)
	s.Copies = s.Copies.override(with.Copies)
	s.Hooks.merge(with.Hooks)
}

func (s *StackHooks) merge(with StackHooks) {
	s.Deploy.merge(with.Deploy)
}

func (s StackHooks) HasTasks() bool {
	return s.Deploy.HasTasks()
}

func createCopies(copies map[string]yamlCopy) Copies {
	res := make(map[string]Copy)
	for cpName, yCop := range copies {
		theCopy := Copy{
			Once:   yCop.Once,
			Labels: yCop.Labels,
		}
		theCopy.Sources = yCop.Sources
		theCopy.Path = yCop.Path
		res[cpName] = theCopy
	}
	return res
}

func (r Copies) override(parent Copies) Copies {
	dst := make(map[string]Copy)
	for k, v := range r {
		// We copy all the original content
		dst[k] = v
	}
	for k, v := range parent {
		// if the parent content is new then we add it
		if _, ok := dst[k]; !ok {
			dst[k] = v
		} else {
			// if it's not new we will merge the patterns/labels from the original content and the parent
			work := dst[k]
			work.Sources = union(work.Sources, v.Sources)
			work.Labels = work.Labels.override(v.Labels)
			if work.Path == "" {
				// Only override path if none specified
				work.Path = v.Path
			}
			if v.Once == true {
				// only override once if true (meaning if it's true, it's forever true in children)
				work.Once = true
			}
			dst[k] = work
		}
	}
	return dst
}

func (r Stacks) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	vErrs := ValidationErrors{}
	if len(r) == 0 {
		vErrs.addWarning("no stack specified", loc)
	}
	return vErrs
}

func (s Stack) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	vErrs := ValidationErrors{}
	if s.isSelfComponent() {
		vErrs.merge(validate(e, loc, s.selfRef))
	} else {
		vErrs.merge(validate(e, loc, s.cRef))
	}
	vErrs.merge(validate(e, loc.appendPath("copies"), s.Copies))
	vErrs.merge(validate(e, loc.appendPath("hooks"), s.Hooks))
	if len(s.Dependencies) > 0 {
		for _, dep := range s.Dependencies {
			if _, ok := e.Stacks[dep]; !ok {
				vErrs.addError(errors.New("no such stack: "+dep), loc.appendPath("dependencies"))
			}
		}
	}
	return vErrs
}

func (s StackHooks) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	return validate(e, loc, s.Deploy)
}

package engine

import (
	_ "log"
)

type ErrorType int

const (
	Warning ErrorType = 0
	Error   ErrorType = 1
)

type GrammarErrors struct {
	Errors []GrammarError
}

type GrammarError struct {
	Location  string
	Message   string
	ErrorType ErrorType
}

func (t ErrorType) String() string {
	names := [...]string{
		"Warning",
		"Error"}
	if t < Warning || t > Error {
		return "Unknown"
	} else {
		return names[t]
	}
}

func (ge GrammarErrors) HasErrors() bool {
	for _, v := range ge.Errors {
		if v.ErrorType == Error {
			return true
		}
	}
	return false
}

func (ge GrammarErrors) HasWarnings() bool {
	for _, v := range ge.Errors {
		if v.ErrorType == Warning {
			return true
		}
	}
	return false
}

func (desc *environmentDef) validate() (result *GrammarErrors) {
	result = &GrammarErrors{}

	// Providers are required
	if b := checkNotEmpty("providers", desc.providers.values, result); b {
		// Check the providers content

	}
	// Nodes are required
	if b := checkNotEmpty("nodes", desc.nodes.values, result); b {
		// Check the nodes content
		for _, v := range desc.nodes.values {
			checkNodes(v.(nodeSetDef), result)
		}

	}
	// Stacks are required
	if b := checkNotEmpty("stacks", desc.stacks.values, result); b {
		// Check the stacks content
		for _, v := range desc.stacks.values {
			checkStacks(v.(stackDef), result)
		}
	}

	// Tasks are not required
	for _, v := range desc.tasks.values {
		checkTasks(v.(taskDef), result)
	}

	// Check general hooks
	checkHook("hooks.init", desc.Hooks.Init, desc, result)
	checkHook("hooks.provision", desc.Hooks.Provision, desc, result)
	checkHook("hooks.deploy", desc.Hooks.Deploy, desc, result)
	checkHook("hooks.undeploy", desc.Hooks.Undeploy, desc, result)
	checkHook("hooks.destroy", desc.Hooks.Destroy, desc, result)
	return
}

func checkNotEmpty(mapName string, m namedMap, e *GrammarErrors) bool {
	if len(m) == 0 {
		e.Errors = append(e.Errors, GrammarError{Location: mapName,
			Message: "Can't be empty", ErrorType: Error})
		return false
	}
	return true
}

func checkNodes(n nodeSetDef, e *GrammarErrors) bool {
	// Check the number of instances
	if n.Instances <= 0 {
		e.Errors = append(e.Errors, GrammarError{Location: "nodes." + n.name + ".instances",
			Message: "Must be greather than 0", ErrorType: Error})
	}

	// We check the provider names integrity only if the descriptor has providers
	if n.desc.GetProviderDescriptions().Count() > 0 {
		if n.Provider.Name == "" {
			e.Errors = append(e.Errors, GrammarError{Location: "nodes." + n.name + ".provider",
				Message: "The provider is required", ErrorType: Error})
		} else {
			if _, b := n.GetProviderDescription(); !b {
				e.Errors = append(e.Errors, GrammarError{Location: "nodes." + n.name + ".provider.name",
					Message: "The provider is unknown", ErrorType: Error})
			}
		}
	}

	checkHook("nodes."+n.name+".hooks.destroy", n.Hooks.Destroy, n.desc, e)
	checkHook("nodes."+n.name+".hooks.provision", n.Hooks.Provision, n.desc, e)
	return true
}

func checkHook(location string, hook hookDef, desc *environmentDef, e *GrammarErrors) {
	for _, v := range hook.Before {
		checkHookName(location+".before", v, desc, e)
	}
	for _, v := range hook.After {
		checkHookName(location+".after", v, desc, e)
	}
}

func checkHookName(location string, name string, desc *environmentDef, e *GrammarErrors) {
	if desc.GetTaskDescriptions().Count() == 0 {
		e.Errors = append(e.Errors, GrammarError{Location: location + ":" + name,
			Message: "The hook task is missing", ErrorType: Error})
	} else {
		if desc.GetTaskDescriptions().Contains(name) == false {
			e.Errors = append(e.Errors, GrammarError{Location: location + ":" + name,
				Message: "The hook task is unknown", ErrorType: Error})
		}
	}
}

func checkStacks(s stackDef, e *GrammarErrors) bool {
	// Check stack repository
	if s.Repository == "" {
		e.Errors = append(e.Errors, GrammarError{Location: "stacks." + s.name + ".repository",
			Message: "The repository is required", ErrorType: Error})
	}

	// Check stack repository
	if s.Version == "" {
		e.Errors = append(e.Errors, GrammarError{Location: "stacks." + s.name + ".version",
			Message: "The version is required", ErrorType: Error})
	}

	if len(s.DeployOn.Names) == 0 {
		e.Errors = append(e.Errors, GrammarError{Location: "stacks." + s.name + ".deployOn",
			Message: "Is empty then this stack won't be deployed", ErrorType: s.desc.getWarningType()})
	} else {
		// We check deployOn names integrity only if the descriptor has nodes
		if s.desc.GetNodeDescriptions().Count() > 0 {
			for _, v := range s.DeployOn.Names {
				if len(s.desc.GetNodeDescriptions().GetNodesByLabel(v)) == 0 {
					e.Errors = append(e.Errors, GrammarError{Location: "stacks." + s.name + ".deployOn:" + v,
						Message: "The label is unknown", ErrorType: Error})
				}
			}
		}
	}
	return true
}

func checkTasks(t taskDef, e *GrammarErrors) bool {
	// Check task playbook
	if t.Task == "" {
		e.Errors = append(e.Errors, GrammarError{Location: "tasks." + t.name + ".task",
			Message: "The playbook is required", ErrorType: Error})
	}

	if len(t.RunOn.Names) > 0 {
		// We check runOn names integrity only if the descriptor has nodes
		if t.desc.GetNodeDescriptions().Count() > 0 {
			for _, v := range t.RunOn.Names {
				if len(t.desc.GetNodeDescriptions().GetNodesByLabel(v)) == 0 {
					e.Errors = append(e.Errors, GrammarError{Location: "tasks." + t.name + ".runOn:" + v,
						Message: "The label is unknown", ErrorType: Error})
				}
			}
		}
	}

	return true
}

package engine

import (
	"encoding/json"
	"fmt"
)

type ErrorType int

const (
	Info    ErrorType = 0
	Warning ErrorType = 1
	Error   ErrorType = 2
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
		"Info",
		"Warning",
		"Error"}
	if t < Info || t > Error {
		return "Unknown"
	} else {
		return names[t]
	}
}

func (ge GrammarErrors) getError() error {
	b, err := json.Marshal(ge)
	if err != nil {
		return err
	}
	return fmt.Errorf(string(b))
}

func (ge GrammarErrors) hasError() (error, bool) {
	for _, v := range ge.Errors {
		if v.ErrorType == Warning || v.ErrorType == Error {
			return ge.getError(), true
		}
	}
	return nil, false
}

func FromError(e error) (GrammarErrors, error) {
	res := GrammarErrors{}
	err := json.Unmarshal([]byte(e.Error()), &res)
	if err != nil {
		return res, err
	}
	return res, nil
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

	return
}

func checkNotEmpty(mapName string, m namedMap, e *GrammarErrors) bool {
	if len(m) == 0 {
		e.Errors = append(e.Errors, GrammarError{Location: mapName, Message: "Can't be empty", ErrorType: Error})
		return false
	}
	return true
}

func checkNodes(n nodeSetDef, e *GrammarErrors) bool {
	// Check the number of instances
	if n.Instances <= 0 {
		e.Errors = append(e.Errors, GrammarError{Location: "nodes." + n.name + ".instances", Message: "Must be greather than 0", ErrorType: Error})
	}

	// We check the provider names integrity only if the descriptor has providers
	if n.desc.GetProviderDescriptions().Count() > 0 {
		if n.Provider.Name == "" {
			e.Errors = append(e.Errors, GrammarError{Location: "nodes." + n.name + ".provider", Message: "The provider is required", ErrorType: Error})
		} else {
			if _, b := n.GetProviderDescription(); !b {
				e.Errors = append(e.Errors, GrammarError{Location: "nodes." + n.name + ".provider.name", Message: "The provider is unknown", ErrorType: Error})
			}
		}
	}
	return true
}

func checkStacks(s stackDef, e *GrammarErrors) bool {
	// Check stack repository
	if s.Repository == "" {
		e.Errors = append(e.Errors, GrammarError{Location: "stacks." + s.name + ".repository", Message: "The repository is required", ErrorType: Error})
	}

	// Check stack repository
	if s.Version == "" {
		e.Errors = append(e.Errors, GrammarError{Location: "stacks." + s.name + ".version", Message: "The version is required", ErrorType: Error})
	}

	if len(s.DeployOn.Names) == 0 {
		e.Errors = append(e.Errors, GrammarError{Location: "stacks." + s.name + ".deployOn", Message: "Is empty then this stack won't be deployed", ErrorType: Error})
	} else {
		// We check deployOn names integrity only if the descriptor has nodes
		if s.desc.GetNodeDescriptions().Count() > 0 {
			for _, v := range s.DeployOn.Names {
				if len(s.desc.GetNodeDescriptions().GetNodesByLabel(v)) == 0 {
					e.Errors = append(e.Errors, GrammarError{Location: "stacks." + s.name + ".deployOn:" + v, Message: "The label is unknown", ErrorType: Error})
				}
			}
		}
	}
	return true
}

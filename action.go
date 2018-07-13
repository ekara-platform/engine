package engine

import (
	"strconv"
)

// EngineAction represents an action available on an environment
type EngineAction int

const (
	// Creation of an environment
	ActionCreate EngineAction = iota
	// Update of an environment
	ActionUpdate
	// Validation of an environment
	ActionCheck
	// Deletion of an environment
	ActionDelete
)

// String returns the string representation of the action
func (a EngineAction) String() string {
	return strconv.Itoa(int(a))

}

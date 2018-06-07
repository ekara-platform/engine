package engine

import (
	"strconv"
)

type EngineAction int

const (
	ActionCreate EngineAction = iota
	ActionUpdate
	ActionCheck
	ActionDelete
)

func (a EngineAction) String() string {
	return strconv.Itoa(int(a))

}

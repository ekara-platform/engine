package model

import (
	"strings"
)

type (
	// Describable represents a part of the environment descriptor
	// which can describe itself with a type and a name.
	//
	// Describable is implemented by :
	//  Nodeset
	//  Provider
	//  Stack
	//  Task
	//
	// The Describable interface can be used for example in logs or even in
	// execution report files in order to provider a human readable vision of what
	// occurred instead providing technical concepts.
	//
	Describable interface {
		//DescType returns the type of the environment part being described,
		//Nodeset, Provider, Stack... Usually the type is harcoded into the
		//implementation
		DescType() string
		// DescName returns the name of the environment part being described
		DescName() string
	}

	chained struct {
		descTypes []string
		descNames []string
	}
)

func (c chained) DescType() string {
	return strings.Join(c.descTypes, "-")
}

func (c chained) DescName() string {
	return strings.Join(c.descNames, "-")
}

// ChainDescribable merges the types and names of several describable
func ChainDescribable(descs ...Describable) Describable {
	r := chained{}
	for _, v := range descs {
		r.descTypes = append(r.descTypes, v.DescType())
		r.descNames = append(r.descNames, v.DescName())
	}
	return r
}

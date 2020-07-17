package model

import (
	"fmt"
	"reflect"
)

type (
	// DescriptorLocation represents the location of any element within the environment
	// descriptor
	DescriptorLocation struct {
		//Descriptor is the url of the descriptor
		Descriptor string
		//Path is the location into the descriptor
		//
		// For example if the descriptor contains this:
		//  level1:
		//    level2:
		//      level3: ekara-platform/aws-provider
		//
		// The location of "level3" will be "level1.level2.level3"
		//
		Path string
	}
)

func (r DescriptorLocation) equals(o DescriptorLocation) bool {
	return reflect.DeepEqual(r, o)
}

func (r DescriptorLocation) empty() bool {
	return r.Descriptor == "" && r.Path == ""
}

func (r DescriptorLocation) appendPath(suffix string) DescriptorLocation {
	newLoc := DescriptorLocation{Path: r.Path, Descriptor: r.Descriptor}
	if newLoc.Path == "" {
		newLoc.Path = suffix
	} else {
		newLoc.Path = newLoc.Path + "." + suffix
	}
	return newLoc
}
func (r DescriptorLocation) appendIndex(i int) DescriptorLocation {
	return DescriptorLocation{Path: r.Path + fmt.Sprintf("[%d]", i), Descriptor: r.Descriptor}
}

package scm

import (
	"fmt"
	"os"

	"github.com/ekara-platform/model"
)

type (

	//FetchedComponent represents a component fetched locally
	FetchedComponent struct {
		ID            string
		LocalPath     string
		Descriptor    string
		DescriptorUrl model.EkUrl
		LocalUrl      model.EkUrl
	}
)

func (fc FetchedComponent) String() string {
	return fmt.Sprintf("id: %s, localPath: %s, descriptor: %s", fc.ID, fc.LocalPath, fc.Descriptor)
}

//HasDescriptor returns true if the fetched component contains a descriptor
func (fc FetchedComponent) HasDescriptor() bool {
	if _, err := os.Stat(fc.DescriptorUrl.AsFilePath()); err == nil {
		return true
	}
	return false
}

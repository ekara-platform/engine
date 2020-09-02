package model

import (
	"errors"
	"regexp"
)

// QualifiedName The Qualified name of an environment.
// This name can be used to identify, using for example Tags or Labels, all the
// content created relatively to the environment on the infrastructure of the desired
// cloud provider.
type QualifiedName struct {
	Name      string
	Qualifier string
}

// isAValidQualifier is the regular expression used to validate the qualified name
// of an environment
//
// The qualified name is the concatenation if the environment name and its qualifier
// separated by a "_"
//
// The name and the qualifier can contain only alphanumerical characters
var isValidQualifiedName = regexp.MustCompile(`^[a-zA-Z0-9_a-zA-Z0-90-9]+$`).MatchString

//String returns the qualified name as string
func (r QualifiedName) String() string {
	res := r.Name
	if len(r.Qualifier) > 0 {
		res = res + "_" + r.Qualifier
	}
	return res
}

func (r QualifiedName) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	vErrs := ValidationErrors{}
	if len(r.Name) == 0 {
		vErrs.addError(errors.New("empty environment name"), loc.appendPath("name"))
	} else if !isValidQualifiedName(r.String()) {
		vErrs.addError(errors.New("the environment name or the qualifier contains a non alphanumeric character"), loc.appendPath("name|qualifier"))
	}
	return vErrs
}

func (r *QualifiedName) merge(with QualifiedName) {
	if r.Name == "" {
		r.Name = with.Name
	}
	if r.Qualifier == "" {
		r.Qualifier = with.Qualifier
	}
}

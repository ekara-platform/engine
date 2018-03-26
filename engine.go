package engine

import (
	"log"

	"gopkg.in/yaml.v2"
)

//go:generate go run ./generate/generate.go

type holder struct {
	// Global state
	logger *log.Logger

	// EnvironmentDef info
	location string
	env      *environmentDef
}

// Create creates an environment descriptor based on the provider location.
//
// The location can be an URL over http or https or even a file system location.
func Create(logger *log.Logger, location string) (lagoon Lagoon, err error) {
	h := holder{logger: logger, location: location}
	desc, err := parseDescriptor(h, location)
	if err != nil {
		return nil, err
	}
	err = desc.adjustAndValidate()
	if err != nil {
		return nil, err
	}
	h.env = &desc
	return &h, nil
}

// Create creates an environment descriptor based on the provider serialized
// content.
//
// The serialized content is typically a fully resolved descriptor without any import left.
func CreateFromContent(logger *log.Logger, content []byte) (lagoon Lagoon, err error) {
	h := holder{logger: logger, location: ""}
	desc, err := parseContent(h, content)
	if err != nil {
		return nil, err
	}
	h.env = &desc
	err = desc.adjustAndValidate()
	if err != nil {
		return nil, err
	}
	return &h, nil
}

// GetContent serializes the content on the environment descriptor
func (h holder) GetContent() ([]byte, error) {
	b, e := yaml.Marshal(h.env)
	return b, e
}

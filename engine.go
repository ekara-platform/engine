package engine

import (
	"log"

	"gopkg.in/yaml.v2"
)

//go:generate go run ./generator/generate.go

var h holder

type holder struct {
	// Global state
	logger *log.Logger

	// EnvironmentDef info
	location string
	env      *environmentDef
}

func Create(logger *log.Logger, location string) (lagoon Lagoon, err error) {
	h = holder{logger: logger, location: location}
	desc, err := parseDescriptor(h.location)
	if err != nil {
		return nil, err
	}
	h.env = &desc
	return &h, nil
}

//GetContent serializes the content on the environment descriptor
func (h holder) GetContent() ([]byte, error) {
	b, e := yaml.Marshal(h.env)
	return b, e
}

package engine

import (
	"log"
)

type Lagoon interface {
	GetEnvironment() Environment
}

type holder struct {
	// Global state
	logger *log.Logger

	// EnvironmentDef info
	location string
	env      *environmentDef
}

func (h *holder) GetEnvironment() (Environment) {
	return h.env
}

func Create(logger *log.Logger, location string) (lagoon Lagoon, err error) {
	h := holder{logger: logger, location: location}
	desc, err := parseDescriptor(h.location)
	if err != nil {
		return nil, err
	}
	h.env = &desc
	return &h, nil
}

package ansible

import (
	"fmt"
)

const (
	ansible_ok                int = 0
	ansible_error             int = 1
	ansible_failed            int = 2
	ansible_unreachable       int = 3
	ansible_parser            int = 4
	ansible_bar_or_incomplete int = 5
	ansible_user_insterrupted int = 99
	ansible_unexpected_error  int = 250
)

func ReturnedError(code int) error {

	switch code {
	case 0:
		return nil // OK or no hosts matched
	case 1:
		return buildError("Error") // Error
	case 2:
		return buildError("One or more hosts failed") //One or more hosts failed
	case 3:
		return buildError("One or more hosts were unreachable") //One or more hosts were unreachable
	case 4:
		return buildError("Parser error") //
	case 5:
		return buildError("Bad or incomplete options") //
	case 99:
		return buildError("User interrupted execution") //
	case 250:
		return buildError("Unexpected error  ") //
	}
	return nil
}

func buildError(s string) error {
	return fmt.Errorf("Ansible returned the following error \"%s\"", s)
}

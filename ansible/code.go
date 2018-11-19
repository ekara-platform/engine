package ansible

import (
	"fmt"
)

const (
	ansibleOk               int = 0
	ansibleError            int = 1
	ansibleFailed           int = 2
	ansibleUreachable       int = 3
	ansibleParser           int = 4
	ansibleBarOrIncomplete  int = 5
	ansibleUserInsterrupted int = 99
	ansibleUnexpectedError  int = 250
)

func ReturnedError(code int) error {

	switch code {
	case ansibleOk:
		return nil // OK or no hosts matched
	case ansibleError:
		return buildError("Error") // Error
	case ansibleFailed:
		return buildError("One or more hosts failed") //One or more hosts failed
	case ansibleUreachable:
		return buildError("One or more hosts were unreachable") //One or more hosts were unreachable
	case ansibleParser:
		return buildError("Parser error") //
	case ansibleBarOrIncomplete:
		return buildError("Bad or incomplete options") //
	case ansibleUserInsterrupted:
		return buildError("User interrupted execution") //
	case ansibleUnexpectedError:
		return buildError("Unexpected error  ") //
	}
	return nil
}

func buildError(s string) error {
	return fmt.Errorf("Ansible returned the following error \"%s\"", s)
}

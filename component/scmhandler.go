package component

import (
	"github.com/ekara-platform/model"
)

//ScmHandler represents the common definition of all Source Control Managment handler
type ScmHandler interface {
	Matches(repository model.Repository, path string) bool
	Fetch(repository model.Repository, path string, auth model.Parameters) error
	Update(path string, auth model.Parameters) error
	Switch(path string, ref string) error
}

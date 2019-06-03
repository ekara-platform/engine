package engine

import (
	"log"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
)

//EkaraMock simulates Ekara for testing purposes
type EkaraMock struct {
	// TODO This mock should be deleted once the logic content of the installer has been
	// refactored and moved into the engine
	Env model.Environment
}

//Init simulates the corresponding method in Ekara for testing purposes
func (e EkaraMock) Init(c LaunchContext) error {
	return nil
}

//Environment simulates the corresponding method in Ekara for testing purposes
func (e EkaraMock) Environment() model.Environment {
	return e.Env
}

//ComponentManager simulates the corresponding method in Ekara for testing purposes
func (e EkaraMock) ComponentManager() component.ComponentManager {
	return CMMock{env: e.Env}
}

//AnsibleManager simulates the corresponding method in Ekara for testing purposes
func (e EkaraMock) AnsibleManager() ansible.AnsibleManager {
	return AMMock{}
}

//Logger simulates the corresponding method in Ekara for testing purposes
func (e EkaraMock) Logger() *log.Logger {
	return nil
}

//BaseDir simulates the corresponding method in Ekara for testing purposes
func (e EkaraMock) BaseDir() string {
	return ""
}

//AMMock simulates the AnsibleManager for testing purposes
type AMMock struct {
	// TODO This mock should be deleted once the logic content of the installer has been
	// refactored and moved into the engine
}

//Execute simulates the corresponding method in the AnsibleManager for testing purposes
func (m AMMock) Execute(cr model.ComponentReferencer, playbook string, extraVars ansible.ExtraVars, envars ansible.EnvVars, inventories string) (int, error) {
	return 0, nil
}

//Contains simulates the corresponding method in the AnsibleManager for testing purposes
func (m AMMock) Contains(cr model.ComponentReferencer, playbook string) bool {
	return true
}

//CMMock simulates the ComponentManager for testing purposes
type CMMock struct {
	// TODO This mock should be deleted once the logic content of the installer has been
	// refactored and moved into the engine
	env model.Environment
}

//RegisterComponent simulates the corresponding method in the ComponentManager for testing purposes
func (m CMMock) RegisterComponent(c model.Component) {
	return
}

//Environment simulates the corresponding method in the ComponentManager for testing purposes
func (m CMMock) Environment() model.Environment {
	return m.env
}

//MatchingDirectories simulates the corresponding method in the ComponentManager for testing purposes
func (m CMMock) MatchingDirectories(dirName string) []string {
	result := make([]string, 0)
	return result
}

//ComponentPath simulates the corresponding method in the ComponentManager for testing purposes
func (m CMMock) ComponentPath(cID string) string {
	return ""
}

//ComponentsPaths  simulates the corresponding method in the ComponentManager for testing purposes
func (m CMMock) ComponentsPaths() map[string]string {
	return make(map[string]string)
}

//SaveComponentsPaths simulates the corresponding method in the ComponentManager for testing purposes
func (m CMMock) SaveComponentsPaths(log *log.Logger, dest util.FolderPath) error {
	_, err := util.SaveFile(log, dest, util.ComponentPathsFileName, []byte("DUMMY_CONTENT"))
	if err != nil {
		return err
	}
	return nil
}

//Ensure simulates the corresponding method in the ComponentManager for testing purposes
func (m CMMock) Ensure() error {
	return nil
}

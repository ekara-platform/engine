package engine

import (
	"log"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
)

type EkaraMock struct {
	// TODO This mock should be deleted once the logic content of the installer has been
	// refactored and moved into the engine
	Env model.Environment
}

func (d EkaraMock) Init(repo string, ref string, descriptor string) error {
	return nil
}
func (d EkaraMock) Environment() model.Environment {
	return d.Env
}
func (d EkaraMock) ComponentManager() component.ComponentManager {
	return CMMock{env: d.Env}
}

func (d EkaraMock) AnsibleManager() ansible.AnsibleManager {
	return AMMock{}
}

func (d EkaraMock) Logger() *log.Logger {
	return nil
}

func (d EkaraMock) BaseDir() string {
	return ""
}

type AMMock struct {
	// TODO This mock should be deleted once the logic content of the installer has been
	// refactored and moved into the engine
}

func (m AMMock) Execute(component model.Component, playbook string, extraVars ansible.ExtraVars, envars ansible.EnvVars, inventories string) (error, int) {
	return nil, 0
}

func (m AMMock) Contains(component model.Component, playbook string) bool {
	return true
}

type CMMock struct {
	// TODO This mock should be deleted once the logic content of the installer has been
	// refactored and moved into the engine
	env model.Environment
}

func (m CMMock) RegisterComponent(c model.Component) {
	return
}

func (m CMMock) Environment() model.Environment {
	return m.env
}

func (m CMMock) MatchingDirectories(dirName string) []string {
	result := make([]string, 0)
	return result
}

func (m CMMock) ComponentPath(cId string) string {
	return ""
}

func (m CMMock) ComponentsPaths() map[string]string {
	return make(map[string]string)
}

func (m CMMock) SaveComponentsPaths(log *log.Logger, dest util.FolderPath) error {
	_, err := util.SaveFile(log, dest, util.ComponentPathsFileName, []byte("DUMMY_CONTENT"))
	if err != nil {
		return err
	}
	return nil
}

func (m CMMock) Ensure() error {
	return nil
}

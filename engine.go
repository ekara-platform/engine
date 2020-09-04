package engine

import (
    "path/filepath"

    "github.com/GroupePSA/componentizer"
    "github.com/ekara-platform/engine/action"
    "github.com/ekara-platform/engine/ansible"
    "github.com/ekara-platform/engine/model"
    "github.com/ekara-platform/engine/util"
)

// Ekara is the facade used to manage environments.
type Ekara interface {
    Init(repo componentizer.Repository) error
    ActionManager() action.Manager
    Environment() model.Environment
}

type engine struct {
    // Base context
    lC        util.LaunchContext
    directory string

    // Environment and its data
    environment model.Environment
    tplC        componentizer.TemplateContext

    // Subsystems
    componentManager componentizer.ComponentManager
    ansibleManager   ansible.Manager
    actionManager    action.Manager
}

// Create creates an engine environment descriptor based on the provided location.
//
//	Parameters:
//		lC: the launch context
//		workDir: the directory where the engine will do its work
func Create(lC util.LaunchContext, workDir string) Ekara {
    eng := engine{
        lC:        lC,
        directory: filepath.Clean(workDir),
        tplC:      model.CreateTemplateContext(lC.ExternalVars()),
    }
    eng.componentManager = componentizer.CreateComponentManager(eng.lC.Log(), filepath.Join(eng.directory, "components"))
    eng.ansibleManager = ansible.CreateAnsibleManager(eng.lC, eng.componentManager)
    eng.actionManager = action.CreateActionManager(eng.lC, eng.componentManager, eng.ansibleManager)
    return &eng
}

func (eng *engine) Init(repo componentizer.Repository) error {
    m, err := eng.componentManager.Init(model.CreateComponent(model.MainComponentId, repo), eng.tplC)
    if err != nil {
        return err
    }
    eng.environment = m.(model.Environment)
    return nil
}

func (eng *engine) ActionManager() action.Manager {
    return eng.actionManager
}

func (eng *engine) Environment() model.Environment {
    return eng.environment
}

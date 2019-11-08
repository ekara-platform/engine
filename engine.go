package engine

import (
	"path/filepath"

	"github.com/ekara-platform/engine/action"
	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/model"
)

//Ekara is the facade used to process environments.
type Ekara interface {
	Init() error
	AnsibleManager() ansible.Manager
	ActionManager() action.Manager
}

type engine struct {
	lC        util.LaunchContext
	directory string

	environment model.Environment
	tplC        *model.TemplateContext

	// Subsystems
	ansibleManager ansible.Manager
	actionManager  action.Manager
}

// Create creates an environment descriptor based on the provided location.
//
// The location can be an URL over http or https or even a file system location.
//
//	Parameters:
//		lC: the launch context
//		workDir: the directory where the engine will do its work
func Create(lC util.LaunchContext, workDir string) (Ekara, error) {
	absWorkDir, err := filepath.Abs(workDir)
	if err != nil {
		return nil, err
	}

	eng := &engine{
		lC:        lC,
		directory: absWorkDir,
		tplC:      model.CreateTemplateContext(lC.ExternalVars()),
	}
	return eng, nil
}

func (eng *engine) Init() (err error) {
	// Get CWD in case the descriptor is local
	repo, ref := util.RepositoryFlavor(eng.lC.Location())
	wdURL, err := model.GetCurrentDirectoryURL(eng.lC.Log())
	if err != nil {
		return
	}

	// Create main component
	mainRep, err := model.CreateRepository(model.Base{Url: wdURL}, repo, ref, eng.lC.DescriptorName())
	if err != nil {
		return
	}
	u := eng.lC.User()
	if u != "" {
		auth := make(map[string]interface{})
		auth["method"] = "basic"
		auth["user"] = u
		auth["password"] = eng.lC.Password()
		mainRep.Authentication = auth
	}
	mainComponent := model.CreateComponent(model.MainComponentId, mainRep)

	finder, env, err := component.Build(mainComponent, eng.tplC, eng.directory, eng.lC.Log())
	if err != nil {
		return err
	}
	eng.environment = env
	// Once the environment is created we can create the ansible and action
	// manager passing them copy of the environment
	eng.ansibleManager = ansible.CreateAnsibleManager(eng.lC.Log(), eng.lC.Verbosity(), finder)
	eng.actionManager = action.CreateActionManager(eng.lC, *eng.tplC, eng.environment, finder, eng.ansibleManager)
	return
}

func (eng *engine) TemplateContext() *model.TemplateContext {
	return eng.tplC
}

func (eng *engine) AnsibleManager() ansible.Manager {
	return eng.ansibleManager
}

func (eng *engine) ActionManager() action.Manager {
	return eng.actionManager
}

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
	ComponentManager() component.Manager
	AnsibleManager() ansible.Manager
	ActionManager() action.Manager
}

type engine struct {
	lC        util.LaunchContext
	tplC      *model.TemplateContext
	directory string

	// Subsystems
	componentManager component.Manager
	ansibleManager   ansible.Manager
	actionManager    action.Manager
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

	// TODO : pass launch context to managers + let templateContext be in the componentManager only
	eng := &engine{
		lC:        lC,
		directory: absWorkDir,
	}

	eng.componentManager = component.CreateComponentManager(lC.Log(), lC.ExternalVars(), absWorkDir)
	eng.ansibleManager = ansible.CreateAnsibleManager(lC.Log(), eng.componentManager)
	eng.actionManager = action.CreateActionManager(lC, eng.componentManager, eng.ansibleManager)
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

	// Discover components starting from the main one
	err = eng.componentManager.Init(mainComponent)
	if err != nil {
		return
	}

	// Then ensure all effectively used components are fetched
	err = eng.componentManager.Ensure()
	if err != nil {
		return
	}

	return
}

func (eng *engine) ComponentManager() component.Manager {
	return eng.componentManager
}

func (eng *engine) AnsibleManager() ansible.Manager {
	return eng.ansibleManager
}

func (eng *engine) ActionManager() action.Manager {
	return eng.actionManager
}

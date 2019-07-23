package engine

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/component"

	"github.com/ekara-platform/model"
)

//Engine  represents the Ekara engine in charge of dealing with the environment
type Engine interface {
	Init(c LaunchContext) error
	Logger() *log.Logger
	BaseDir() string
	ComponentManager() *component.ComponentManager
	ReferenceManager() *component.ReferenceManager
	AnsibleManager() ansible.AnsibleManager
}

type context struct {
	// Base attributes
	logger    *log.Logger
	directory string

	// Subsystems
	componentManager *component.ComponentManager
	referenceManager *component.ReferenceManager
	ansibleManager   ansible.AnsibleManager
	actionManager    ActionManager
}

// Create creates an environment descriptor based on the provided location.
//
// The location can be an URL over http or https or even a file system location.
//
//	Parameters:
//		logger: the logger
//		baseDir: the directory where the environment will take place among its
//				 inclusions and related components
//		data: the user data for templating the environment descriptor
func Create(logger *log.Logger, workDir string, data *model.TemplateContext) (Engine, error) {
	absWorkDir, err := filepath.Abs(workDir)
	if err != nil {
		return nil, err
	}

	ctx := &context{
		logger:    logger,
		directory: absWorkDir,
	}

	ctx.componentManager = component.CreateComponentManager(ctx.logger, data, absWorkDir)
	ctx.referenceManager = component.CreateReferenceManager(ctx.componentManager)
	ctx.ansibleManager = ansible.CreateAnsibleManager(ctx.logger, *ctx.componentManager)
	ctx.actionManager = CreateActionManager()
	return ctx, nil
}

//repositoryFlavor returns the repository flavor, branchn tag ..., based on the
// presence of '@' into the given url
func repositoryFlavor(url string) (string, string) {

	if strings.Contains(url, "@") {
		s := strings.Split(url, "@")
		return s[0], s[1]
	}
	return url, ""
}

func (ctx *context) Init(c LaunchContext) (err error) {
	repo, ref := repositoryFlavor(c.Location())
	wdURL, err := model.GetCurrentDirectoryURL(ctx.logger)
	if err != nil {
		return
	}

	// Register main component
	mainRep, err := model.CreateRepository(model.Base{Url: wdURL}, repo, ref, c.Name())
	if err != nil {
		return
	}
	u := c.User()
	if u != "" {
		auth := make(map[string]interface{})
		auth["method"] = "basic"
		auth["user"] = u
		auth["password"] = c.Password()
		mainRep.Authentication = auth
	}

	mainComponent := model.CreateComponent(model.MainComponentId, mainRep)
	// Parse upward all the references composing the environment
	err = ctx.referenceManager.Init(mainComponent)
	if err != nil {
		return
	}
	err = ctx.ReferenceManager().Ensure()
	if err != nil {
		return
	}
	return
}

func (ctx *context) Logger() *log.Logger {
	return ctx.logger
}

func (ctx *context) BaseDir() string {
	return ctx.directory
}

func (ctx *context) ComponentManager() *component.ComponentManager {
	return ctx.componentManager
}

func (ctx *context) ReferenceManager() *component.ReferenceManager {
	return ctx.referenceManager
}

func (ctx *context) AnsibleManager() ansible.AnsibleManager {
	return ctx.ansibleManager
}

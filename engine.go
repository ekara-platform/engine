package engine

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/model"
)

//Engine  represents the Ekara engine in charge of dealing with the environment
type Engine interface {
	Init(repo string, ref string, descriptor string) error
	Logger() *log.Logger
	BaseDir() string
	ComponentManager() component.ComponentManager
	AnsibleManager() ansible.AnsibleManager
}

type context struct {
	// Base attributes
	logger    *log.Logger
	directory string

	// Subsystems
	componentManager component.ComponentManager
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
func Create(logger *log.Logger, workDir string, data map[string]interface{}) (Engine, error) {
	absWorkDir, err := filepath.Abs(workDir)
	if err != nil {
		return nil, err
	}

	ctx := context{
		logger:    logger,
		directory: absWorkDir,
	}

	ctx.componentManager = component.CreateComponentManager(ctx.logger, data, absWorkDir)
	ctx.ansibleManager = ansible.CreateAnsibleManager(ctx.logger, ctx.componentManager)
	ctx.actionManager = CreateActionManager()
	return &ctx, nil
}

func (ctx *context) Init(repo string, ref string, descriptor string) (err error) {
	wdURL, err := model.GetCurrentDirectoryURL(ctx.logger)
	if err != nil {
		return
	}

	// Register main component
	mainRep, err := model.CreateRepository(model.Base{Url: wdURL}, repo, ref, descriptor)
	if err != nil {
		return
	}
	mainComponent := model.CreateComponent("__main__", mainRep)
	ctx.componentManager.RegisterComponent(mainComponent)

	// Ensure the main component is present
	err = ctx.componentManager.Ensure()
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

func (ctx *context) ComponentManager() component.ComponentManager {
	return ctx.componentManager
}

func (ctx *context) AnsibleManager() ansible.AnsibleManager {
	return ctx.ansibleManager
}

//CheckProxy returns the proxy setting from environment variables
//
// See:
//		github.com/ekara-platform/engine/util.HttpProxyEnvVariableKey
//		github.com/ekara-platform/engine/util.HttpsProxyEnvVariableKey
//		github.com/ekara-platform/engine/util.NoProxyEnvVariableKey
func CheckProxy() (httpProxy string, httpsProxy string, noProxy string) {
	httpProxy = os.Getenv(util.HttpProxyEnvVariableKey)
	httpsProxy = os.Getenv(util.HttpsProxyEnvVariableKey)
	noProxy = os.Getenv(util.NoProxyEnvVariableKey)
	return
}

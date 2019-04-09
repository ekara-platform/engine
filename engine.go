package engine

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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
	actionManager    actionManager
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
		directory: absWorkDir}

	ctx.componentManager = component.CreateComponentManager(ctx.logger, data, absWorkDir)
	ctx.ansibleManager = ansible.CreateAnsibleManager(ctx.logger, ctx.componentManager)
	ctx.actionManager = CreateActionManager()
	return &ctx, nil
}

func (ctx *context) Init(repo string, ref string, descriptor string) error {
	wdUrl, err := getCurrentDirectoryURL(ctx)
	if err != nil {
		return err
	}
	wdUrl, err = model.NormalizeUrl(wdUrl)
	if err != nil {
		return err
	}

	// Register main component
	mainComponent, err := model.CreateComponent(wdUrl, "__main__", repo, ref)
	if err != nil {
		return err
	}
	if descriptor == "" {
		ctx.componentManager.RegisterComponent(mainComponent)
	} else {
		ctx.componentManager.RegisterComponent(mainComponent, descriptor)
	}
	if err != nil {
		return err
	}

	// Ensure the main component is present
	err = ctx.componentManager.Ensure()
	if err != nil {
		return err
	}

	return nil
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

// BuildDescriptorUrl builds the url of environment descriptor based on the
// url received has parameter
func BuildDescriptorUrl(url url.URL, fileName string) url.URL {
	if strings.HasSuffix(url.Path, "/") {
		url.Path = url.Path + fileName
	} else {
		url.Path = url.Path + "/" + fileName
	}
	return url
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

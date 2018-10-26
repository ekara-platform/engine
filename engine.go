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
	_ "gopkg.in/yaml.v2"
)

type Engine interface {
	Init(repo string, ref string, descriptor string) error
	Logger() *log.Logger
	BaseDir() string
	ComponentManager() component.ComponentManager
	AnsibleManager() ansible.AnsibleManager
	Environment() model.Environment
}

type context struct {
	// Base attributes
	logger    *log.Logger
	directory string

	// Ekara environment
	environment *model.Environment

	// Subsystems
	componentManager component.ComponentManager
	ansibleManager   ansible.AnsibleManager
}

// Create creates an environment descriptor based on the provided location.
//
// The location can be an URL over http or https or even a file system location.
//
//	Parameters:
//		logger: the logger
//		baseDir: the directory where the environment will take place among its
//				 inclusions and related components
func Create(logger *log.Logger, baseDir string, data map[string]interface{}) (Engine, error) {
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, err
	}

	ctx := context{
		logger:      logger,
		directory:   absBaseDir,
		environment: &model.Environment{}}

	ctx.componentManager = component.CreateComponentManager(ctx.logger, ctx.environment, data, absBaseDir)
	ctx.ansibleManager = ansible.CreateAnsibleManager(ctx.logger, ctx.componentManager)

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
	var cDescriptor = descriptor
	if descriptor == "" {
		cDescriptor = util.DescriptorFileName
	}
	ctx.componentManager.RegisterComponent(mainComponent, cDescriptor)
	if err != nil {
		return err
	}

	// Ensure the main component is present
	err = ctx.componentManager.Ensure()
	if err != nil {
		return err
	}

	// Register the core component
	ctx.logger.Println("Registering core")
	ctx.componentManager.RegisterComponent(ctx.environment.Ekara.Component.Resolve(), util.DescriptorFileName)

	// Register the orchestrator component
	ctx.logger.Println("Registering orchestrator")
	ctx.componentManager.RegisterComponent(ctx.environment.Orchestrator.Component.Resolve(), util.DescriptorFileName)

	// Register provider components
	for pName, pComp := range ctx.environment.Providers {
		ctx.logger.Println("Registering provider " + pName)
		ctx.componentManager.RegisterComponent(pComp.Component.Resolve(), util.DescriptorFileName)
	}

	// Register stack components
	for sName, sComp := range ctx.environment.Stacks {
		ctx.logger.Println("Registering stack " + sName)
		ctx.componentManager.RegisterComponent(sComp.Component.Resolve(), util.DescriptorFileName)
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

func (ctx *context) Environment() model.Environment {
	return *ctx.environment
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
//		engine.HttpProxyEnvVariableKey
//		engine.HttpsProxyEnvVariableKey
//		engine.NoProxyEnvVariableKey
func CheckProxy() (httpProxy string, httpsProxy string, noProxy string) {
	httpProxy = os.Getenv(util.HttpProxyEnvVariableKey)
	httpsProxy = os.Getenv(util.HttpsProxyEnvVariableKey)
	noProxy = os.Getenv(util.NoProxyEnvVariableKey)
	return
}

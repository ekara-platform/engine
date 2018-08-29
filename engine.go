package engine

import (
	"fmt"
	"hash/crc64"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lagoon-platform/model"
	_ "gopkg.in/yaml.v2"
)

type Lagoon interface {
	Init(repo string, ref string) error
	Environment() model.Environment
	ComponentManager() ComponentManager
}

type context struct {
	// Base attributes
	logger    *log.Logger
	directory string

	// Lagoon environment
	environment *model.Environment
	data        map[string]interface{}

	// Subsystems
	componentManager ComponentManager
}

// Create creates an environment descriptor based on the provided location.
//
// The location can be an URL over http or https or even a file system location.
//
//	Parameters:
//		logger: the logger
//		baseDir: the directory where the environment will take place among its
//				 inclusions and related components
func Create(logger *log.Logger, baseDir string, data map[string]interface{}) (Lagoon, error) {
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, err
	}
	ctx := context{
		logger:      logger,
		directory:   absBaseDir,
		environment: &model.Environment{},
		data:        data,
	}
	ctx.componentManager = createComponentManager(&ctx)
	return &ctx, nil
}

func (ctx *context) Init(repo string, ref string) error {

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	absWd, err := filepath.Abs(wd)
	if err != nil {
		return err
	}
	wdUrl, err := url.Parse("file://" + absWd)
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
	ctx.componentManager.RegisterComponent(mainComponent)

	// Ensure the main component is present
	err = ctx.componentManager.Ensure()
	if err != nil {
		return err
	}

	// Register the core component
	ctx.logger.Println("Registering core")
	ctx.componentManager.RegisterComponent(ctx.environment.Lagoon.Component.Resolve())

	// Register the orchestrator component
	ctx.logger.Println("Registering orchestrator")
	ctx.componentManager.RegisterComponent(ctx.environment.Orchestrator.Component.Resolve())

	// Register provider components
	for pName, pComp := range ctx.environment.Providers {
		ctx.logger.Println("Registering provider " + pName)
		ctx.componentManager.RegisterComponent(pComp.Component.Resolve())
	}

	// Register stack components
	for sName, sComp := range ctx.environment.Stacks {
		ctx.logger.Println("Registering stack " + sName)
		ctx.componentManager.RegisterComponent(sComp.Component.Resolve())
	}

	// Ensure all components are present
	err = ctx.componentManager.Ensure()
	if err != nil {
		return err
	}

	// Use context as Lagoon facade
	return nil
}

func (ctx *context) Environment() model.Environment {
	return *ctx.environment
}

func (ctx *context) ComponentManager() ComponentManager {
	return ctx.componentManager
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

func GetUId() string {
	sIp := getOutboundIP().String()
	sTime := time.Now().UTC().String()
	s := sIp + sTime
	aStringToHash := []byte(s)
	crc64Table := crc64.MakeTable(0xC96C5795D7870F42)
	crc64Int := crc64.Checksum(aStringToHash, crc64Table)
	return strconv.FormatUint(crc64Int, 16)
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

//SaveFile saves the given bytes into a fresh new file specified by its folder
//and name.
//
//If the file already exists then it will be replaced.
func SaveFile(logger *log.Logger, folder FolderPath, name string, b []byte) (string, error) {
	l := filepath.Join(folder.Path(), name)
	logger.Printf(LOG_SAVING, l)
	os.Remove(l)
	if _, err := os.Stat(name); os.IsNotExist(err) {
		e := os.MkdirAll(folder.Path(), 0700)
		if e != nil {
			logger.Printf(ERROR, e.Error())
			return l, e
		}

		logger.Printf(LOG_CREATING_FILE, l)

		f, e := os.Create(l)
		if e != nil {
			logger.Printf(ERROR, e.Error())
			return l, fmt.Errorf(ERROR_CREATING_CONFIG_FILE, name, e.Error())
		}
		defer f.Close()
		_, e = f.Write(b)
		if e != nil {
			return l, e
		}
	}
	return l, nil
}

//CheckProxy returns the proxy setting from environment variables
//
// See:
//		engine.HttpProxyEnvVariableKey
//		engine.HttpsProxyEnvVariableKey
//		engine.NoProxyEnvVariableKey
func CheckProxy() (httpProxy string, httpsProxy string, noProxy string) {
	httpProxy = os.Getenv(HttpProxyEnvVariableKey)
	httpsProxy = os.Getenv(HttpsProxyEnvVariableKey)
	noProxy = os.Getenv(NoProxyEnvVariableKey)
	return
}

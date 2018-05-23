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
	"gopkg.in/yaml.v2"
)

type Lagoon interface {
	Environment() model.Environment
	ComponentManager() ComponentManager
}

type context struct {
	logger      *log.Logger
	workDir     string
	environment model.Environment

	// Subsystems
	componentManager ComponentManager
}

// Create creates an environment descriptor based on the provider location.
//
// The location can be an URL over http or https or even a file system location.
func Create(logger *log.Logger, baseDir string, location string, tag string) (engine Lagoon, err error) {
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return
	}

	ctx := context{
		logger:  logger,
		workDir: absBaseDir}

	// Create component manager
	ctx.componentManager, err = createComponentManager(&ctx)
	if err != nil {
		return
	}

	// Fetch the main component
	envPath, err := ctx.componentManager.Fetch(location, tag)
	if err != nil {
		return
	}

	// Parse the environment descriptor from the main component
	ctx.environment, err = model.Parse(logger, filepath.Join(envPath, DescriptorFileName))
	if err != nil {
		switch err.(type) {
		case model.ValidationErrors:
			err.(model.ValidationErrors).Log(ctx.logger)
			if err.(model.ValidationErrors).HasErrors() {
				return
			}
		default:
			return
		}
	}

	// Register all environment components
	for pName, pComp := range ctx.environment.Providers {
		ctx.logger.Println("Registering provider " + pName)
		ctx.componentManager.RegisterComponent(pComp.Component)
	}

	// Use context as Lagoon facade
	engine = &ctx

	return
}

// BuildDescriptorUrl builds the url of environment descriptor based on the
// url received has parameter
//
func BuildDescriptorUrl(url url.URL) url.URL {
	if strings.HasSuffix(url.Path, "/") {
		url.Path = url.Path + DescriptorFileName
	} else {
		url.Path = url.Path + "/" + DescriptorFileName
	}
	return url
}

func (c *context) Environment() model.Environment {
	return c.environment
}

func (c *context) ComponentManager() ComponentManager {
	return c.componentManager
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
func SaveFile(logger *log.Logger, folder string, name string, b []byte) error {
	l := filepath.Join(folder, name)
	logger.Printf(LOG_SAVING, l)
	os.Remove(l)
	if _, err := os.Stat(name); os.IsNotExist(err) {
		e := os.MkdirAll(folder, 0700)
		if e != nil {
			logger.Printf(ERROR, e.Error())
			return e
		}

		logger.Printf(LOG_CREATING_CONFIG_FILE, l)

		f, e := os.Create(l)
		if e != nil {
			logger.Printf(ERROR, e.Error())
			return fmt.Errorf(ERROR_CREATING_CONFIG_FILE, name, e.Error())
		}
		defer f.Close()
		_, e = f.Write(b)
		if e != nil {
			return e
		}
	}
	return nil
}

// Proxy describes the structure used to Marshal the content of the proxy file configuration
type Proxy struct {
	// The root of the proxy specification
	ProxyEnv      ProxyEnv `yaml:"proxy_env"`
	ProxyHost     string   `yaml:"proxy_host"`
	ProxyPort     string   `yaml:"proxy_port"`
	ProxyUser     string   `yaml:"proxy_user"`
	ProxyPassword string   `yaml:"proxy_password"`
}

// ProxyEnv contains the proxy file configuration
type ProxyEnv struct {
	// The "HTTP PROXY" specification
	Http string `yaml:"http_proxy"`
	// The "HTTPS PROXY" specification
	Https string `yaml:"https_proxy"`
	// The "NO PROXY" specification
	No string `yaml:"no_proxy"`
}

//proxyConfig returns the content of the proxy configuration file
func proxyConfig(http string, https string, no string) (b []byte, e error) {
	proxy := Proxy{ProxyEnv: ProxyEnv{Http: http, Https: https, No: no}, ProxyHost: "http.internetpsa.inetpsa.com", ProxyPort: "80", ProxyUser: "mzplagww", ProxyPassword: "wwlag00n"}
	b, e = yaml.Marshal(&proxy)
	return
}

// SaveProxy creates and saves the proxy configuration file.
//
// The file will be saved into the given folder with the name:
//  engine.ProxyConfigFileName
func SaveProxy(logger *log.Logger, folder string, httpProxy string, httpsProxy string, noProxy string) error {
	b, e := proxyConfig(httpProxy, httpsProxy, noProxy)
	if e != nil {
		logger.Printf(ERROR_GENERATING_PROXY_CONFIG, e.Error())
		return e
	}

	e = SaveFile(logger, folder, ProxyConfigFileName, b)
	if e != nil {
		return e
	}
	return nil
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

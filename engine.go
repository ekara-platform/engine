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
	Environment() model.Environment
	ComponentManager() ComponentManager
}

type context struct {
	// Base attributes
	logger    *log.Logger
	directory string

	// Lagoon environment
	environment model.Environment

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
//		location: the location of the environment descriptor
//		ref: the tag/branch reference to use
func Create(logger *log.Logger, baseDir string, location string, ref string) (Lagoon, error) {
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, err
	}

	locationUrl, err := url.Parse(location)
	if err != nil {
		return nil, err
	}
	locationUrl, err = model.NormalizeUrl(locationUrl)
	if err != nil {
		return nil, err
	}

	ctx := context{
		logger:    logger,
		directory: absBaseDir}

	// Create component manager
	ctx.componentManager = createComponentManager(&ctx)

	if ref == "" {
		// Try to directly parse the descriptor if no ref is provided
		ctx.logger.Println("parsing descriptor at " + location)
		ctx.environment, err = model.Parse(logger, model.EnsurePathSuffix(locationUrl, DescriptorFileName))
	}
	if ref != "" || err != nil {
		// If no ref is provided or direct parsing is not possible, try fetching the repository
		if err != nil {
			ctx.logger.Println("descriptor is not directly parsable, fetching repository at " + location)
		} else {
			ctx.logger.Println("fetching repository at " + location)
		}
		var envUrl *url.URL
		envUrl, err = ctx.componentManager.Fetch(location, ref)
		if err != nil {
			return nil, err
		}
		ctx.environment, err = model.Parse(logger, model.EnsurePathSuffix(envUrl, DescriptorFileName))
	}

	// If only warnings are issued, allow to continue
	if err != nil {
		switch err.(type) {
		case model.ValidationErrors:
			err.(model.ValidationErrors).Log(ctx.logger)
			if err.(model.ValidationErrors).HasErrors() {
				return nil, err
			}
		default:
			return nil, err
		}
	}

	// Register all environment components
	for pName, pComp := range ctx.environment.Providers {
		ctx.logger.Println("Registering provider " + pName)
		ctx.componentManager.RegisterComponent(pComp.Component)
	}

	// Use context as Lagoon facade
	return &ctx, nil
}

// BuildDescriptorUrl builds the url of environment descriptor based on the
// url received has parameter
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
func SaveFile(logger *log.Logger, folder FolderPath, name string, b []byte) error {
	l := filepath.Join(folder.Path(), name)
	logger.Printf(LOG_SAVING, l)
	os.Remove(l)
	if _, err := os.Stat(name); os.IsNotExist(err) {
		e := os.MkdirAll(folder.Path(), 0700)
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

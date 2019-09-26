package util

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/ekara-platform/model"
)

type (
	//MockLaunchContext simulates the LaunchContext for testing purposes
	MockLaunchContext struct {
		exFolder             ExchangeFolder
		logger               *log.Logger
		qualifiedName        string
		locationContent      string
		sshPublicKeyContent  string
		sshPrivateKeyContent string
		name                 string
		user                 string
		password             string
		parameters           model.Parameters
		ekError              error
	}
)

func CreateMockLaunchContext(mainPath string, withLog bool) LaunchContext {
	p, _ := model.CreateParameters(map[string]interface{}{})
	return CreateMockLaunchContextWithData(mainPath, p, withLog)
}

func CreateMockLaunchContextWithData(mainPath string, params model.Parameters, withLog bool) LaunchContext {
	ef, _ := CreateExchangeFolder("dummy", "dummy")
	c := CreateMockLaunchContextWithDataAndFolder(mainPath, params, ef, withLog)
	c.(*MockLaunchContext).sshPublicKeyContent = "sshPublicKey_content"
	c.(*MockLaunchContext).sshPrivateKeyContent = "sshPrivateKey_content"
	return c
}

func CreateMockLaunchContextWithDataAndFolder(mainPath string, params model.Parameters, ef ExchangeFolder, withLog bool) LaunchContext {
	c := MockLaunchContext{
		locationContent: mainPath,
		parameters:      params,
		exFolder:        ef,
	}
	if withLog {
		c.logger = log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	} else {
		c.logger = log.New(ioutil.Discard, "TEST: ", log.Ldate|log.Ltime)
	}
	return &c
}

//Name simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Name() string {
	return lC.name
}

//User simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) User() string {
	return lC.user
}

//Password simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Password() string {
	return lC.password
}

//Log simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Log() *log.Logger {
	return lC.logger
}

//Ef simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Ef() ExchangeFolder {
	return lC.exFolder
}

//QualifiedName simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) QualifiedName() string {
	return lC.qualifiedName
}

//Location simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Location() string {
	return lC.locationContent
}

//Proxy simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Proxy() model.Proxy {
	return model.Proxy{}
}

//SSHPublicKey simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) SSHPublicKey() string {
	return lC.sshPublicKeyContent
}

//SSHPrivateKey simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) SSHPrivateKey() string {
	return lC.sshPrivateKeyContent
}

//Error simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Error() error {
	return lC.ekError
}

//ParamsFile simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) ParamsFile() model.Parameters {
	return lC.parameters
}

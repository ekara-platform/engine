package util

import (
	"github.com/ekara-platform/engine/model"

	"io/ioutil"
	"log"
	"os"
)

type (
	//MockLaunchContext simulates the LaunchContext for testing purposes
	MockLaunchContext struct {
		logger               *log.Logger
		fN                   FeedbackNotifier
		ef                   ExchangeFolder
		externalVars         model.Parameters
		sshPublicKeyContent  string
		sshPrivateKeyContent string
	}
)

func CreateMockLaunchContext(withLog bool) LaunchContext {
	p := model.CreateParameters(map[string]interface{}{})
	return CreateMockLaunchContextWithData(p, withLog)
}

func CreateMockLaunchContextWithData(params model.Parameters, withLog bool) LaunchContext {
	ef, _ := CreateExchangeFolder("dummy", "dummy")
	c := CreateMockLaunchContextWithDataAndFolder(params, ef, withLog)
	c.(*MockLaunchContext).sshPublicKeyContent = "sshPublicKey_content"
	c.(*MockLaunchContext).sshPrivateKeyContent = "sshPrivateKey_content"
	return c
}

func CreateMockLaunchContextWithDataAndFolder(params model.Parameters, ef ExchangeFolder, withLog bool) LaunchContext {
	c := MockLaunchContext{
		externalVars: params,
		ef:           ef,
	}
	if withLog {
		c.logger = log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	} else {
		c.logger = log.New(ioutil.Discard, "TEST: ", log.Ldate|log.Ltime)
	}
	c.fN = CreateLoggingProgressNotifier(c.logger)
	return &c
}

//Skip simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Skipping() int {
	return 0
}

//Verbosity simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Verbosity() int {
	return 0
}

//Progress simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Feedback() FeedbackNotifier {
	return lC.fN
}

//Log simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Log() *log.Logger {
	return lC.logger
}

//Ef simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Ef() ExchangeFolder {
	return lC.ef
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

//ParamsFile simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) ExternalVars() model.Parameters {
	return lC.externalVars
}

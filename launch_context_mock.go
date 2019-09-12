package engine

import (
	"log"

	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
)

type (
	//MockLaunchContext simulates the LaunchContext for testing purposes
	MockLaunchContext struct {
		efolder              *util.ExchangeFolder
		logger               *log.Logger
		qualifiedNameContent string
		locationContent      string
		sshPublicKeyContent  string
		sshPrivateKeyContent string
		engine               Engine
		name                 string
		user                 string
		password             string
		data                 model.Parameters
		ekaraError           error
	}
)

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
func (lC MockLaunchContext) Ef() *util.ExchangeFolder {
	return lC.efolder
}

//Ekara simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Ekara() Engine {
	return lC.engine
}

//QualifiedName simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) QualifiedName() string {
	return lC.qualifiedNameContent
}

//Location simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Location() string {
	return lC.locationContent
}

//HTTPProxy simulates the corresponding method in LaunchContext for testing purposes
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
	return lC.ekaraError
}

//ParamsFile simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) ParamsFile() model.Parameters {
	return lC.data
}

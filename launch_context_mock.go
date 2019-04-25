package engine

import (
	"log"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/util"
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
		cliparams            ansible.ParamContent
		ekaraError           error
	}
)

//Name simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Name() string {
	return lC.name
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

//HttpProxy simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) HttpProxy() string {
	return ""
}

//HttpsProxy simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) HttpsProxy() string {
	return ""
}

//NoProxy simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) NoProxy() string {
	return ""
}

//SshPublicKey simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) SshPublicKey() string {
	return lC.sshPublicKeyContent
}

//SshPrivateKey simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) SshPrivateKey() string {
	return lC.sshPrivateKeyContent
}

//Cliparams simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Cliparams() ansible.ParamContent {
	return lC.cliparams
}

//Error simulates the corresponding method in LaunchContext for testing purposes
func (lC MockLaunchContext) Error() error {
	return lC.ekaraError
}

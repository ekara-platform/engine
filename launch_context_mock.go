package engine

import (
	"log"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/util"
)

type (
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

func (lC MockLaunchContext) Name() string {
	return lC.name
}

func (lC MockLaunchContext) Log() *log.Logger {
	return lC.logger
}

func (lC MockLaunchContext) Ef() *util.ExchangeFolder {
	return lC.efolder
}

func (lC MockLaunchContext) Ekara() Engine {
	return lC.engine
}

func (lC MockLaunchContext) QualifiedName() string {
	return lC.qualifiedNameContent
}

func (lC MockLaunchContext) Location() string {
	return lC.locationContent
}

func (lC MockLaunchContext) HttpProxy() string {
	return ""
}

func (lC MockLaunchContext) HttpsProxy() string {
	return ""
}

func (lC MockLaunchContext) NoProxy() string {
	return ""
}

func (lC MockLaunchContext) SshPublicKey() string {
	return lC.sshPublicKeyContent
}

func (lC MockLaunchContext) SshPrivateKey() string {
	return lC.sshPrivateKeyContent
}

func (c MockLaunchContext) Cliparams() ansible.ParamContent {
	return c.cliparams
}

func (c MockLaunchContext) Error() error {
	return c.ekaraError
}

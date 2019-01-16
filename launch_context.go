package engine

import (
	"log"

	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/util"
)

type (
	LaunchContext interface {
		Name() string
		Log() *log.Logger
		Ef() *util.ExchangeFolder
		Ekara() Engine
		QualifiedName() string
		Location() string
		HttpProxy() string
		HttpsProxy() string
		NoProxy() string
		SshPublicKey() string
		SshPrivateKey() string
		Cliparams() ansible.ParamContent
		Error() error
	}
)

func BuildBaseParam(c LaunchContext, nodeSetName string, provider string) ansible.BaseParam {
	return ansible.BuildBaseParam(c.Ekara().ComponentManager().Environment(), nodeSetName, provider, c.SshPublicKey(), c.SshPrivateKey())
}

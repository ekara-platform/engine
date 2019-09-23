package util

import (
	"log"

	"github.com/ekara-platform/model"
)

type (
	//LaunchContext Represents the informations required to run the engine
	LaunchContext interface {
		//Name The environment descriptor name to process
		Name() string
		//User The user to log into the descriptor repository
		User() string
		//Password The user to log into the descriptor repository
		Password() string
		//Location specifies where to look for the environment descriptor
		Location() string
		//Log the looger to used during the ekara execution
		Log() *log.Logger
		//Ef the exchange folder where to find informations required
		// to complete the Ekara execution of where to write informations
		// produced by the execution
		Ef() ExchangeFolder
		//QualifiedName The qualifie name of the environment to process
		QualifiedName() string
		//Proxy returns launch context proxy settings
		Proxy() model.Proxy
		//SSHPublicKey the public key used by the engine during the process execution to
		// connect the created nodes
		SSHPublicKey() string
		//SSHPrivateKey the private key used by the engine during the process execution to
		// connect the created nodes
		SSHPrivateKey() string
		//ParamsFile returns the content the parameters provided by the user to fill the environment descriptor as a template
		ParamsFile() model.Parameters
		//Error if any, the error which has broken the process
		Error() error
	}
)

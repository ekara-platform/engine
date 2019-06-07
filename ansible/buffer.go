package ansible

import (
	"log"

	"github.com/ekara-platform/engine/util"

	"gopkg.in/yaml.v2"
)

// Buffer contains the data to be passed from one playbook to another one
type Buffer struct {
	// The environment variables to pass
	Envvars map[string]string
	// The extra vars to pass
	Extravars map[string]string
	// The inventories to pass
	//TODO CHECK IF THIS CAN BE DELETED
	Inventories map[string]string
	// The parameters to pass.
	// The parameters can be any valid yaml content
	Param map[string]interface{}
}

// Params returns the yaml content of the buffer
func (b Buffer) Params() (by []byte, e error) {
	by, e = yaml.Marshal(b.Param)
	return
}

//CreateBuffer creates an empty buffer
func CreateBuffer() Buffer {
	b := Buffer{}
	b.Envvars = make(map[string]string)
	b.Extravars = make(map[string]string)
	b.Inventories = make(map[string]string)
	b.Param = make(map[string]interface{})
	return b
}

// GetBuffer returns a buffer based on the content of the given folder.
//
/*
 - The "Envvars" of the buffer will be filled by the content of the file name like  EnvYamlFileName.
 - The "Extravars" of the buffer will be filled by the content of the file name like  ExtraVarYamlFileName.
 - The "Inventories" of the buffer will be filled by the content of the file name like  InventoryYamlFileName.
 - The "Param" of the buffer will be filled by the content of the file name like  ParamYamlFileName.
*/
//
// Parameters:
//		f: the folder where to look for the buffer content
//		logger: the logger
//		location: a string used to specify, into the log messages, to which concept the buffer is related
//
//
// See Also :
//  EnvYamlFileName
//  ExtraVarYamlFileName
//  InventoryYamlFileName
//  ParamYamlFileName
//
func GetBuffer(f *util.FolderPath, logger *log.Logger, location string) (buffer Buffer, err error) {
	buffer = CreateBuffer()

	var ok bool
	var b []byte

	if ok, b, err = f.ContainsEnvYaml(); ok {
		logger.Printf("Consuming %s for %s", util.EnvYamlFileName, location)
		if err != nil {
			return
		}
		err = yaml.Unmarshal([]byte(b), buffer.Envvars)
		if err != nil {
			return
		}
	} else {
		logger.Printf("No %s located...", util.EnvYamlFileName)
	}

	if ok, b, err = f.ContainsExtraVarsYaml(); ok {
		logger.Printf("Consuming %s for %s", util.ExtraVarYamlFileName, location)
		if err != nil {
			return
		}
		err = yaml.Unmarshal([]byte(b), buffer.Extravars)
		if err != nil {
			return
		}
	} else {
		logger.Printf("No %s located...", util.ExtraVarYamlFileName)
	}

	if ok, b, err = f.ContainsInventoryYamlFileName(); ok {
		logger.Printf("Consuming %s for %s", util.InventoryYamlFileName, location)
		if err != nil {
			return
		}
		err = yaml.Unmarshal([]byte(b), buffer.Inventories)
		if err != nil {
			return
		}
	} else {
		logger.Printf("No %s located...", util.InventoryYamlFileName)
	}

	if ok, b, err = f.ContainsParamYaml(); ok {
		logger.Printf("Consuming %s for  %s", util.ParamYamlFileName, location)
		if err != nil {
			return
		}
		err = yaml.Unmarshal([]byte(b), buffer.Param)
		if err != nil {
			logger.Printf("Error consuming the param %s", err.Error())
			return
		}
	} else {
		logger.Printf("No %s located...", util.ParamYamlFileName)
	}
	return
}

func (b Buffer) isEmpy() bool {
	return !b.hasEnvvars() && !b.hasExtravars() && !b.hasInventories() && !b.hasParam()
}

func (b Buffer) hasEnvvars() bool {
	return len(b.Envvars) == 0
}

func (b Buffer) hasExtravars() bool {
	return len(b.Extravars) == 0
}

func (b Buffer) hasInventories() bool {
	return len(b.Inventories) == 0
}

func (b Buffer) hasParam() bool {
	return len(b.Param) == 0
}

package engine

import (
	"log"

	"gopkg.in/yaml.v2"
)

// Buffer contains the data to be passed from one playbook to another one
type Buffer struct {
	// The environment variables to pass
	Envvars map[string]string
	// The extra vars to pass
	Extravars map[string]string
	// The parameters to pass.
	// The parameters can be any valid yaml content
	Param map[string]interface{}
}

// Params returns the yaml content of the buffer
func (bu Buffer) Params() (b []byte, e error) {
	b, e = yaml.Marshal(bu.Param)
	return
}

//CreateBuffer creates an empty buffer
func CreateBuffer() Buffer {
	b := Buffer{}
	b.Envvars = make(map[string]string)
	b.Extravars = make(map[string]string)
	b.Param = make(map[string]interface{})
	return b
}

// GetBuffer returns a buffer based on the content of the given folder.
//
/*
 - The "Envvars" of the buffer will be filled by the content of the file name like  EnvYamlFileName
 - The "Extravars" of the buffer will be filled by the content of the file name like  ExtraVarYamlFileName
 - The "Param" of the buffer will be filled by the content of the file name like  ParamYamlFileName
*/
//
//
//	Parameters:
//		f: the folder where to look for the buffer content
//		logger: the logger
//		location: a string used to specify, into the log messages, to which concept the buffer is related
//
//
// See Also :
//  EnvYamlFileName
//  ExtraVarYamlFileName
//  ParamYamlFileName
//
func GetBuffer(f *FolderPath, logger *log.Logger, location string) (err error, buffer Buffer) {
	buffer = CreateBuffer()

	var ok bool
	var b []byte

	if ok, b, err = f.ContainsEnvYaml(); ok {
		logger.Printf("Consuming %s for %s", EnvYamlFileName, location)
		if err != nil {
			return
		}
		err = yaml.Unmarshal([]byte(b), buffer.Envvars)
		if err != nil {
			return
		}
	} else {
		logger.Printf("No %s located...", EnvYamlFileName)
	}

	if ok, b, err = f.ContainsExtraVarsYaml(); ok {
		logger.Printf("Consuming %s for %s", ExtraVarYamlFileName, location)
		if err != nil {
			return
		}
		err = yaml.Unmarshal([]byte(b), buffer.Extravars)
		if err != nil {
			return
		}
	} else {
		logger.Printf("No %s located...", ExtraVarYamlFileName)
	}

	if ok, b, err = f.ContainsParamYaml(); ok {
		logger.Printf("Consuming %s for  %s", ParamYamlFileName, location)
		if err != nil {
			return
		}
		err = yaml.Unmarshal([]byte(b), buffer.Param)
		if err != nil {
			logger.Printf("Error consuming the param %s", err.Error())
			return
		}
	} else {
		logger.Printf("No %s located...", ParamYamlFileName)
	}
	return
}

package ansible

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/ekara-platform/model"
	"gopkg.in/yaml.v2"
)

type ParamValues map[string]string

type keyValue struct {
	key   string
	value string
}

// BaseParam contains the extra vars to be passed to a playbook
//
// The BaseParam content is supposed to be serialized in yaml in order to be passed
// to a playbook
type BaseParam struct {
	// The content of the BaseParam
	Body model.Parameters
}

// BuildBaseParam the common parameters required by all playbooks.
//
// Parameters:
//		name: the qualified name of the environment
//		nodesetName: the name of the nodeset we are working with
//		provider: the name of the provider where to create the nodeset
//		pubK: the public SSH key to connect on the created nodeset ( the name of the file)
//		privK: the private SSH key to connect on the created nodeset ( the name of the file)
//
func BuildBaseParam(env *model.Environment, pubK string, privK string, nodesetName string) BaseParam {
	baseParam := BaseParam{}
	baseParam.Body = make(map[string]interface{})

	sshM := make(map[string]interface{})
	if pubK != "" {
		sshM["public_key"] = pubK
	}
	if privK != "" {
		sshM["private_key"] = privK
	}
	baseParam.Body["ssh"] = sshM

	environmentM := make(map[string]interface{})
	if env.QualifiedName().String() != "" {
		environmentM["id"] = env.QualifiedName().String()
	}
	if env.Name != "" {
		environmentM["name"] = env.Name
	}
	if env.Qualifier != "" {
		environmentM["qualifier"] = env.Qualifier
	}
	if nodesetName != "" {
		environmentM["nodeset"] = nodesetName
	}
	baseParam.Body["environment"] = environmentM

	return baseParam
}

//AddNamedMap adds a parameter of type map[string]interface{} for the given name
func (bp *BaseParam) AddNamedMap(name string, c map[string]interface{}) {
	bp.Body[name] = c
}

//AddMap adds parameters of interface{} for all the given map entries
func (bp *BaseParam) AddMap(m map[string]interface{}) {
	for k, v := range m {
		bp.AddInterface(k, v)
	}
}

//AddInterface adds a parameter of type interface{} for the given name
func (bp *BaseParam) AddInterface(name string, i interface{}) {
	bp.Body[name] = i
}

//AddInt adds a int parameter for the given name
func (bp *BaseParam) AddInt(name string, i int) {
	bp.Body[name] = i
}

//AddString adds a string parameter for the given name
func (bp *BaseParam) AddString(name string, s string) {
	bp.Body[name] = s
}

// AddBuffer adds the parameters coming from the given buffer
//
// Only the "Param" content of the buffer will be processed.
func (bp *BaseParam) AddBuffer(b Buffer) {
	bp.AddMap(b.Param)
}

// Content returns the yaml representation of the content
func (bp BaseParam) Content() (b []byte, e error) {
	b, e = yaml.Marshal(&bp.Body)
	return
}

// Copy returns a copy of the base parameter
func (bp BaseParam) Copy() BaseParam {
	result := BaseParam{}
	result.Body = make(map[string]interface{})
	for k, v := range bp.Body {
		val, okType := v.(map[string]interface{})
		if okType {
			m := make(map[string]interface{})
			for kVal, vVal := range val {
				m[kVal] = vVal
			}
			result.Body[k] = m
		} else {
			result.Body[k] = v
		}
	}
	return result
}

//ParseParamValues parses a yaml file into a map of "key:value"
// All the nested yaml levels will be concatenated to create the keys of the map.
//
// As example :
// 	level1:
// 	  level2:value
//
// Will generate the following key/value
// 	level1.level2=value
func ParseParamValues(path string) (ParamValues, error) {
	r := make(map[string]string)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return r, err
	}
	cKv := make(chan keyValue)
	exit := make(chan bool)
	env := make(map[interface{}]interface{})
	err = yaml.Unmarshal(b, &env)
	if err != nil {
		panic(err)
	}
	go readMap(cKv, exit, "", env)

	for {
		select {
		case <-exit:
			return r, nil
		case kv := <-cKv:
			r[kv.key] = kv.value
		}
	}
}

func readMap(cKv chan keyValue, exit chan bool, location string, m map[interface{}]interface{}) {
	if location != "" {
		location += "."
	}

	for k, v := range m {
		ks := fmt.Sprintf("%v", k)
		if reflect.ValueOf(v).Kind() == reflect.Map {
			// The value is a map so we go deeper...
			readMap(cKv, exit, location+ks, v.(map[interface{}]interface{}))
		} else {
			if v != nil {
				vs := fmt.Sprintf("%v", v)
				// The located leaf is returned...
				cKv <- keyValue{
					key:   location + ks,
					value: vs,
				}
			}
		}
	}
	// Only the root map can trigger the exit
	if location == "" {
		exit <- true
	}
}

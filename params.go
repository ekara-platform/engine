package engine

import (
	"gopkg.in/yaml.v2"
)

type ParamContent map[string]interface{}

// BaseParam contains the extra vars to be passed to a playbook
//
// The BaseParam content is supposed to be serialized in yaml in order to be passed
// to a playbook
type BaseParam struct {
	// The content of the BaseParam
	Body ParamContent
}

// BuilBaseParam the common parameters required by all playbooks.
//
// Parameters:
//		client: the name of the client
//		uid: the unique id of the nodeset we are working with
//		provider: the name of the provider where to create the nodeset
//		pubK: the public SSH key to connect on the created nodeset ( the name of the file)
//		privK: the private SSH key to connect on the created nodeset ( the name of the file)
//
func BuilBaseParam(client string, uid string, provider string, pubK string, privK string) BaseParam {

	baseParam := BaseParam{}
	baseParam.Body = make(map[string]interface{})

	connectionM := make(map[string]interface{})
	if provider != "" {
		connectionM["provider"] = provider
	}
	if pubK != "" {
		connectionM["machine_public_key"] = pubK
	}
	if privK != "" {
		connectionM["machine_private_key"] = privK
	}
	baseParam.Body["connectionConfig"] = connectionM

	clientM := make(map[string]interface{})
	if client != "" {
		clientM["name"] = client
	}
	if uid != "" {
		clientM["uid"] = uid
	}
	baseParam.Body["client"] = clientM

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

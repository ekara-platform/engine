package engine

import (
	"gopkg.in/yaml.v2"
)

type ParamContent map[string]interface{}

type BaseParam struct {
	Body ParamContent
}

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

func (bp *BaseParam) AddNamedMap(name string, c map[string]interface{}) {
	bp.Body[name] = c
}

func (bp *BaseParam) AddMap(m map[string]interface{}) {
	for k, v := range m {
		bp.AddInterface(k, v)
	}
}

func (bp *BaseParam) AddInterface(name string, i interface{}) {
	bp.Body[name] = i
}

func (bp *BaseParam) AddInt(name string, i int) {
	bp.Body[name] = i
}

func (bp *BaseParam) AddString(name string, s string) {
	bp.Body[name] = s
}

// AddBuffer adds the parameters coming from the given buffer
func (bp *BaseParam) AddBuffer(b Buffer) {
	bp.AddMap(b.Param)
}

func (bp BaseParam) Content() (b []byte, e error) {
	b, e = yaml.Marshal(&bp.Body)
	return
}

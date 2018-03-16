package main

import (
	"encoding/json"
	"log"
	"os"
	_ "os/exec"
	"text/template"
)

type Interface struct {

	// The name of the generated interface
	Name string `json:"interface_name"`

	// The list of inherited interfaces
	Inherited []string `json:"inherited_interfaces"`

	// The list of methods exposed by the interface
	Methods []Method `json:"exposed_methods"`
}

type Method struct {

	// The type returned by the method ( if missing the returned type will be "string" )
	Ret string `json:"return_type"`

	// The name of the method. ( if missing the generated name will be :
	// "Get" + returned type
	Name string

	// Parameters to pass to the generated method
	// Ex :
	// "parameters":"...string" will generate "(...string)"
	Params string `json:"parameters"`

	// The list of supported implementation to generate
	Implentations []Implementation `json:"implemented_by"`
}

type Implementation struct {

	// The type implementing the interface
	//
	// The type will be named "e" so if you provide a custom implementation
	// it must like :
	// "impl":"return CreateNodes(e.Nodes)"
	// or
	// "impl":"return e.DoSomething()"

	Type string `json:"invoker_type"`

	// The attribute(s), of the type implementing the interface, allowing to
	// get the value returned by the implementation:
	//
	// If specified the implentation body will be : "return e"SubType".Att"
	// If not the implentation body will be : "return e.Att"
	//
	// The SubType must start with a "."
	// Ex :
	// "invoker_attribute":".Proxy" wil generate "return e.Proxy.Att"
	// "invoker_attribute":".Proxy.Http.Blablabla" wil generate "return e..Proxy.Http.Blablabla.Att"

	SubType string `json:"invoker_attribute"`

	// 1 if the implementation must be done on a pointer (defaulted to "0")

	Pointer int `json:"on_pointer"`

	// The value returned by the implementation
	Att string `json:"attribute"`

	// The explicit implementation of the method body
	// Ex :
	// "impl":"return CreateNodes(e.Nodes)"
	Impl string `json:"impl"`
}

func (m Method) Signature() string {
	if m.Ret == "" {
		m.Ret = "string"
	}
	var r string
	if m.Name == "" {
		r = "Get" + m.Ret
	} else {
		r = m.Name
	}
	if m.Params == "" {
		r += "() "
	} else {
		r += "(" + m.Params + ") "
	}
	r += m.Ret
	return r
}

func (i Implementation) ImplSignature(m Method) string {
	r := "func (e "
	if i.Pointer == 1 {
		r += "*"
	}
	r += i.Type + ") "
	r += m.Signature()
	return r
}

func (i Implementation) Body() string {
	if i.Impl != "" {
		return i.Impl
	}

	r := "return "
	if i.Att == "" {
		r += "e" + i.SubType
	} else {
		r += "e" + i.SubType + "." + i.Att
	}
	return r
}

func main() {
	r, err := os.Open("model_interface.json")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	defer r.Close()

	v := []Interface{}
	err = json.NewDecoder(r).Decode(&v)
	if err != nil {
		panic(err)
	}

	log.Printf("loaded %v", v)

	w, err := os.Create("interfaces_generated.go")
	if err != nil {
		panic(err)
	}
	defer w.Close()

	t, err := template.ParseFiles("model_interface_template.txt")
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, v)
	if err != nil {
		panic(err)
	}
}

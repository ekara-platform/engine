package main

import (
	"encoding/json"
	"os"
	"text/template"
)

type GenString string

type Interface struct {

	// The name of the generated interface
	Name string `json:"name"`

	// The list of inherited interfaces
	Inherited []string `json:"inherited"`

	// The list of methods exposed by the interface
	Methods []Method `json:"methods"`
}

type Method struct {

	// The type returned by the method ( if missing the returned type will be "string" )
	// The returned type will be processed "as is".
	// It should include parentheses in case of multiple returns.
	// Ex:
	//    "int"
	//    "(bool, error)"
	Ret string `json:"returns"`

	// The name of the method. ( if missing the generated name will be :
	// "Get" + "type returned". )
	Name string

	// Parameters to pass to the generated method
	// Ex :
	// "parameters":"...string" will generate "(...string)"
	Params string `json:"parameters"`

	// The list of supported implementation to generate
	Implentations []Implementation `json:"implemented_by"`
}

type Implementation struct {

	// The types implementing the interface
	//
	// The type will be named "e" so if you provide a custom implementation
	// it must like :
	// "impl":"return CreateNodes(e.Nodes)"
	// or
	// "impl":"return e.DoSomething()"
	// One implementation will be generated for each type.
	Types []GenString `json:"types"`

	// The attribute(s), of the type implementing the interface, allowing to
	// get the value returned by the implementation:
	//
	// If specified the implentation body will be : "return e + ".SubType"+ ".Att"
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

func (m Method) Signature(toImpl bool) string {
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
		if toImpl {
			r += "(p " + m.Params + ") "
		} else {
			r += "(" + m.Params + ") "
		}
	}
	r += m.Ret
	return r
}

func (t GenString) ImplSignature(m Method, i Implementation) string {
	r := "func (e "
	if i.Pointer == 1 {
		r += "*"
	}
	r += string(t) + ") "
	r += m.Signature(true)
	return r
}

func (t GenString) Body(i Implementation) string {
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

func (i Interface) HasImplentations() bool {
	for _, vm := range i.Methods {
		if len(vm.Implentations) > 0 {
			return true
		}
	}
	return false
}

func main() {
	r, err := os.Open("generate/engine_api.json")
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

	w, err := os.Create("engine_api.generated.go")
	if err != nil {
		panic(err)
	}
	defer w.Close()

	t, err := template.ParseFiles("generate/engine_api.tmpl")
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, v)
	if err != nil {
		panic(err)
	}
}

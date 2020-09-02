package model

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"reflect"

	"gopkg.in/yaml.v2"
)

//ParametersAware represents the common behavior of parameters holders
type Parameterized interface {
	Parameters() Parameters
}

// Parameters represents the parameters coming from a descriptor
type Parameters map[string]interface{}

// CreateEmptyParameters builds an empty Parameters structure.
func CreateEmptyParameters() Parameters {
	return make(map[string]interface{})
}

// CreateParameters builds Parameters from the specified map
func CreateParameters(src map[string]interface{}) Parameters {
	dst := make(map[string]interface{})
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// ParseParameters parses a yaml file into a Parameters
func ParseParameters(path string) (Parameters, error) {
	r := make(Parameters)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return r, err
	}
	err = yaml.Unmarshal(b, &r)
	if err != nil {
		return r, err
	}
	return r, nil
}

// CloneParameters deep-copy the entire parameters
func CloneParameters(other Parameters) Parameters {
	cp := make(map[string]interface{})
	for k, v := range other {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = CloneParameters(vm)
		} else {
			cp[k] = v
		}
	}
	return cp
}

// TODO: terrible! change this
func (r Parameters) Override(with Parameters) Parameters {
	withG := make(map[interface{}]interface{})
	for k, v := range with {
		withG[k] = v
	}
	rG := make(map[interface{}]interface{})
	for k, v := range r {
		rG[k] = v
	}
	dst := make(map[interface{}]interface{})
	merge(dst, rG)
	merge(dst, withG)
	ret := make(map[string]interface{})
	for k, v := range dst {
		ret[fmt.Sprintf("%v", k)] = v
	}
	return ret
}

func merge(dst map[interface{}]interface{}, src map[interface{}]interface{}) {
	for k, v := range src {
		vv := reflect.ValueOf(v)
		if vv.Kind() == reflect.Map {
			// The value is a map so we try to go deeper if they have the same key type
			// Otherwise we overwrite the destination map with the source one
			vd := reflect.ValueOf(dst[k])
			if vd.Kind() != reflect.Map || vd.Type().Key() != vv.Type().Key() {
				dst[k] = make(map[interface{}]interface{})
			}
			merge(dst[k].(map[interface{}]interface{}), v.(map[interface{}]interface{}))
		} else if vv.Kind() == reflect.Slice {
			// The value is a slice so we try to concatenate if they have the same element type
			// Otherwise we overwrite the destination slice with the source one
			vd := reflect.ValueOf(dst[k])
			if vd.Kind() != reflect.Slice || vd.Type().Elem() != vv.Type().Elem() {
				dst[k] = reflect.MakeSlice(reflect.SliceOf(vv.Type().Elem()), 0, vv.Len()).Interface()
				vd = reflect.ValueOf(dst[k])
			}
			dst[k] = reflect.AppendSlice(vd, vv).Interface()
		} else {
			if v != nil {
				dst[k] = v
			}
		}
	}
}

// ToYAML returns the Parameters content as yaml
func (r Parameters) ToYAML() ([]byte, error) {
	return yaml.Marshal(r)
}

//IdentYamlMap convert parameters into yaml with content based on the given identation
func (r Parameters) IdentYamlMap(ident int) string {
	ret := ""

	cKv := make(chan string)
	exit := make(chan bool)

	go readMap(cKv, exit, ident, ident, r)

	for {
		select {
		case <-exit:
			return ret
		case s := <-cKv:
			ret = ret + s
		}
	}
}

//IndentYaml convert a content into yaml based on the given identation
func IndentYaml(ident int, v interface{}) string {
	r := ""

	cKv := make(chan string)
	exit := make(chan bool)

	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Map {
		go readMap(cKv, exit, ident, ident, v.(map[string]interface{}))
	} else if vv.Kind() == reflect.Slice {
		go readSlice(cKv, exit, ident, ident, v.([]interface{}))
	} else {
		sp := ""
		for i := 0; i < ident; i++ {
			sp = sp + " "
		}
		r = r + sp + fmt.Sprintf("%v", v)
		return r
	}

	for {
		select {
		case <-exit:
			return r
		case s := <-cKv:
			r = r + s
		}
	}
}

func readMap(cKv chan string, exit chan bool, i int, ident int, src map[string]interface{}) {
	sp := ""
	for i := 0; i < ident; i++ {
		sp = sp + " "
	}
	for k, v := range src {
		vv := reflect.ValueOf(v)
		if vv.Kind() == reflect.Map {
			cKv <- sp + fmt.Sprintf("%v:\n", k)
			readMap(cKv, exit, i, ident+2, v.(map[string]interface{}))
		} else if vv.Kind() == reflect.Slice {
			cKv <- sp + fmt.Sprintf("%v:\n", k)
			readSlice(cKv, exit, i, ident+2, v.([]interface{}))
		} else {
			cKv <- sp + fmt.Sprintf("%v: %v\n", k, v)
		}
	}
	if ident == i {
		exit <- true
	}
}

func readSlice(cKv chan string, exit chan bool, i int, ident int, src []interface{}) {
	sp := ""
	for i := 0; i < ident; i++ {
		sp = sp + " "
	}
	for _, v := range src {
		vv := reflect.ValueOf(v)
		if vv.Kind() == reflect.Map {
			readMap(cKv, exit, i, ident+2, v.(map[string]interface{}))
		} else if vv.Kind() == reflect.Slice {
			readSlice(cKv, exit, i, ident+2, v.([]interface{}))
		} else {
			cKv <- sp + fmt.Sprintf("- %v\n", v)
		}
	}
	if ident == i {
		exit <- true
	}
}

func Json(v interface{}) template.HTML {
	strB, err := json.Marshal(v)
	if err != nil {
		return template.HTML(err.Error())
	}
	return template.HTML(strB)
}

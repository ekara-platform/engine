package descriptor

import (
	"gopkg.in/yaml.v2"
	"github.com/imdario/mergo"
	"os"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"fmt"
	"strings"
)

type Task struct {
	Playbook string
	Cron     string
	RunOn    []string `yaml:"runOn"`
}

type Hook struct {
	Before []string
	After  []string
}

type Provider struct {
	parameters map[string]string
}

type NodeSet struct {
	Platform bool
	Provider string
	Instances struct {
		Min int
		Max int
	}
	Config struct {
		Cpu    int
		Memory int
		Disk   int
	}
	Hooks struct {
		Installation Hook
	}
}

type Stack struct {
	Url      string
	DeployOn []string `yaml:"deployOn"`
}

type Descriptor struct {
	Name         string
	Description  string
	Imports      []string
	BaseLocation string
	Nodes        map[string]NodeSet
	Stacks       map[string]Stack
	Tasks        map[string]Task
	Hooks struct {
		Provisioning Hook
		Installation Hook
		Deployment   Hook
	}
}

func ParseDescriptor(location string) (desc Descriptor, err error) {
	var content []byte

	desc.BaseLocation, content, err = readDescriptor(location)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(content, &desc)
	if err != nil {
		return
	}

	err = processImports(&desc)
	if err != nil {
		return
	}

	return
}

func readDescriptor(location string) (base string, content []byte, err error) {
	if strings.Index(location, "http") == 0 {
		fmt.Println("Loading URL", location)
		var response *http.Response
		response, err = http.Get(location)
		if err != nil {
			return
		}
		defer response.Body.Close()
		content, err = ioutil.ReadAll(response.Body)
		base = path.Dir(location) + "/"
	} else {
		fmt.Println("Loading file", location)
		var file *os.File
		file, err = os.Open(location)
		if err != nil {
			return
		}
		defer file.Close()
		content, err = ioutil.ReadAll(file)
		if err != nil {
			return
		}
		var absLocation string
		absLocation, err = filepath.Abs(location)
		if err != nil {
			return
		}
		base = filepath.Dir(absLocation) + string(filepath.Separator)
	}
	return
}

func processImports(desc *Descriptor) error {
	if len(desc.Imports) > 0 {
		fmt.Println("Processing imports", desc.Imports)
		for _, val := range desc.Imports {
			importedDesc, err := ParseDescriptor(desc.BaseLocation + val)
			if err != nil {
				return err
			}
			mergo.Merge(desc, importedDesc)
		}
		desc.Imports = nil
	} else {
		fmt.Println("No import to process")
	}
	return nil
}

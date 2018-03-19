package engine

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

type hookDef struct {
	Before []string
	After  []string
}

type labelsDef struct {
	Labels []string
}

type paramsDef struct {
	Params map[string]string
}

type platformDef struct {
	Version  string
	Registry string
	Proxy    struct {
		Http    string
		Https   string
		NoProxy string `yaml:"noProxy"`
	}
}

type providerDef struct {
	name      string `yaml:"-"`
	labelsDef `yaml:",inline"`
	paramsDef `yaml:",inline"`
}

type nodeSetDef struct {
	name      string `yaml:"-"`
	labelsDef `yaml:",inline"`

	Provider struct {
		paramsDef `yaml:",inline"`
		Name      string
	}
	Instances int
	Hooks     struct {
		Provision hookDef
		Destroy   hookDef
	}
}

type stackDef struct {
	name      string `yaml:"-"`
	labelsDef `yaml:",inline"`

	Repository string
	Version    string
	DeployOn   []string `yaml:"deployOn"`
	Hooks      struct {
		Deploy   hookDef
		Undeploy hookDef
	}
}

type taskDef struct {
	name      string `yaml:"-"`
	labelsDef `yaml:",inline"`

	Playbook string
	Cron     string
	RunOn    []string `yaml:"runOn"`
	Hooks    struct {
		Execute hookDef
	}
}

type environmentDef struct {
	labelsDef `yaml:",inline"`

	// Global attributes
	Name         string
	Description  string
	Version      string
	BaseLocation string

	// Imports
	Imports []string

	// PlatformDef attributes
	Lagoon platformDef

	// Providers
	Providers map[string]providerDef
	providers providers `yaml:"-"`

	// Node sets
	Nodes map[string]nodeSetDef
	nodes nodes `yaml:"-"`

	// Software stacks
	Stacks map[string]stackDef
	stacks stacks `yaml:"-"`

	// Custom tasks
	Tasks map[string]taskDef
	tasks tasks `yaml:"-"`

	// Global hooks
	Hooks struct {
		Init      hookDef
		Provision hookDef
		Deploy    hookDef
		Undeploy  hookDef
		Destroy   hookDef
	}
}

func parseDescriptor(location string) (desc environmentDef, err error) {
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

	desc.providers = CreateProviders(desc.Providers)
	desc.nodes = CreateNodes(desc.Nodes)
	desc.stacks = CreateStacks(desc.Stacks)
	desc.tasks = CreateTasks(desc.Tasks)

	return
}

func readDescriptor(location string) (base string, content []byte, err error) {
	if strings.Index(location, "http") == 0 {
		h.logger.Println("Loading URL", location)

		_, err = url.Parse(location)
		if err != nil {
			return
		}

		var response *http.Response
		response, err = http.Get(location)
		if err != nil {
			return
		}
		defer response.Body.Close()
		content, err = ioutil.ReadAll(response.Body)

		i := strings.LastIndex(location, "/")
		base = location[0 : i+1]
	} else {
		h.logger.Println("Loading file", location)
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

func processImports(desc *environmentDef) error {
	if len(desc.Imports) > 0 {
		h.logger.Println("Processing imports", desc.Imports)
		for _, val := range desc.Imports {
			importedDesc, err := parseDescriptor(desc.BaseLocation + val)
			if err != nil {
				return err
			}
			mergo.Merge(desc, importedDesc)
		}
		desc.Imports = nil
	} else {
		h.logger.Println("No import to process")
	}
	return nil
}

type namedMap map[string]interface{}

func mapContains(m namedMap, candidate string) bool {
	_, ok := m[candidate]
	return ok
}

func mapMultipleContains(m namedMap, candidates []string) bool {
	for _, c := range candidates {
		if b := mapContains(m, c); b == false {
			return false
		}
	}
	return true
}

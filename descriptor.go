package engine

import (
	"fmt"
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
	name string          `yaml:"-"`
	desc *environmentDef `yaml:"-"`

	labelsDef `yaml:",inline"`
	paramsDef `yaml:",inline"`
}

type nodeSetDef struct {
	name string          `yaml:"-"`
	desc *environmentDef `yaml:"-"`

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
	name string          `yaml:"-"`
	desc *environmentDef `yaml:"-"`

	labelsDef `yaml:",inline"`

	Repository string
	Version    string
	DeployOn   deployDef `yaml:",inline"`
	Hooks      struct {
		Deploy   hookDef
		Undeploy hookDef
	}
}

type deployDef struct {
	Names []string `yaml:"deployOn"`
}

type taskDef struct {
	name      string          `yaml:"-"`
	desc      *environmentDef `yaml:"-"`
	labelsDef `yaml:",inline"`

	Task  string
	Cron  string
	RunOn runDef `yaml:",inline"`
	Hooks struct {
		Execute hookDef
	}
}

type runDef struct {
	Names []string `yaml:"runOn"`
}

type environmentDef struct {
	labelsDef `yaml:",inline"`

	// Global attributes
	Name                  string
	Description           string
	Version               string
	BaseLocation          string
	TreatWarningsAsErrors bool `yaml:"treatWarningsAsErrors"`

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

// parseDescriptor parses the descriptor based on a location.
// The location can be an URL ( HTTP and HTTPS are supported )
// or a location on the local file system
func parseDescriptor(h holder, location string) (desc environmentDef, err error) {
	var content []byte
	desc.BaseLocation, content, err = readDescriptor(h, location)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(content, &desc)
	if err != nil {
		return
	}

	err = processImports(h, &desc)
	if err != nil {
		return
	}
	return
}

// parseDescriptor parses the descriptor based previously serialized content.
func parseContent(h holder, content []byte) (desc environmentDef, err error) {
	err = yaml.Unmarshal(content, &desc)
	if err != nil {
		return
	}
	// Note the parseContent won't process any import because they have been
	// resolved before the serialization
	return
}

// adjustAndValidate adjusts the descriptor content and validate its grammar and
// its required content
func (desc *environmentDef) adjustAndValidate() (ge GrammarErrors) {
	desc.providers = CreateProviders(desc)
	desc.nodes = CreateNodes(desc)
	desc.stacks = CreateStacks(desc)
	desc.tasks = CreateTasks(desc)
	ge = *desc.validate()
	return
}

/*
getWarningType returns the error type correponding to natural warning.

The returned ErrorType will depend on "TreatWarningsAsErrors"

	TreatWarningsAsErrors = true then a "Error" will be returned
	TreatWarningsAsErrors = false then a "Warning" will be returned

*/
func (desc *environmentDef) getWarningType() ErrorType {
	if desc.TreatWarningsAsErrors {
		return Error
	}
	return Warning
}

// readDescriptor reads the descriptor based on a location.
// The location can be an URL ( HTTP and HTTPS are supported )
// or a location on the local file system
func readDescriptor(h holder, location string) (base string, content []byte, err error) {
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

		if response.StatusCode != 200 {
			err = fmt.Errorf("HTTP Error getting the environment descriptor , error code %d", response.StatusCode)
			return
		}

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

func processImports(h holder, desc *environmentDef) error {
	if len(desc.Imports) > 0 {
		h.logger.Println("Processing imports", desc.Imports)
		for _, val := range desc.Imports {
			importedDesc, err := parseDescriptor(h, desc.BaseLocation+val)
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

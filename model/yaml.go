package model

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type (
	// yaml tag for the proxy details
	yamlProxy struct {
		Http    string `yaml:"http_proxy"`
		Https   string `yaml:"https_proxy"`
		NoProxy string `yaml:"no_proxy"`
	}

	// yaml tag for stuff to be copied on volumes
	yamlCopy struct {
		//Once indicates if the copy should be done only on one node matching the targeted labels
		Once bool
		// The volume path where to copy the content
		Path string
		// Labels to restrict the copy to some node sets
		yamlLabel `yaml:",inline"`
		// The list of path patterns identifying content to be copied
		Sources []string `yaml:"sources"`
	}

	// yaml tag for parameters
	yamlParams struct {
		Params map[string]interface{} `yaml:",omitempty"`
	}

	// yaml tag for variables
	yamlVars struct {
		Vars map[string]interface{} `yaml:",omitempty"`
	}

	// yaml tag for authentication parameters
	yamlAuth struct {
		Auth map[string]string `yaml:",omitempty"`
	}

	// yaml tag for environment variables
	yamlEnv struct {
		Env map[string]string `yaml:",omitempty"`
	}

	// yaml tag for labels on nodesets
	yamlLabel struct {
		Labels map[string]string `yaml:",omitempty"`
	}

	// yaml tag for custom playbooks
	yamlPlaybooks struct {
		Playbooks map[string]string `yaml:",omitempty"`
	}

	// yaml tag for component
	yamlComponent struct {
		// The source repository where the component lives
		Repository string
		// The ref (branch or tag) of the component to use
		Ref string
		// The authentication parameters
		yamlAuth `yaml:",inline"`
	}

	// yaml tag for a volume and its parameters
	yamlVolume struct {
		// The mounting path of the created volume
		Path string
		// The parameters required to create the volume (typically provider dependent)
		yamlParams `yaml:",inline"`
	}

	// yaml tag for a shared volume content
	yamlVolumeContent struct {
		// The component holding the content to copy into the volume
		Component string
		// The path of the content to copy
		Path string
	}

	// yaml reference to provider
	yamlProviderRef struct {
		Name string
		// The overriding provider parameters
		yamlParams `yaml:",inline"`
		// The overriding provider environment variables
		yamlEnv `yaml:",inline"`
		// The overriding provider proxy
		Proxy yamlProxy
	}

	// yaml reference to orchestrator
	yamlOrchestratorRef struct {
		// The overriding orchestrator parameters
		yamlParams `yaml:",inline"`
		// The overriding orchestrator environment variables
		yamlEnv `yaml:",inline"`
	}

	// yaml reference to task
	yamlTaskRef struct {
		// The referenced task
		Task string
		// Prefix, optional string used to prefix the stored hook results.*
		Prefix string
		// The overriding parameters
		yamlParams `yaml:",inline"`
		// The overriding environment variables
		yamlEnv `yaml:",inline"`
	}

	//yaml tag for hooks
	yamlHook struct {
		// Hooks to be executed before the corresponding process step
		Before []yamlTaskRef `yaml:",omitempty"`
		// Hooks to be executed after the corresponding process step
		After []yamlTaskRef `yaml:",omitempty"`
	}

	yamlEkara struct {
		// Base for all non-absolute components
		Base string `yaml:",omitempty"`
		// Parent component
		Parent yamlComponent
		// Components declared
		Components map[string]yamlComponent
		// The list of path patterns where to apply the template mechanism
		Templates []string
		// The list of custom playbooks
		yamlPlaybooks `yaml:",inline"`
	}

	yamlNode struct {
		// The number of instances to create within the node set
		Instances int
		// The provider used to create the node set and its settings
		Provider yamlProviderRef
		// The orchestrator settings for this node set
		Orchestrator yamlOrchestratorRef
		// The orchestrator settings for this node set
		Volumes []yamlVolume
		// The Hooks to be executed while creating the node set
		Hooks struct {
			Create  yamlHook `yaml:",omitempty"`
			Destroy yamlHook `yaml:",omitempty"`
		} `yaml:",omitempty"`

		// The labels associated with the nodeset
		yamlLabel `yaml:",inline"`
	}

	// Definition of the Ekara environment
	yamlEnvironment struct {
		// The name of the environment
		Name string
		// The qualifier of the environment
		Qualifier string `yaml:",omitempty"`

		// The description of the environment
		Description string `yaml:",omitempty"`

		// The Ekara platform used to interact with the environment
		Ekara yamlEkara

		// The descriptor variables
		yamlVars `yaml:",inline"`

		// Tasks which can be run on the created environment
		Tasks map[string]struct {
			// Name of the task component
			Component string
			// The task parameters
			yamlParams `yaml:",inline"`
			// The task environment variables
			yamlEnv `yaml:",inline"`
			// The name of the playbook to launch the task
			Playbook string `yaml:",omitempty"`
			// The Hooks to be executed in addition the the main task playbook
			Hooks struct {
				Execute yamlHook `yaml:",omitempty"`
			} `yaml:",omitempty"`
		}

		// Global definition of the orchestrator to install on the environment
		Orchestrator struct {
			// Name of the orchestrator component
			Component string
			// The orchestrator parameters
			yamlParams `yaml:",inline"`
			// The orchestrator environment variables
			yamlEnv `yaml:",inline"`
		}

		// The list of all cloud providers required to create the environment
		Providers map[string]struct {
			// Name of the provider component
			Component string
			// The provider parameters
			yamlParams `yaml:",inline"`
			// The provider environment variables
			yamlEnv `yaml:",inline"`
			// The provider proxy
			Proxy yamlProxy
		}

		// The list of node sets to create
		Nodes map[string]yamlNode

		// Software stacks to be installed on the environment
		Stacks map[string]struct {
			// Name of the stack component
			Component string
			// The name of the stacks on which this one depends
			Dependencies []string `yaml:",omitempty"`
			// The Hooks to be executed while deploying the stack
			Hooks struct {
				Deploy yamlHook `yaml:",omitempty"`
			} `yaml:",omitempty"`

			// The parameters
			yamlParams `yaml:",inline"`
			// The environment variables
			yamlEnv `yaml:",inline"`

			// The stack content to be copied on volumes
			Copies map[string]yamlCopy

			// Custom playbook
			Playbook string
		}

		// Global hooks
		Hooks struct {
			Init    yamlHook `yaml:",omitempty"`
			Create  yamlHook `yaml:",omitempty"`
			Install yamlHook `yaml:",omitempty"`
			Deploy  yamlHook `yaml:",omitempty"`
			Delete  yamlHook `yaml:",omitempty"`
		} `yaml:",omitempty"`

		// Global volumes
		Volumes map[string]struct {
			Content []yamlVolumeContent `yaml:",omitempty"`
		} `yaml:",omitempty"`
	}

	yamlRefs struct {
		Ekara    yamlEkara
		yamlVars `yaml:",inline"`

		Orchestrator struct {
			Component string
		}

		Providers map[string]struct {
			Component string
		}

		Nodes map[string]struct {
			Provider struct {
				Name string
			}
		}

		Stacks map[string]struct {
			Component string
		}

		Tasks map[string]struct {
			Component string
		}
	}
)

// parseYaml parses the url content, template it and unmarshal it into the out struct.
func parseYaml(path string, tplC *TemplateContext, out interface{}) error {
	// Read descriptor content
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Parse just the "vars:" section of the descriptor and fill the template context with it
	err = parseVars(content, tplC)
	if err != nil {
		err = fmt.Errorf("yaml error in %s: %s", path, err.Error())
		return err
	}

	// Template the content of the environment descriptor with the updated template context
	templated, err := tplC.Execute(string(content))
	if err != nil {
		return err
	}

	// Unmarshal the resulting YAML into output structure
	err = yaml.Unmarshal([]byte(templated), out)
	if err != nil {
		err = fmt.Errorf("yaml error in %s : %s", path, err.Error())
		return err
	}

	return nil
}

// parseVars parses the "vars:" section of the descriptor
func parseVars(content []byte, tplC *TemplateContext) error {
	readVars := func(b []byte) (yamlVars, error) {
		yVars := yamlVars{}
		err := yaml.Unmarshal(b, &yVars)
		if err != nil {
			return yVars, err
		}
		return yVars, nil
	}

	// Read raw vars
	vars, err := readVars(content)
	if err != nil {
		return err
	}

	// Marshal variables to be able to template only this part
	onlyVars, err := yaml.Marshal(vars)
	if err != nil {
		return err
	}

	// Apply template
	templatedVars, err := tplC.Execute(string(onlyVars))
	if err != nil {
		return err
	}

	// Read templated vars the templated vars again to merge them into the template context
	vars, err = readVars([]byte(templatedVars))
	if err != nil {
		return err
	}

	if len(vars.Vars) > 0 {
		tplC.addVars(vars.Vars)
	}

	return nil
}

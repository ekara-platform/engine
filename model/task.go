package model

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/GroupePSA/componentizer"
)

type (
	//Task represent an task executable on the built environment
	Task struct {
		// The component containing the task
		cRef componentRef
		// Name of the task
		Name string
		// The playbook to execute
		Playbook string
		// The task parameters
		params Parameters
		// The task environment variables
		envVars EnvVars
		//The hooks linked to the task lifecycle events
		Hooks TaskHooks
	}

	//Tasks represent all the tasks of an environment
	Tasks map[string]Task

	//TaskHooks represents hooks associated to a task
	TaskHooks struct {
		//Execute specifies the hook tasks to run when a task is executed
		Execute Hook
	}

	circularRefTracking map[string]interface{}
)

func (r Task) DescType() string {
	return "Task"
}

func (r Task) DescName() string {
	return r.Name
}

func (r Task) Parameters() Parameters {
	return r.params
}

func (r Task) EnvVars() EnvVars {
	return r.envVars
}

func (r Task) ComponentId() string {
	return r.cRef.ComponentId()
}

func (r Task) Component(model interface{}) (componentizer.Component, error) {
	return r.cRef.Component(model)
}

func (r *Task) merge(with Task) {
	if with.cRef.ref != "" {
		r.cRef = with.cRef
	}
	if with.Playbook != "" {
		r.Playbook = with.Playbook
	}
	r.Hooks.Execute.merge(with.Hooks.Execute)
	r.params = r.params.Override(with.params)
	r.envVars = r.envVars.Override(with.envVars)
}

func createTasks(yamlEnv yamlEnvironment) Tasks {
	res := Tasks{}
	for name, yamlTask := range yamlEnv.Tasks {
		res[name] = Task{
			Name:     name,
			Playbook: yamlTask.Playbook,
			cRef:     componentRef{ref: yamlTask.Component},
			params:   CreateParameters(yamlTask.Params),
			envVars:  CreateEnvVars(yamlTask.Env),
			Hooks: TaskHooks{
				Execute: createHook("execute", yamlTask.Hooks.Execute),
			},
		}
	}
	return res
}

func (r *Tasks) merge(with Tasks) {
	for id, t := range with {
		if task, ok := (*r)[id]; ok {
			task.merge(t)
			(*r)[id] = task
		} else {
			(*r)[id] = t
		}
	}
}

func (r circularRefTracking) String() string {
	b := new(bytes.Buffer)
	for key := range r {
		fmt.Fprintf(b, "%s -> ", key)
	}
	return b.String()
}

//HasTasks returns true if the hook contains at least one task reference
func (r TaskHooks) HasTasks() bool {
	return r.Execute.HasTasks()
}

func (r *TaskHooks) merge(with TaskHooks) {
	r.Execute.merge(with.Execute)
}

func (r Task) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	vErrs := ValidationErrors{}
	if len(r.Playbook) == 0 {
		vErrs.addError(errors.New("no playbook specified"), loc.appendPath("playbook"))
	}
	vErrs.merge(validate(e, loc.appendPath("hooks"), r.Hooks))
	return vErrs
}

func (r TaskHooks) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	vErrs := ValidationErrors{}
	vErrs.merge(validate(e, loc, r.Execute))
	return vErrs
}

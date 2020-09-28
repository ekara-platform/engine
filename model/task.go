package model

import (
	"errors"
	"fmt"
	"github.com/GroupePSA/componentizer"
	"strings"
)

type (
	//Task represent an task executable on the built environment
	Task struct {
		// The component containing the task
		cRef componentRef
		// The component the stack was created from
		selfRef componentRef
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
	if r.isSelfComponent() {
		return r.selfRef.ComponentId()
	} else {
		return r.cRef.ComponentId()
	}
}

func (r Task) Component(model interface{}) (componentizer.Component, error) {
	if r.isSelfComponent() {
		return r.selfRef.Component(model)
	} else {
		return r.cRef.Component(model)
	}
}

func (r Task) isSelfComponent() bool {
	return r.cRef.ref == "" || r.cRef.ref == "_"
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

func createTasks(from component, yamlEnv yamlEnvironment) Tasks {
	res := Tasks{}
	for name, yamlTask := range yamlEnv.Tasks {
		res[name] = Task{
			Name:     name,
			Playbook: yamlTask.Playbook,
			cRef:     componentRef{ref: yamlTask.Component},
			selfRef:  componentRef{ref: from.Id},
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
	builder := strings.Builder{}
	for key := range r {
		builder.WriteString(fmt.Sprintf("%s -> ", key))
	}
	return builder.String()
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
	if r.isSelfComponent() {
		vErrs.merge(validate(e, loc, r.selfRef))
	} else {
		vErrs.merge(validate(e, loc, r.cRef))
	}
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

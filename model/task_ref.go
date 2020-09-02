package model

import (
	"errors"
	"fmt"
)

type (
	//TaskRef represents a reference to a task
	TaskRef struct {
		ref    string
		Prefix string

		// Overriding items
		params  Parameters
		envVars EnvVars
	}
)

//reference return a validatable representation of the reference on a task
func (r TaskRef) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	vErrs := ValidationErrors{}
	if _, ok := e.Tasks[r.ref]; !ok {
		vErrs.addError(fmt.Errorf("no such task: %s", r.ref), loc.appendPath("task"))
	}
	return vErrs
}

func (r *TaskRef) merge(with TaskRef) error {
	if r.ref == "" {
		r.ref = with.ref
	}
	r.Prefix = with.Prefix
	r.params = r.params.Override(with.params)
	r.envVars = r.envVars.Override(with.envVars)
	return nil
}

// Resolve returns a resolved reference to a task containing all the
// inherited content from the referenced task
func (r TaskRef) Resolve(env Environment) (Task, error) {
	task, ok := env.Tasks[r.ref]
	if !ok {
		return Task{}, fmt.Errorf("no such task: %s", r.ref)
	}
	return Task{
		Name:     task.Name,
		cRef:     task.cRef,
		Playbook: task.Playbook,
		Hooks:    task.Hooks,
		params:   task.params.Override(r.params),
		envVars:  task.envVars.Override(r.envVars)}, nil
}

func createTaskRef(tRef yamlTaskRef) TaskRef {
	return TaskRef{
		ref:     tRef.Task,
		Prefix:  tRef.Prefix,
		params:  CreateParameters(tRef.Params),
		envVars: CreateEnvVars(tRef.Env),
	}
}

func checkCircularRefs(taskRefs []TaskRef, alreadyEncountered *circularRefTracking) error {
	for _, taskRef := range taskRefs {
		if _, ok := (*alreadyEncountered)[taskRef.ref]; ok {
			return errors.New("circular task reference: " + alreadyEncountered.String() + taskRef.ref)
		}
	}
	return nil
}

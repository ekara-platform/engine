package model

import (
	"reflect"
)

type (
	// Hook represents tasks to be executed linked to an ekara life cycle event
	Hook struct {
		// Name of the hook
		Name string
		//Before specifies the tasks to run before the ekara life cycle event occurs
		Before []TaskRef
		//After specifies the tasks to run once the ekara life cycle event has occurred
		After []TaskRef
	}
)

func createHook(name string, yamlHook yamlHook) Hook {
	hook := Hook{
		Name:   name,
		Before: make([]TaskRef, len(yamlHook.Before)),
		After:  make([]TaskRef, len(yamlHook.After))}
	for i, yamlRef := range yamlHook.Before {
		hook.Before[i] = createTaskRef(yamlRef)
	}

	for i, yamlRef := range yamlHook.After {
		hook.After[i] = createTaskRef(yamlRef)
	}
	return hook
}

func (r Hook) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	vErrs := ValidationErrors{}
	vErrs.merge(validate(e, loc.appendPath(r.Name).appendPath("before"), r.Before))
	vErrs.merge(validate(e, loc.appendPath(r.Name).appendPath("after"), r.After))
	return vErrs
}

func (r *Hook) merge(with Hook) {
	if !reflect.DeepEqual(r, &with) {
		r.Before = append(r.Before, with.Before...)
		r.After = append(r.After, with.After...)
	}
}

func (r Hook) HasTasks() bool {
	return len(r.Before) > 0 || len(r.After) > 0
}

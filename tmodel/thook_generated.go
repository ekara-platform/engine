package engine

import(
	"github.com/ekara-platform/model"
)

//*****************************************************************************
//
// This file is autogenerated by "go generate .". Do not modify.
//
//*****************************************************************************

// ----------------------------------------------------
// THook is a read only hooks 
// ----------------------------------------------------
type THook interface {
    //After returns the references of tasks to run after the hooked action
    After() []TTaskRef
    //Before returns the references of tasks to run before the hooked action
    Before() []TTaskRef
	
}

// ----------------------------------------------------
// Implementation(s) of THook  
// ----------------------------------------------------

//THookOnHookHolder is the struct containing the Hook in order to implement THook  
type THookOnHookHolder struct {
	h 	model.Hook
}

//CreateTHookForHook returns an holder of Hook implementing THook
func CreateTHookForHook(o model.Hook) THookOnHookHolder {
	return THookOnHookHolder{
		h: o,
	}
}

//After returns the references of tasks to run after the hooked action
func (r THookOnHookHolder) After() []TTaskRef{
	    result := make([]TTaskRef, 0, 0)
    for _ , val := range r.h.After{
        result = append(result, CreateTTaskRefForTaskRef(val))
    }
    return result

}

//Before returns the references of tasks to run before the hooked action
func (r THookOnHookHolder) Before() []TTaskRef{
	    result := make([]TTaskRef, 0, 0)
    for _ , val := range r.h.Before{
        result = append(result, CreateTTaskRefForTaskRef(val))
    }
    return result

}

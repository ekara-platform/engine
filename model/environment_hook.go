package model

type (
	//EnvironmentHooks represents hooks associated to the environment
	EnvironmentHooks struct {
		//Init specifies the hook tasks to run before the environment is created
		Init Hook
		//Create specifies the hook tasks to run when the environment is created
		Create Hook
		//Install specifies the hook tasks to run at the environment installation
		Install Hook
		//Deploy specifies the hook tasks to run at the environment deployment
		Deploy Hook
		//Destroy specifies the hook tasks to run at the environment destruction
		Destroy Hook
	}
)

func createEnvHooks(yamlEnv yamlEnvironment) EnvironmentHooks {
	return EnvironmentHooks{
		Init:    createHook("init", yamlEnv.Hooks.Init),
		Create:  createHook("create", yamlEnv.Hooks.Create),
		Install: createHook("install", yamlEnv.Hooks.Install),
		Deploy:  createHook("deploy", yamlEnv.Hooks.Deploy),
		Destroy: createHook("delete", yamlEnv.Hooks.Delete),
	}
}

//HasTasks returns true if the hook contains at least one task reference
func (r EnvironmentHooks) HasTasks() bool {
	return r.Init.HasTasks() ||
		r.Create.HasTasks() ||
		r.Install.HasTasks() ||
		r.Deploy.HasTasks() ||
		r.Destroy.HasTasks()
}

func (r EnvironmentHooks) validate(e Environment, loc DescriptorLocation) ValidationErrors {
	return validate(e, loc, r.Init, r.Create, r.Install, r.Deploy, r.Destroy)
}

func (r *EnvironmentHooks) merge(with EnvironmentHooks) {
	r.Init.merge(with.Init)
	r.Create.merge(with.Create)
	r.Install.merge(with.Install)
	r.Deploy.merge(with.Deploy)
	r.Destroy.merge(with.Destroy)
}

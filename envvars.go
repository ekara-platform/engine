package engine

// Buffer contains the extra vars to be passed to a playbook
type EnvVars struct {
	content map[string]string
}

func BuildEnvVars() EnvVars {
	r := EnvVars{}
	r.content = make(map[string]string)
	return r
}

func (ev *EnvVars) Add(key, value string) {
	ev.content[key] = value
}

// AddBuffer adds the environment variable coming from the given buffer
func (ev *EnvVars) AddBuffer(b Buffer) {
	for k, v := range b.Envvars {
		ev.content[k] = v
	}
}

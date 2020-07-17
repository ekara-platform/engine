package model

//EnvVarsAware represents the common behavior of environment variables holders
type EnvVarsAware interface {
	EnvVars() EnvVars
}

//EnvVars Represents environment variable
type EnvVars map[string]string

// CreateEmptyParameters builds an empty Parameters structure.
func CreateEmptyEnvVars() EnvVars {
	return make(map[string]string)
}

func CreateEnvVars(src map[string]string) EnvVars {
	dst := make(map[string]string)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (r EnvVars) Override(with map[string]string) EnvVars {
	dst := make(map[string]string)
	for k, v := range r {
		dst[k] = v
	}
	for k, v := range with {
		dst[k] = v
	}
	return dst
}

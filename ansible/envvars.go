package ansible

import (
	"os"
	"runtime"
	"strings"

	"github.com/ekara-platform/model"
)

// EnvVars contains the extra vars to be passed to a playbook
type envVars struct {
	Content map[string]string
}

//CreateEnvVars creates and empy EnvVars
func createEnvVars() envVars {
	r := envVars{}
	r.Content = make(map[string]string)
	return r
}

func (ev envVars) String() string {
	r := ""
	for k, v := range ev.Content {
		r += k + "=" + v + " "
	}
	return strings.Trim(r, " ")
}

// Add adds the given key and value.
//
// If the key already exists then its content will be overwritten by the by the value
func (ev *envVars) add(key, value string) {
	ev.Content[key] = value
}

// AddProxy adds the proxy information if any
//
// If proxy info is already present, it will be overwritten or removed if empty proxy values are passed
func (ev *envVars) addProxy(proxy model.Proxy) {
	if proxy.Http != "" {
		ev.Content["http_proxy"] = proxy.Http
	} else {
		delete(ev.Content, "http_proxy")
	}
	if proxy.Https != "" {
		ev.Content["https_proxy"] = proxy.Https
	} else {
		delete(ev.Content, "https_proxy")
	}
	if proxy.NoProxy != "" {
		ev.Content["no_proxy"] = proxy.NoProxy
	} else {
		delete(ev.Content, "no_proxy")
	}
}

// AddDefaultOsVars adds the HOSTNAME, PATH, TERM and HOME OS variables
func (ev *envVars) addDefaultOsVars() {
	ev.addOsVars("HOSTNAME", "PATH", "TERM", "HOME")
}

// AddOsVars adds the current OS value of the specified variables
//
// If the var list is empty, all os variables are added
func (ev *envVars) addOsVars(vars ...string) {
	osVars := os.Environ()
	for _, osV := range osVars {
		splitOsV := strings.Split(osV, "=")
		if len(vars) == 0 {
			ev.Content[splitOsV[0]] = splitOsV[1]
		} else {
			for _, v := range vars {
				if runtime.GOOS == "windows" {
					if strings.EqualFold(splitOsV[0], v) {
						ev.Content[v] = splitOsV[1]
						break
					}
				} else {
					if splitOsV[0] == v {
						ev.Content[splitOsV[0]] = splitOsV[1]
						break
					}
				}

			}
		}
	}
}

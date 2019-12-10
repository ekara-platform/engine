package ansible

import (
	"encoding/json"
	"strings"

	"github.com/ekara-platform/engine/util"
)

//ExtraVars extra vars, received into the buffer, to be passed to a playbook
type ExtraVars struct {
	Content map[string]interface{}
}

//CreateExtraVars builds the extra var to pass to a playbook
func CreateExtraVars(inputFolder util.FolderPath, outputFolder util.FolderPath) ExtraVars {
	r := ExtraVars{}
	r.Content = make(map[string]interface{}, 0)
	in := strings.Trim(inputFolder.Path(), " ")
	out := strings.Trim(outputFolder.Path(), " ")
	if in != "" {
		r.Content["input_dir"] = in
	}
	if out != "" {
		r.Content["output_dir"] = out
	}
	return r
}

//String returns the extra var content as string
func (ev ExtraVars) String() (string, error) {
	b, e := json.Marshal(ev.Content)
	return string(b), e
}

//Empty returns true if the extra var has no content
func (ev ExtraVars) Empty() bool {
	return len(ev.Content) == 0
}

// Add adds the given key and value.
//
// If the key already exists then its content will be overwritten by the by the value
func (ev *ExtraVars) Add(key, value string) {
	ev.Content[key] = value
}

// AddArray adds the given key and values passed as an array.
//
// If the key already exists then its content will be overwritten by the by the provided values
func (ev *ExtraVars) AddArray(key string, values []string) {
	ev.Content[key] = values
}

// AddMap adds the given key and values passed as an array.
//
// If the key already exists then its content will be overwritten by the by the provided values
func (ev *ExtraVars) AddMap(key string, values map[string]string) {
	ev.Content[key] = values
}

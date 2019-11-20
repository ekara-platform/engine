package ansible

import (
	"strings"

	"github.com/ekara-platform/engine/util"
)

//ExtraVars extra vars, received into the buffer, to be passed to a playbook
type ExtraVars struct {
	Content map[string]string
}

//BuildExtraVars builds the extra var to pass to a playbook
func CreateExtraVars(inputFolder util.FolderPath, outputFolder util.FolderPath) ExtraVars {
	r := ExtraVars{}
	r.Content = make(map[string]string, 0)
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

func (ev ExtraVars) String() string {
	r := ""
	for k, v := range ev.Content {
		r += k + "=" + v + " "
	}
	return strings.Trim(r, " ")
}

func (ev ExtraVars) Empty() bool {
	return len(ev.Content) == 0
}

// Add adds the given key and value.
//
// If the key already exists then its content will be overwritten by the by the value
func (ev *ExtraVars) Add(key, value string) {
	ev.Content[key] = value
}

package ansible

import (
	"strings"

	"github.com/ekara-platform/engine/util"
)

//ExtraVars extra vars, received into the buffer, to be passed to a playbook
type ExtraVars struct {
	// Indicates if the struct has content or not
	Bool bool
	// The extra vars content
	Vals []string
}

func (ev ExtraVars) String() string {
	r := ""
	for _, v := range ev.Vals {
		r = r + " " + v
	}
	return r
}

//BuildExtraVars builds the extra var to pass to a playbook
func BuildExtraVars(extraVars string, inputFolder util.FolderPath, outputFolder util.FolderPath, b Buffer) ExtraVars {
	r := ExtraVars{}
	r.Vals = make([]string, 0)

	vars := strings.Trim(extraVars, " ")
	in := strings.Trim(inputFolder.Path(), " ")
	out := strings.Trim(outputFolder.Path(), " ")
	if vars == "" && in == "" && out == "" && len(b.Extravars) == 0 {
		r.Bool = false
	} else {
		r.Bool = true
		r.Vals = append(r.Vals, "--extra-vars")
		if vars != "" {
			r.Vals = append(r.Vals, vars)
		}

		if in != "" {
			r.Vals = append(r.Vals, "input_dir="+in)
		}

		if out != "" {
			r.Vals = append(r.Vals, "output_dir="+out)
		}

		if len(b.Extravars) > 0 {
			for k, v := range b.Extravars {
				r.Vals = append(r.Vals, k+"="+v)
			}
		}
	}
	return r
}

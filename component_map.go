package engine

import (
	"fmt"

	"github.com/ekara-platform/engine/util"
)

func saveComponentMap(c LaunchContext, dest *util.FolderPath, sr *StepResult) bool {
	e := c.Ekara().ComponentManager().SaveComponentsPaths(c.Log(), *dest)
	if e != nil {
		FailsOnCode(sr, e, fmt.Sprintf("An error occurred saving the components file into :%v", dest.Path()), nil)
		return true
	}
	return false
}

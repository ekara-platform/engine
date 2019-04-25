package engine

import (
	"fmt"
	"log"

	"github.com/ekara-platform/engine/util"
)

func createChildExchangeFolder(parent *util.FolderPath, name string, sr *StepResult, log *log.Logger) (*util.ExchangeFolder, bool) {
	ef, e := parent.AddChildExchangeFolder(name)
	if e != nil {
		err := fmt.Errorf(ErrorAddingExchangeFolder, name, e.Error())
		log.Printf("%s\n", err.Error())
		FailsOnCode(sr, e, err.Error(), nil)
		return ef, true
	}
	e = ef.Create()
	if e != nil {
		err := fmt.Errorf(ErrorCreatingExchangeFolder, name, e.Error())
		log.Printf("%s\n", err.Error())
		FailsOnCode(sr, e, err.Error(), nil)
		return ef, true
	}
	return ef, false
}

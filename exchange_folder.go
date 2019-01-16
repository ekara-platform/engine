package engine

import (
	"fmt"
	"log"

	"github.com/ekara-platform/engine/util"
)

func createChildExchangeFolder(parent *util.FolderPath, name string, sr *StepResult, log *log.Logger) (*util.ExchangeFolder, bool) {
	ef, e := parent.AddChildExchangeFolder(name)
	if e != nil {
		err := fmt.Errorf(ERROR_ADDING_EXCHANGE_FOLDER, name, e.Error())
		log.Printf(err.Error())
		FailsOnCode(sr, e, err.Error(), nil)
		return ef, true
	}
	e = ef.Create()
	if e != nil {
		err := fmt.Errorf(ERROR_CREATING_EXCHANGE_FOLDER, name, e.Error())
		log.Printf(err.Error())
		FailsOnCode(sr, e, err.Error(), nil)
		return ef, true
	}
	return ef, false
}

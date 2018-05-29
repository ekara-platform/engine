package engine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestComponentManager_Fetch(t *testing.T) {
	ctx := createTestContext()
	manager := createComponentManager(&ctx)
	mainPath, e := manager.Fetch("testdata/components/lagoon-platform/core/", "v1.0.1")
	assert.Nil(t, e)
	fmt.Println(mainPath)
	assert.NotNil(t, mainPath)
}

func createTestContext() context {
	os.RemoveAll("testdata/work")
	return context{
		logger:    log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime),
		directory: "testdata/work"}
}

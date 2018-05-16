package engine

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestComponentManager_Fetch(t *testing.T) {
	ctx := createTestContext()
	manager, e := createComponentManager(&ctx)
	assert.Nil(t, e)

	mainPath, e := manager.Fetch("./testdata/components/main", "1.0.1")
	assert.Nil(t, e)
	assert.NotNil(t, mainPath)
}

func createTestContext() context {
	return context{
		logger:  log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime),
		baseDir: "testdata/work"}
}

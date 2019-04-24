package component

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"
)

func TestComponentManager_DirectoriesMatching(t *testing.T) {

	wd, e := os.Getwd()
	assert.Nil(t, e)
	rawManager := filepath.Join(wd, "testdata")
	rawManagerWork := filepath.Join(rawManager, "work")

	manager := buildComponentManager(t, rawManager, rawManagerWork)
	e = manager.Ensure()
	assert.Nil(t, e)
	moduleDirs := manager.MatchingDirectories("modules")
	if assert.Equal(t, len(moduleDirs), 2) {
		assert.Contains(t, moduleDirs, filepath.Join(rawManagerWork, "components", "c1", "modules"))
		assert.Contains(t, moduleDirs, filepath.Join(rawManagerWork, "components", "c3", "modules"))
	}
	inventoryDirs := manager.MatchingDirectories("inventory")
	if assert.Equal(t, len(inventoryDirs), 2) {
		assert.Contains(t, inventoryDirs, filepath.Join(rawManagerWork, "components", "c2", "inventory"))
		assert.Contains(t, inventoryDirs, filepath.Join(rawManagerWork, "components", "c3", "inventory"))
	}
}

func buildComponentManager(t *testing.T, rawManager, rawManagerWork string) ComponentManager {

	os.RemoveAll(rawManagerWork)

	manager := CreateComponentManager(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), map[string]interface{}{}, rawManagerWork)
	registerComponent(t, manager, rawManager, "c1")
	registerComponent(t, manager, rawManager, "c2")
	registerComponent(t, manager, rawManager, "c3")
	registerComponent(t, manager, rawManager, "c4")
	return manager
}

func registerComponent(t *testing.T, manager ComponentManager, rawManagerWork string, id string) {
	rawbase := filepath.Join(rawManagerWork, "components")
	b, e := model.CreateBase(rawbase)
	assert.Nil(t, e)
	assert.Equal(t, b.Url.UpperScheme(), model.SchemeFile)
	assert.True(t, len(b.Url.AsFilePath()) > 0)
	r, e := model.CreateRepository(b, "ekara-platform/"+id, "v1.0.0", "")
	assert.Nil(t, e)
	assert.Equal(t, r.Url.UpperScheme(), model.SchemeFile)
	assert.True(t, len(r.Url.AsFilePath()) > 0)
	component := model.CreateComponent(id, r)
	manager.RegisterComponent(component)
}

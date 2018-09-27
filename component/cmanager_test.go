package component

import (
	"github.com/lagoon-platform/model"
	"github.com/stretchr/testify/assert"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestComponentManager_DirectoriesMatching(t *testing.T) {
	manager := buildComponentManager(t)
	e := manager.Ensure()
	assert.Nil(t, e)
	moduleDirs := manager.MatchingDirectories("modules")
	assert.Contains(t, moduleDirs, "testdata/work/components/c1/modules")
	assert.Contains(t, moduleDirs, "testdata/work/components/c3/modules")
	inventoryDirs := manager.MatchingDirectories("inventory")
	assert.Contains(t, inventoryDirs, "testdata/work/components/c2/inventory")
	assert.Contains(t, inventoryDirs, "testdata/work/components/c3/inventory")
}

func buildComponentManager(t *testing.T) ComponentManager {
	os.RemoveAll("testdata/work")
	manager := CreateComponentManager(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), &model.Environment{}, "testdata/work")
	wd, e := os.Getwd()
	assert.Nil(t, e)
	base, e := url.Parse(filepath.Join(wd, "testdata", "components") + "/")
	assert.Nil(t, e)
	registerComponent(t, manager, base, "c1")
	registerComponent(t, manager, base, "c2")
	registerComponent(t, manager, base, "c3")
	registerComponent(t, manager, base, "c4")
	return manager
}

func registerComponent(t *testing.T, manager ComponentManager, base *url.URL, id string) {
	component, e := model.CreateComponent(base, id, "lagoon-platform/"+id, "1.0.0")
	assert.Nil(t, e)
	manager.RegisterComponent(component)
}

package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManagerInitialGrossContent(t *testing.T) {

	am := CreateActionManager()

	assert.NotNil(t, am.actions)

	// Check actions preloaded into the manager
	assert.False(t, am.empty())
	assert.Equal(t, len(am.actions), 6)

	v, err := am.get(ActionFailId)
	assert.Nil(t, err)
	check(t, v, ActionFailId, ActionNilId, "FailOnError")

	v, err = am.get(ActionReportId)
	assert.Nil(t, err)
	check(t, v, ActionReportId, ActionFailId, "Report")

	v, err = am.get(ActionCreateId)
	assert.Nil(t, err)
	check(t, v, ActionCreateId, ActionReportId, "Create")

	v, err = am.get(ActionInstallId)
	assert.Nil(t, err)
	check(t, v, ActionInstallId, ActionCreateId, "Install")

	v, err = am.get(ActionDeployId)
	assert.Nil(t, err)
	check(t, v, ActionDeployId, ActionInstallId, "Deploy")

	v, err = am.get(ActionCheckId)
	assert.Nil(t, err)
	check(t, v, ActionCheckId, ActionNilId, "Check")
	// The nil action shouldn't be strored into the manager
	_, err = am.get(ActionNilId)
	assert.NotNil(t, err)
}

func check(t *testing.T, a action, id ActionId, depends ActionId, name string) {
	assert.Equal(t, a.id, id)
	assert.Equal(t, a.dependsOn, depends)
	assert.Equal(t, a.name, name)
}

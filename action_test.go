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
	assert.Equal(t, len(am.actions), 7)

	v, err := am.get(ActionFailID)
	assert.Nil(t, err)
	check(t, v, ActionFailID, ActionNilID, "FailOnError")

	v, err = am.get(ActionReportID)
	assert.Nil(t, err)
	check(t, v, ActionReportID, ActionFailID, "Report")

	v, err = am.get(ActionCreateID)
	assert.Nil(t, err)
	check(t, v, ActionCreateID, ActionReportID, "Create")

	v, err = am.get(ActionInstallID)
	assert.Nil(t, err)
	check(t, v, ActionInstallID, ActionCreateID, "Install")

	v, err = am.get(ActionDeployID)
	assert.Nil(t, err)
	check(t, v, ActionDeployID, ActionInstallID, "Deploy")

	v, err = am.get(ActionCheckID)
	assert.Nil(t, err)
	check(t, v, ActionCheckID, ActionNilID, "Check")

	v, err = am.get(ActionDumpID)
	assert.Nil(t, err)
	check(t, v, ActionDumpID, ActionCheckID, "Dump")

	// The nil action shouldn't be strored into the manager
	_, err = am.get(ActionNilID)
	assert.NotNil(t, err)
}

func check(t *testing.T, a Action, id ActionID, depends ActionID, name string) {
	assert.Equal(t, a.id, id)
	assert.Equal(t, a.dependsOn, depends)
	assert.Equal(t, a.name, name)
}

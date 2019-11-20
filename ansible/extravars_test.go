package ansible

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/stretchr/testify/assert"
)

func TestNoExtraVars(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath(""), util.CreateFolderPath(""))
	assert.Equal(t, false, !ev.Empty())
}

func TestExtraVarsString(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath(""), util.CreateFolderPath(""))
	ev.Add("aa", "bb")
	assert.Equal(t, true, !ev.Empty())
	assert.Equal(t, 1, len(ev.Content))
	assert.Equal(t, ev.Content["aa"], "bb")
}

func TestExtraVarsInputFolder(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath("aa/bb"), util.CreateFolderPath(""))
	assert.Equal(t, true, !ev.Empty())
	assert.Equal(t, 1, len(ev.Content))
	assert.Equal(t, ev.Content["input_dir"], "aa/bb")
}

func TestExtraVarsOutputFolder(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath(""), util.CreateFolderPath("aa/bb"))
	assert.Equal(t, true, !ev.Empty())
	assert.Equal(t, 1, len(ev.Content))
	assert.Equal(t, ev.Content["output_dir"], "aa/bb")
}

func TestExtraVarsInputOutputFolder(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath("aa/bb"), util.CreateFolderPath("aa/bb"))
	assert.Equal(t, true, !ev.Empty())
	assert.Equal(t, 2, len(ev.Content))
	assert.Equal(t, ev.Content["input_dir"], "aa/bb")
	assert.Equal(t, ev.Content["output_dir"], "aa/bb")
}

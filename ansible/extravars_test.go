package ansible

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/stretchr/testify/assert"
)

func TestNoExtraVars(t *testing.T) {
	ev := BuildExtraVars("", util.CreateFolderPath(""), util.CreateFolderPath(""), Buffer{})
	assert.Equal(t, false, ev.Bool)
}

func TestExtraVarsString(t *testing.T) {
	ev := BuildExtraVars("aa=bb", util.CreateFolderPath(""), util.CreateFolderPath(""), Buffer{})
	assert.Equal(t, true, ev.Bool)
	assert.Equal(t, 2, len(ev.Vals))
	assert.Equal(t, "--extra-vars", ev.Vals[0])
	assert.Equal(t, "aa=bb", ev.Vals[1])
}

func TestExtraVarsInputFolder(t *testing.T) {
	ev := BuildExtraVars("", util.CreateFolderPath("aa/bb"), util.CreateFolderPath(""), Buffer{})
	assert.Equal(t, true, ev.Bool)
	assert.Equal(t, 2, len(ev.Vals))
	assert.Equal(t, "--extra-vars", ev.Vals[0])
	assert.Equal(t, "input_dir=aa/bb", ev.Vals[1])
}

func TestExtraVarsOutputFolder(t *testing.T) {
	ev := BuildExtraVars("", util.CreateFolderPath(""), util.CreateFolderPath("aa/bb"), Buffer{})
	assert.Equal(t, true, ev.Bool)
	assert.Equal(t, 2, len(ev.Vals))
	assert.Equal(t, "--extra-vars", ev.Vals[0])
	assert.Equal(t, "output_dir=aa/bb", ev.Vals[1])
}

func TestExtraVarsInputOutputFolder(t *testing.T) {
	ev := BuildExtraVars("", util.CreateFolderPath("aa/bb"), util.CreateFolderPath("aa/bb"), Buffer{})
	assert.Equal(t, true, ev.Bool)
	assert.Equal(t, 3, len(ev.Vals))
	assert.Equal(t, "--extra-vars", ev.Vals[0])
	assert.Equal(t, "input_dir=aa/bb", ev.Vals[1])
	assert.Equal(t, "output_dir=aa/bb", ev.Vals[2])
}

func TestExtraVarsBuffer(t *testing.T) {
	b := Buffer{}
	extraVars := make(map[string]string)
	extraVars["key1"] = "val1"
	b.Extravars = extraVars
	ev := BuildExtraVars("", util.CreateFolderPath(""), util.CreateFolderPath(""), b)
	assert.Equal(t, true, ev.Bool)
	assert.Equal(t, 2, len(ev.Vals))
	assert.Equal(t, "--extra-vars", ev.Vals[0])
	assert.Equal(t, "key1=val1", ev.Vals[1])
}

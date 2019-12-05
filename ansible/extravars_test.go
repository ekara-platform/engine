package ansible

import (
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/stretchr/testify/assert"
)

func TestNoExtraVars(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath(""), util.CreateFolderPath(""))
	assert.Equal(t, false, !ev.Empty())
	s, e := ev.String()
	assert.Nil(t, e)
	assert.Equal(t, s, "{}")
}

func TestExtraVarsString(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath(""), util.CreateFolderPath(""))
	ev.Add("aa", "bb")
	assert.Equal(t, true, !ev.Empty())
	assert.Equal(t, 1, len(ev.Content))
	s, e := ev.String()
	assert.Nil(t, e)
	assert.Equal(t, s, "{\"aa\":\"bb\"}")
}

func TestExtraVarsInputFolder(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath("aa/bb"), util.CreateFolderPath(""))
	assert.Equal(t, true, !ev.Empty())
	assert.Equal(t, 1, len(ev.Content))
	s, e := ev.String()
	assert.Nil(t, e)
	assert.Equal(t, s, "{\"input_dir\":\"aa/bb\"}")
}

func TestExtraVarsOutputFolder(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath(""), util.CreateFolderPath("aa/bb"))
	assert.Equal(t, true, !ev.Empty())
	assert.Equal(t, 1, len(ev.Content))
	s, e := ev.String()
	assert.Nil(t, e)
	assert.Equal(t, s, "{\"output_dir\":\"aa/bb\"}")
}

func TestExtraVarsInputOutputFolder(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath("aa/bb"), util.CreateFolderPath("aa/bb"))
	assert.Equal(t, true, !ev.Empty())
	assert.Equal(t, 2, len(ev.Content))
	s, e := ev.String()
	assert.Nil(t, e)
	assert.Equal(t, s, "{\"input_dir\":\"aa/bb\",\"output_dir\":\"aa/bb\"}")
}

func TestExtraVarsEmptryArray(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath(""), util.CreateFolderPath(""))
	values := []string{}
	ev.AddArray("aa", values)
	assert.Equal(t, true, !ev.Empty())
	assert.Equal(t, 1, len(ev.Content))
	s, e := ev.String()
	assert.Nil(t, e)
	assert.Equal(t, s, "{\"aa\":[]}")
}

func TestExtraVarsArrayOneValue(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath(""), util.CreateFolderPath(""))
	values := []string{"str1"}
	ev.AddArray("aa", values)
	assert.Equal(t, true, !ev.Empty())
	assert.Equal(t, 1, len(ev.Content))
	s, e := ev.String()
	assert.Nil(t, e)
	assert.Equal(t, s, "{\"aa\":[\"str1\"]}")
}

func TestExtraVarsArrayValues(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath(""), util.CreateFolderPath(""))
	values := []string{"str1", "str2", "str3"}
	ev.AddArray("aa", values)
	assert.Equal(t, true, !ev.Empty())
	assert.Equal(t, 1, len(ev.Content))
	s, e := ev.String()
	assert.Nil(t, e)
	assert.Equal(t, s, "{\"aa\":[\"str1\",\"str2\",\"str3\"]}")
}

func TestExtraVarsEmptryMap(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath(""), util.CreateFolderPath(""))
	values := map[string]string{}

	ev.AddMap("aa", values)
	assert.Equal(t, true, !ev.Empty())
	assert.Equal(t, 1, len(ev.Content))
	s, e := ev.String()
	assert.Nil(t, e)
	assert.Equal(t, s, "{\"aa\":{}}")

}

func TestExtraVarsMapOneValue(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath(""), util.CreateFolderPath(""))
	values := map[string]string{
		"key1": "val1",
	}
	ev.AddMap("aa", values)
	assert.Equal(t, true, !ev.Empty())
	assert.Equal(t, 1, len(ev.Content))
	s, e := ev.String()
	assert.Nil(t, e)
	assert.Equal(t, s, "{\"aa\":{\"key1\":\"val1\"}}")
}

func TestExtraVarsMapValues(t *testing.T) {
	ev := CreateExtraVars(util.CreateFolderPath(""), util.CreateFolderPath(""))
	values := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
	}

	ev.AddMap("aa", values)
	assert.Equal(t, true, !ev.Empty())
	assert.Equal(t, 1, len(ev.Content))
	s, e := ev.String()
	assert.Nil(t, e)
	assert.Equal(t, s, "{\"aa\":{\"key1\":\"val1\",\"key2\":\"val2\",\"key3\":\"val3\"}}")

}

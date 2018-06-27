package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExchangeFolder(t *testing.T) {
	ef, err := CreateExchangeFolder(".", "testExFolder")
	assert.Nil(t, err)

	ok, err := ef.Location.Exixts()
	assert.False(t, ok)

	ok, err = ef.Input.Exixts()
	assert.False(t, ok)

	ok, err = ef.Output.Exixts()
	assert.False(t, ok)

	// Physical creation
	err = ef.Create()
	assert.Nil(t, err)

	ok, err = ef.Location.Exixts()
	assert.True(t, ok)

	ok, err = ef.Input.Exixts()
	assert.True(t, ok)

	ok, err = ef.Output.Exixts()
	assert.True(t, ok)

	// Deletion
	err = ef.Delete()
	assert.Nil(t, err)

	ok, err = ef.Location.Exixts()
	assert.False(t, ok)

	ok, err = ef.Input.Exixts()
	assert.False(t, ok)

	ok, err = ef.Output.Exixts()
	assert.False(t, ok)
}

func TestChildren(t *testing.T) {
	ef, err := CreateExchangeFolder(".", "testExFolder")
	assert.Nil(t, err)

	assert.Equal(t, len(ef.Location.Children), 0)
	assert.Equal(t, len(ef.Input.Children), 0)
	assert.Equal(t, len(ef.Output.Children), 0)

	child, err := ef.Input.AddChildExchangeFolder("child")
	assert.Nil(t, err)

	assert.Equal(t, len(ef.Location.Children), 0)
	assert.Equal(t, len(ef.Input.Children), 1)
	assert.Equal(t, len(ef.Output.Children), 0)

	_, ok := (ef.Input.Children)["child"]
	assert.True(t, ok)

	ok, err = child.Location.Exixts()
	assert.False(t, ok)

	ok, err = child.Input.Exixts()
	assert.False(t, ok)

	ok, err = child.Output.Exixts()
	assert.False(t, ok)

	// Physical creation of the parent
	// This is supposed to also create the children
	err = ef.Create()
	assert.Nil(t, err)

	ok, err = child.Location.Exixts()
	assert.True(t, ok)

	ok, err = child.Input.Exixts()
	assert.True(t, ok)

	ok, err = child.Output.Exixts()
	assert.True(t, ok)

	// Deletion of the child only
	err = child.Delete()
	assert.Nil(t, err)

	ok, err = child.Location.Exixts()
	assert.False(t, ok)

	ok, err = child.Input.Exixts()
	assert.False(t, ok)

	ok, err = child.Output.Exixts()
	assert.False(t, ok)

	// we just recreate the child
	err = ef.Create()

	ok, err = child.Location.Exixts()
	assert.True(t, ok)

	ok, err = child.Input.Exixts()
	assert.True(t, ok)

	ok, err = child.Output.Exixts()
	assert.True(t, ok)

	// Physical creation of the parent
	// This is supposed to also delete the children
	err = ef.Delete()
	assert.Nil(t, err)

	ok, err = child.Location.Exixts()
	assert.False(t, ok)

	ok, err = child.Input.Exixts()
	assert.False(t, ok)

	ok, err = child.Output.Exixts()
	assert.False(t, ok)

}

func TestContent(t *testing.T) {
	ef, err := CreateExchangeFolder(".", "testExFolder")
	assert.Nil(t, err)
	defer ef.Delete()

	b := []byte("Dummy file Content \n")
	err = ef.Input.Write(b, "test.txt")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "The directory is not created")

	ef.Create()
	err = ef.Input.Write(b, "test.txt")
	assert.Nil(t, err)

	ok := ef.Input.Contains("test.txt")
	assert.True(t, ok)

	ok = ef.Output.Contains("test.txt")
	assert.False(t, ok)

	err = ef.Input.Copy("test.txt", *ef.Output)
	assert.Nil(t, err)

	ok = ef.Output.Contains("test.txt")
	assert.True(t, ok)

	err = ef.CleanAll()
	assert.Nil(t, err)

	ok = ef.Input.Contains("test.txt")
	assert.False(t, ok)

	ok = ef.Output.Contains("test.txt")
	assert.False(t, ok)
}

package engine_test

import (
	"log"
	"os"
	"testing"

	"github.com/lagoon-platform/engine"
	"github.com/stretchr/testify/assert"
)

func TestHttpGetDescriptor(t *testing.T) {
	_, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "https://raw.githubusercontent.com/lagoon-platform/engine/master/testdata/complete_descriptor.yaml")
	//_, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/complete_descriptor.yaml")

	// no error occurred
	assert.Nil(t, e)
}

func TestHttpGetNoDescriptor(t *testing.T) {
	_, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "https://raw.githubusercontent.com/lagoon-platform/engine/master/testdata/DUMMY.yaml")

	// an error occurred
	assert.NotNil(t, e)

	// the error code should be 404
	assert.Equal(t, "HTTP Error getting the environment descriptor , error code 404", e.Error())
}

package engine_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/lagoon-platform/engine"
	"github.com/stretchr/testify/assert"
)

func TestNoProviders(t *testing.T) {
	testEmpyContent(t, "providers")
}

func TestNodes(t *testing.T) {
	testEmpyContent(t, "nodes")
}

func TestStacks(t *testing.T) {
	testEmpyContent(t, "stacks")
}

func testEmpyContent(t *testing.T, name string) {
	file := fmt.Sprintf("./testdata/grammar/no_%s.yaml", name)
	_, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), file)
	assert.NotNil(t, e)

	ges, err := engine.FromError(e)
	assert.Nil(t, err)
	assert.NotNil(t, ges)
	assert.Equal(t, 1, len(ges.Errors))
	assert.Equal(t, name, ges.Errors[0].Location)
	assert.Equal(t, "Can't be empty", ges.Errors[0].Message)
	assert.Equal(t, engine.Error, ges.Errors[0].ErrorType)
}

func TestNoNodesProvidersStacks(t *testing.T) {
	_, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/no_nodes_stacks_providers.yaml")
	assert.NotNil(t, e)

	ges, err := engine.FromError(e)
	assert.Nil(t, err)
	assert.NotNil(t, ges)
	assert.Equal(t, 3, len(ges.Errors))
	testEmptyErrorsContent(t, ges, "providers", "nodes", "stacks")
}

func testEmptyErrorsContent(t *testing.T, ges engine.GrammarErrors, names ...string) {
	for i, name := range names {
		testEmptyErrorContent(t, ges.Errors[i], name)
	}
}

func testEmptyErrorContent(t *testing.T, ge engine.GrammarError, name string) {
	assert.Equal(t, name, ge.Location)
	assert.Equal(t, "Can't be empty", ge.Message)
	assert.Equal(t, engine.Error, ge.ErrorType)
}

func TestNoNodesInstance(t *testing.T) {
	_, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/no_nodes_instance.yaml")
	assert.NotNil(t, e)

	ges, err := engine.FromError(e)
	assert.Nil(t, err)
	assert.NotNil(t, ges)
	assert.Equal(t, 1, len(ges.Errors))
	assert.Equal(t, "nodes.managers.instances", ges.Errors[0].Location)
	assert.Equal(t, "Must be greather than 0", ges.Errors[0].Message)
	assert.Equal(t, engine.Error, ges.Errors[0].ErrorType)
}

func TestNoNodesProvider(t *testing.T) {
	_, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/no_nodes_provider.yaml")
	assert.NotNil(t, e)

	ges, err := engine.FromError(e)
	assert.Nil(t, err)
	assert.NotNil(t, ges)
	assert.Equal(t, 1, len(ges.Errors))
	assert.Equal(t, "nodes.managers.provider", ges.Errors[0].Location)
	assert.Equal(t, "The provider is required", ges.Errors[0].Message)
	assert.Equal(t, engine.Error, ges.Errors[0].ErrorType)
}

func TestNodesUnknownProvider(t *testing.T) {
	_, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/nodes_unknown_provider.yaml")
	assert.NotNil(t, e)

	ges, err := engine.FromError(e)
	assert.Nil(t, err)
	assert.NotNil(t, ges)
	assert.Equal(t, 1, len(ges.Errors))
	assert.Equal(t, "nodes.managers.provider.name", ges.Errors[0].Location)
	assert.Equal(t, "The provider is unknown", ges.Errors[0].Message)
	assert.Equal(t, engine.Error, ges.Errors[0].ErrorType)
}

func TestStacksNoDeployOn(t *testing.T) {
	_, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/stacks_no_deploy_on.yaml")
	assert.NotNil(t, e)

	ges, err := engine.FromError(e)
	assert.Nil(t, err)
	assert.NotNil(t, ges)
	assert.Equal(t, 1, len(ges.Errors))
	assert.Equal(t, "stacks.monitoring.deployOn", ges.Errors[0].Location)
	assert.Equal(t, "Is empty then this stack won't be deployed", ges.Errors[0].Message)
	assert.Equal(t, engine.Error, ges.Errors[0].ErrorType)
}

func TestStacksUnknownDeployOn(t *testing.T) {
	_, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/stacks_unknow_deploy_on.yaml")
	assert.NotNil(t, e)

	ges, err := engine.FromError(e)
	assert.Nil(t, err)
	assert.NotNil(t, ges)
	assert.Equal(t, 1, len(ges.Errors))
	assert.Equal(t, "stacks.monitoring.deployOn:DUMMY", ges.Errors[0].Location)
	assert.Equal(t, "The label is unknown", ges.Errors[0].Message)
	assert.Equal(t, engine.Error, ges.Errors[0].ErrorType)
}

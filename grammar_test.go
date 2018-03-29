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
	_, ges, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), file)
	assert.Nil(t, e)

	assert.NotNil(t, ges)
	assert.Equal(t, true, ges.HasErrors())
	assert.Equal(t, false, ges.HasWarnings())
	assert.Equal(t, 1, len(ges.Errors))
	assert.Equal(t, name, ges.Errors[0].Location)
	assert.Equal(t, "Can't be empty", ges.Errors[0].Message)
	assert.Equal(t, engine.Error, ges.Errors[0].ErrorType)
}

func TestNoNodesProvidersStacks(t *testing.T) {
	_, ges, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/no_nodes_stacks_providers.yaml")
	assert.Nil(t, e)

	assert.NotNil(t, ges)
	assert.Equal(t, true, ges.HasErrors())
	assert.Equal(t, false, ges.HasWarnings())
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
	_, ges, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/no_nodes_instance.yaml")
	assert.Nil(t, e)

	assert.NotNil(t, ges)
	assert.Equal(t, true, ges.HasErrors())
	assert.Equal(t, false, ges.HasWarnings())
	assert.Equal(t, 1, len(ges.Errors))
	assert.Equal(t, "nodes.managers.instances", ges.Errors[0].Location)
	assert.Equal(t, "Must be greather than 0", ges.Errors[0].Message)
	assert.Equal(t, engine.Error, ges.Errors[0].ErrorType)
}

func TestNoNodesProvider(t *testing.T) {
	_, ges, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/no_nodes_provider.yaml")
	assert.Nil(t, e)

	assert.NotNil(t, ges)
	assert.Equal(t, true, ges.HasErrors())
	assert.Equal(t, false, ges.HasWarnings())
	assert.Equal(t, 1, len(ges.Errors))
	assert.Equal(t, "nodes.managers.provider", ges.Errors[0].Location)
	assert.Equal(t, "The provider is required", ges.Errors[0].Message)
	assert.Equal(t, engine.Error, ges.Errors[0].ErrorType)
}

func TestNodesUnknownProvider(t *testing.T) {
	_, ges, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/nodes_unknown_provider.yaml")
	assert.Nil(t, e)

	assert.NotNil(t, ges)
	assert.Equal(t, true, ges.HasErrors())
	assert.Equal(t, false, ges.HasWarnings())

	assert.Equal(t, 1, len(ges.Errors))
	assert.Equal(t, "nodes.managers.provider.name", ges.Errors[0].Location)
	assert.Equal(t, "The provider is unknown", ges.Errors[0].Message)
	assert.Equal(t, engine.Error, ges.Errors[0].ErrorType)
}

func TestNodesUnknownHook(t *testing.T) {
	_, ges, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/nodes_unknown_hook.yaml")
	assert.Nil(t, e)

	assert.NotNil(t, ges)
	assert.Equal(t, true, ges.HasErrors())
	assert.Equal(t, false, ges.HasWarnings())
	assert.Equal(t, 2, len(ges.Errors))
	testHook(t, "nodes.managers.hooks.provision.before:DUMMY", 0, ges)
	testHook(t, "nodes.managers.hooks.provision.after:DUMMY", 1, ges)
}

func TestStacksNoDeployOnError(t *testing.T) {
	_, ges, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/stacks_no_deploy_on_error.yaml")
	assert.Nil(t, e)

	assert.NotNil(t, ges)
	assert.Equal(t, true, ges.HasErrors())
	assert.Equal(t, false, ges.HasWarnings())
	assert.Equal(t, 1, len(ges.Errors))
	assert.Equal(t, "stacks.monitoring.deployOn", ges.Errors[0].Location)
	assert.Equal(t, "Is empty then this stack won't be deployed", ges.Errors[0].Message)
	assert.Equal(t, engine.Error, ges.Errors[0].ErrorType)
}

func TestStacksNoDeployOnWarning(t *testing.T) {
	_, ges, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/stacks_no_deploy_on_warning.yaml")
	assert.Nil(t, e)

	assert.NotNil(t, ges)
	assert.Equal(t, false, ges.HasErrors())
	assert.Equal(t, true, ges.HasWarnings())
	assert.Equal(t, 1, len(ges.Errors))
	assert.Equal(t, "stacks.monitoring.deployOn", ges.Errors[0].Location)
	assert.Equal(t, "Is empty then this stack won't be deployed", ges.Errors[0].Message)
	assert.Equal(t, engine.Warning, ges.Errors[0].ErrorType)
}

func TestStacksUnknownDeployOn(t *testing.T) {
	_, ges, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/stacks_unknow_deploy_on.yaml")
	assert.Nil(t, e)

	assert.NotNil(t, ges)
	assert.Equal(t, true, ges.HasErrors())
	assert.Equal(t, false, ges.HasWarnings())
	assert.Equal(t, 1, len(ges.Errors))
	assert.Equal(t, "stacks.monitoring.deployOn:DUMMY", ges.Errors[0].Location)
	assert.Equal(t, "The label is unknown", ges.Errors[0].Message)
	assert.Equal(t, engine.Error, ges.Errors[0].ErrorType)
}

func TestTasksNoPlayBook(t *testing.T) {
	_, ges, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/no_task_playbook.yaml")
	assert.Nil(t, e)

	assert.NotNil(t, ges)
	assert.Equal(t, true, ges.HasErrors())
	assert.Equal(t, false, ges.HasWarnings())
	assert.Equal(t, 1, len(ges.Errors))
	assert.Equal(t, "tasks.task1.task", ges.Errors[0].Location)
	assert.Equal(t, "The playbook is required", ges.Errors[0].Message)
	assert.Equal(t, engine.Error, ges.Errors[0].ErrorType)
}

func TestTasksUnknownRunOn(t *testing.T) {
	_, ges, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/tasks_unknown_run_on.yaml")
	assert.Nil(t, e)

	assert.NotNil(t, ges)
	assert.Equal(t, true, ges.HasErrors())
	assert.Equal(t, false, ges.HasWarnings())
	assert.Equal(t, 1, len(ges.Errors))
	assert.Equal(t, "tasks.task1.runOn:DUMMY", ges.Errors[0].Location)
	assert.Equal(t, "The label is unknown", ges.Errors[0].Message)
	assert.Equal(t, engine.Error, ges.Errors[0].ErrorType)
}

func TestUnknownGlobalHooks(t *testing.T) {
	_, ges, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "./testdata/grammar/unknown_global_hook.yaml")
	assert.Nil(t, e)

	assert.NotNil(t, ges)
	assert.Equal(t, true, ges.HasErrors())
	assert.Equal(t, false, ges.HasWarnings())
	assert.Equal(t, 10, len(ges.Errors))

	testHook(t, "hooks.init.before:DUMMY", 0, ges)
	testHook(t, "hooks.init.after:DUMMY", 1, ges)
	testHook(t, "hooks.provision.before:DUMMY", 2, ges)
	testHook(t, "hooks.provision.after:DUMMY", 3, ges)
	testHook(t, "hooks.deploy.before:DUMMY", 4, ges)
	testHook(t, "hooks.deploy.after:DUMMY", 5, ges)
	testHook(t, "hooks.undeploy.before:DUMMY", 6, ges)
	testHook(t, "hooks.undeploy.after:DUMMY", 7, ges)
	testHook(t, "hooks.destroy.before:DUMMY", 8, ges)
	testHook(t, "hooks.destroy.after:DUMMY", 9, ges)

}

func testHook(t *testing.T, msg string, index int, ges engine.GrammarErrors) {
	assert.Equal(t, msg, ges.Errors[index].Location)
	assert.Equal(t, "The hook task is unknown", ges.Errors[index].Message)
	assert.Equal(t, engine.Error, ges.Errors[index].ErrorType)
}

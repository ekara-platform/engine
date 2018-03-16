package engine_test

import (
	"log"
	"os"
	"testing"

	"github.com/lagoon-platform/engine"
	"github.com/stretchr/testify/assert"
)

func TestCreateEngine(t *testing.T) {
	lagoon, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "testdata/lagoon.yaml")
	assert.Nil(t, e) // no error occurred

	env := lagoon.GetEnvironment()
	assert.Equal(t, "testEnvironment", env.GetName())                               // importing file have has precedence
	assert.Equal(t, "This is my awesome Lagoon environment.", env.GetDescription()) // imported files are merged
	assert.Equal(t, []string{"tag1", "tag2"}, env.GetLabels().AsStrings())
	// FIXME assert.Contains(t, "task1", "task2", "task3", env.Hooks.Provision.After)        // order matters
}

func TestCreateEngineComplete(t *testing.T) {
	lagoon, e := engine.Create(log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime), "testdata/complete_descriptor.yaml")
	assert.Nil(t, e)

	env := lagoon.GetEnvironment()
	assert.Equal(t, "name_value", env.GetName())
	assert.Equal(t, "description_value", env.GetDescription())
	assert.Equal(t, "baselocation_value", env.GetBaseLocation())

	// Environment Version
	v, e := env.GetVersion()
	assert.Nil(t, e)
	assert.Equal(t, 1, v.Major())
	assert.Equal(t, 2, v.Minor())
	assert.Equal(t, 3, v.Micro())
	assert.Equal(t, "1.2.3", v.AsString())

	// Environment Labels
	labels := env.GetLabels()
	assert.NotNil(t, labels)
	assert.Equal(t, 3, len(labels.AsStrings()))
	assert.Equal(t, true, labels.Contains("root_label1", "root_label2", "root_label3"))

	//------------------------------------------------------------
	// Lagoon Plateform
	//------------------------------------------------------------
	pla := env.GetLagoonPlatform()
	assert.NotNil(t, pla)
	assert.Equal(t, "version_value", pla.GetVersion())
	assert.Equal(t, "registry_value", pla.GetRegistry())

	// Lagoon Plateform Proxy
	proxy := pla.GetProxy()
	assert.NotNil(t, proxy)
	assert.Equal(t, "proxy_http_value", proxy.GetHttp())
	assert.Equal(t, "proxy_https_value", proxy.GetHttps())
	assert.Equal(t, "proxy_noproxy_value", proxy.GetNoProxy())

	//------------------------------------------------------------
	// Environment Providers
	//------------------------------------------------------------
	providerDescs := env.GetProviderDescriptions()
	assert.NotNil(t, providerDescs)

	b := providerDescs.Contains("aws")
	assert.Equal(t, true, b)
	b = providerDescs.Contains("aws", "azure")
	assert.Equal(t, true, b)

	// Environment Provider
	p, _ := providerDescs.GetProvider("aws")
	assert.NotNil(t, p)

	// Environment Provider Labels
	labels = p.GetLabels()
	assert.NotNil(t, labels)
	assert.Equal(t, 2, len(labels.AsStrings()))
	assert.Equal(t, true, labels.Contains("aws_tag1_value", "aws_tag2_value"))

	// Environment Provider Parameters
	params := p.GetParameters().AsMap()
	assert.NotNil(t, params)
	param := params["aws_param_key1"]
	assert.NotNil(t, param)
	assert.Equal(t, "aws_param_key1_value", param)

	param = params["aws_param_key2"]
	assert.NotNil(t, param)
	assert.Equal(t, "aws_param_key2_value", param)

	// Environment Provider
	p, _ = providerDescs.GetProvider("azure")
	assert.NotNil(t, p)

	// Environment Provider Labels
	labels = p.GetLabels()
	assert.NotNil(t, labels)
	assert.Equal(t, 2, len(labels.AsStrings()))
	assert.Equal(t, true, labels.Contains("azure_tag1_value", "azure_tag2_value"))

	// Environment Provider Parameters
	params = p.GetParameters().AsMap()
	assert.NotNil(t, params)
	param = params["azure_param_key1"]
	assert.NotNil(t, param)
	assert.Equal(t, "azure_param_key1_value", param)

	param = params["azure_param_key2"]
	assert.NotNil(t, param)
	assert.Equal(t, "azure_param_key2_value", param)

	//------------------------------------------------------------
	// Environment Nodes
	//------------------------------------------------------------
	nodeDescs := env.GetNodeDescriptions()
	assert.NotNil(t, nodeDescs)

	b = nodeDescs.Contains("node1")
	assert.Equal(t, true, b)
	b = nodeDescs.Contains("node1", "node2")
	assert.Equal(t, true, b)

	// Environment Node
	n, _ := nodeDescs.GetNode("node1")
	assert.NotNil(t, n)
	assert.Equal(t, 10, n.GetInstances())

	// Environment Node Labels
	labels = n.GetLabels()
	assert.NotNil(t, labels)
	assert.Equal(t, 3, len(labels.AsStrings()))
	assert.Equal(t, true, labels.Contains("node1_label1", "node1_label2", "node1_label3"))

	// Environment Node Provider
	name := n.GetProviderName()
	assert.NotNil(t, n)
	assert.Equal(t, "provider_node1", name.GetName())

	// Environment Node Provider Parameters
	params = name.GetParameters().AsMap()
	assert.NotNil(t, params)
	param = params["provider_node1_param_key1"]
	assert.NotNil(t, param)
	assert.Equal(t, "provider_node1_param_key1_value", param)
	param = params["provider_node1_param_key2"]
	assert.NotNil(t, param)
	assert.Equal(t, "provider_node1_param_key2_value", param)

	// Environment Node
	n, _ = nodeDescs.GetNode("node2")
	assert.NotNil(t, n)
	assert.Equal(t, 20, n.GetInstances())

	// Environment Node Labels
	labels = n.GetLabels()
	assert.NotNil(t, labels)
	assert.Equal(t, 3, len(labels.AsStrings()))
	assert.Equal(t, true, labels.Contains("node2_label1", "node2_label2", "node2_label3"))

	// Environment Node Provider
	name = n.GetProviderName()
	assert.NotNil(t, n)
	assert.Equal(t, "provider_node2", name.GetName())

	// Environment Node Provider Parameters
	params = name.GetParameters().AsMap()
	assert.NotNil(t, params)
	param = params["provider_node2_param_key1"]
	assert.NotNil(t, param)
	assert.Equal(t, "provider_node2_param_key1_value", param)
	param = params["provider_node2_param_key2"]
	assert.NotNil(t, param)
	assert.Equal(t, "provider_node2_param_key2_value", param)

	//------------------------------------------------------------
	// Environment Stacks
	//------------------------------------------------------------
	stackDescs := env.GetStackDescriptions()
	assert.NotNil(t, stackDescs)

	b = stackDescs.Contains("stack1")
	assert.Equal(t, true, b)
	b = stackDescs.Contains("stack1", "stack2")
	assert.Equal(t, true, b)

	// Environment Stack
	s, _ := stackDescs.GetStack("stack1")
	assert.NotNil(t, s)
	assert.Equal(t, "stack1_repository", s.GetRepository())
	assert.Equal(t, "stack1_version", s.GetVersion())

	// Environment Stack Labels
	labels = s.GetLabels()
	assert.NotNil(t, labels)
	assert.Equal(t, 3, len(labels.AsStrings()))
	assert.Equal(t, true, labels.Contains("stack1_label1", "stack1_label2", "stack1_label3"))

	// Environment Stack
	s, _ = stackDescs.GetStack("stack2")
	assert.NotNil(t, s)
	assert.Equal(t, "stack2_repository", s.GetRepository())
	assert.Equal(t, "stack2_version", s.GetVersion())

	// Environment Stack Labels
	labels = s.GetLabels()
	assert.NotNil(t, labels)
	assert.Equal(t, 3, len(labels.AsStrings()))
	assert.Equal(t, true, labels.Contains("stack2_label1", "stack2_label2", "stack2_label3"))

	//------------------------------------------------------------
	// Environment Tasks
	//------------------------------------------------------------
	taskDescs := env.GetTaskDescriptions()
	assert.NotNil(t, taskDescs)

	b = taskDescs.Contains("task1")
	assert.Equal(t, true, b)
	b = taskDescs.Contains("task1", "task2")
	assert.Equal(t, true, b)

	// Environment Task
	ts, _ := taskDescs.GetTask("task1")
	assert.NotNil(t, ts)
	assert.Equal(t, "task1_cron", ts.GetCron())
	assert.Equal(t, "task1_playbook", ts.GetPlaybook())

	// Environment Task Labels
	labels = ts.GetLabels()
	assert.NotNil(t, labels)
	assert.Equal(t, 3, len(labels.AsStrings()))
	assert.Equal(t, true, labels.Contains("task1_label1", "task1_label2", "task1_label3"))

	// Environment Task
	ts, _ = taskDescs.GetTask("task2")
	assert.NotNil(t, ts)
	assert.Equal(t, "task2_cron", ts.GetCron())
	assert.Equal(t, "task2_playbook", ts.GetPlaybook())

	// Environment Task Labels
	labels = ts.GetLabels()
	assert.NotNil(t, labels)
	assert.Equal(t, 3, len(labels.AsStrings()))
	assert.Equal(t, true, labels.Contains("task2_label1", "task2_label2", "task2_label3"))

}

package engine

import (
	_ "log"
	"testing"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestDebugDemo(t *testing.T) {

	p, _ := model.CreateParameters(map[string]interface{}{
		"ek": map[interface{}]interface{}{
			"aws": map[interface{}]interface{}{
				"region": "dummy",
				"accessKey": map[interface{}]interface{}{
					"id":     "dummy",
					"secret": "dummy",
				},
			},
		},
		"app": map[interface{}]interface{}{
			"visualizer": map[interface{}]interface{}{
				"port": "8080",
			},
		},
	})
	mainPath := "./testdata/gittest/descriptor"
	tc := model.CreateContext(p)
	c := &MockLaunchContext{locationContent: mainPath, templateContext: tc}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDesc := tester.createRep(mainPath)

	descContent := `
name: ekara-demo-al3
qualifier: dev

ekara:
  parent:
    repository: ekara-platform/distribution
  components:
    visualizer:
      repository: ekara-platform/swarm-visualizer-stack
nodes:
  "*":
    provider:
      name: ek-aws
      params:
        instance_type: "t2.micro"
        ami_id: "ami-f90a4880"
        vpc_id: "vpc-880aeaef"
        security_groups:
          - name: app
            rules:
              - proto: tcp
                ports:
                  - {{ .Vars.app.visualizer.port }}
                cidr_ip: 0.0.0.0/0
                rule_desc: allow all on port {{ .Vars.app.visualizer.port }} for the swarm visualizer 
    volumes:
      - path: /data1
        params:
          device_name: xvdf
          volume_type: gp2
          volume_size: 9
          delete_on_termination: true
      - path: /var/lib/docker
        params:
          device_name: xvdg
          volume_type: standard
          volume_size: 1
          delete_on_termination: true
          tags:
            Type: Docker
            TestTagKey: TestTagValue
  nodeset1:
    instances: 2
  nodeset2:
    instances: 1

stacks:
  visualizer:
    component: visualizer
`
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err := tester.initEngine()
	assert.Nil(t, err)
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	_, err = yaml.Marshal(env)
	assert.Nil(t, err)
	//log.Printf("--> yaml content %s", yamlContent)
}

package component

import (
	_ "log"
	"testing"

	"github.com/ekara-platform/engine/util"

	"github.com/ekara-platform/model"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestDebugDemo(t *testing.T) {

	p := model.CreateParameters(map[string]interface{}{
		"ekara": map[interface{}]interface{}{
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

	c := util.CreateMockLaunchContextWithData(mainPath, p, false)
	tester := CreateComponentTester(t, c)
	defer tester.Clean()

	repDesc := tester.CreateRep(mainPath)

	descContent := `
  name: ekaraDemoVar
  qualifier: dev
  
  ekara:
    parent:
      repository: ekara-platform/distribution
      components:
      visualizer:
        repository: ekara-platform/swarm-visualizer-stack
  
  nodes:
    _:
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
                  rule_desc: allow all on port {{ .Vars.app.visualizer.port }}
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
    nodeSet1:
      provider:
        name: ek-aws
      instances: 3
    nodeSet2:
      provider:
        name: ek-aws
      instances: 4
  
  stacks:
    visualizer:
      component: visualizer
  
`
	repDesc.WriteCommit("ekara.yaml", descContent)

	err := tester.Init()
	assert.Nil(t, err)
	env := tester.Env()
	assert.NotNil(t, env)

	_, err = yaml.Marshal(env)
	assert.Nil(t, err)
	refM := tester.rM
	assert.True(t, refM.usedReferences.IdUsed("ek-aws"))
	assert.True(t, refM.usedReferences.IdUsed("ek-swarm"))
	assert.True(t, refM.usedReferences.IdUsed("visualizer"))

	// comp1 should be downloaded because it's used as orchestrator into the parent
	// comp2 should not be downloaded because it's referenced by a component
	tester.AssertComponentsContainsExactly(model.MainComponentId, model.EkaraComponentId+"1", "ek-core", "ek-aws", "ek-swarm")
}

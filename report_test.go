package engine

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
	"github.com/stretchr/testify/assert"
)

type MockDescriber struct {
}

func TestReportContent(t *testing.T) {

	r := ExecutionReport{}
	_, err := r.Content()
	assert.Nil(t, err)

}

func TestReportContentSingleStep(t *testing.T) {

	r := ExecutionReport{}
	_, err := r.Content()
	assert.Nil(t, err)

	sc := InitCodeStepResult("DUMMY_STEP", nil, NoCleanUpRequired)
	r.Steps = sc.Array()
	_, err = r.Content()
	assert.Nil(t, err)

}

func TestReportContentSingleStepInstaller(t *testing.T) {

	r := ExecutionReport{}
	_, err := r.Content()
	assert.Nil(t, err)

	sc := InitCodeStepResult("DUMMY_STEP", nil, NoCleanUpRequired)
	FailsOnCode(&sc, fmt.Errorf("DUMMY_ERROR"), "", nil)
	r.Steps = sc.Array()
	_, err = r.Content()
	assert.Nil(t, err)

}

func TestReportContentSingleStepInstallerNilError(t *testing.T) {

	r := ExecutionReport{}
	_, err := r.Content()
	assert.Nil(t, err)

	sc := InitCodeStepResult("DUMMY_STEP", nil, NoCleanUpRequired)
	FailsOnCode(&sc, nil, "", nil)
	r.Steps = sc.Array()
	_, err = r.Content()
	assert.Nil(t, err)

}

func TestReportContentSingleStepDescriptor(t *testing.T) {

	r := ExecutionReport{}
	_, err := r.Content()
	assert.Nil(t, err)

	sc := InitCodeStepResult("DUMMY_STEP", nil, NoCleanUpRequired)
	FailsOnDescriptor(&sc, fmt.Errorf("DUMMY_ERROR"), "", nil)
	r.Steps = sc.Array()
	_, err = r.Content()
	assert.Nil(t, err)

}

func TestReportContentSingleStepDescriptorNilError(t *testing.T) {

	r := ExecutionReport{}
	_, err := r.Content()
	assert.Nil(t, err)

	sc := InitCodeStepResult("DUMMY_STEP", nil, NoCleanUpRequired)
	FailsOnDescriptor(&sc, nil, "", nil)
	r.Steps = sc.Array()
	_, err = r.Content()
	assert.Nil(t, err)

}

func TestReportContentMultipleSteps(t *testing.T) {

	r := ExecutionReport{}
	_, err := r.Content()
	assert.Nil(t, err)

	sCs := InitStepResults()
	sc1 := InitCodeStepResult("DUMMY_STEP1", nil, NoCleanUpRequired)
	sCs.Add(sc1)
	sc2 := InitCodeStepResult("DUMMY_STEP2", nil, NoCleanUpRequired)
	sCs.Add(sc2)
	sc3 := InitCodeStepResult("DUMMY_STEP2", nil, NoCleanUpRequired)
	sCs.Add(sc3)

	_, err = r.Content()
	assert.Nil(t, err)

}

func TestReadReport(t *testing.T) {
	var err error

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
	})
	mainPath := "./testdata/gittest/descriptor"
	tc := model.CreateContext(p)

	ef, err := util.CreateExchangeFolder("./testdata/report", "")
	assert.Nil(t, err)
	assert.NotNil(t, ef)
	assert.Nil(t, err)

	c := &MockLaunchContext{
		locationContent:      mainPath,
		templateContext:      tc,
		efolder:              ef,
		logger:               log.New(ioutil.Discard, "Test", log.Ldate|log.Ltime|log.Lmicroseconds),
		sshPublicKeyContent:  "sshPublicKey_content",
		sshPrivateKeyContent: "sshPrivateKey_content",
	}
	tester := gitTester(t, c, false)
	defer tester.clean()

	repDesc := tester.createRep(mainPath)

	descContent := `
name: NameContent
qualifier: QualifierContent

`
	repDesc.writeCommit(t, "ekara.yaml", descContent)

	err = tester.initEngine()
	assert.Nil(t, err)
	err = tester.context.engine.ComponentManager().Ensure()
	assert.Nil(t, err)
	env := tester.env()
	assert.NotNil(t, env)

	rC := &runtimeContext{}

	ok := ef.Output.Contains(ReportOutputFile)
	assert.True(t, ok)

	stepC := freport(c, rC)
	assert.NotNil(t, stepC)
	assert.NotNil(t, rC.report)
	assert.Equal(t, 3, len(rC.report.Results))
	has, cpt := rC.report.hasFailure()
	assert.True(t, has)
	assert.Equal(t, 0, len(cpt.otherFailures))
	assert.Equal(t, 1, len(cpt.playBookFailures))
}

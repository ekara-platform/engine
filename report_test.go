package engine

import (
	"fmt"
	"log"
	"os"
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

	ef, err := util.CreateExchangeFolder("./testdata/report", "")
	assert.Nil(t, err)

	c := MockLaunchContext{
		efolder:              ef,
		logger:               log.New(os.Stdout, util.InstallerLogPrefix, log.Ldate|log.Ltime|log.Lmicroseconds),
		sshPublicKeyContent:  "sshPublicKey_content",
		sshPrivateKeyContent: "sshPrivateKey_content",
		engine: EkaraMock{
			Env: model.Environment{
				Name:      "NameContent",
				Qualifier: "QualifierContent",
			},
		},
	}

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

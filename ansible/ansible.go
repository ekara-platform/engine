package ansible

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/lagoon-platform/engine/component"
	"github.com/lagoon-platform/engine/util"
	"github.com/lagoon-platform/model"
)

type AnsibleManager interface {
	Execute(component model.Component, playbook string, extraVars ExtraVars, envars EnvVars, inventories string) (error, int)
}

type ExecutionReport struct {
}

type context struct {
	logger           *log.Logger
	componentManager component.ComponentManager
}

func CreateAnsibleManager(logger *log.Logger, componentManager component.ComponentManager) AnsibleManager {
	return &context{
		logger:           logger,
		componentManager: componentManager}
}

func (ctx context) Execute(component model.Component, playbook string, extraVars ExtraVars, envars EnvVars, inventories string) (error, int) {
	// Path of the component where the plabook is supposed to be located
	path := ctx.componentManager.ComponentPath(component.Id)

	// We check if the component contains mmodules
	module := util.JoinPaths(path, util.ComponentModuleFolder)
	if _, err := os.Stat(module); err != nil {
		if os.IsNotExist(err) {
			ctx.logger.Printf("No module located in %s", module)
			module = ""
		}
	} else {
		ctx.logger.Printf("Module located in %s", module)
	}

	ctx.logger.Printf(util.LOG_LAUNCHING_PLAYBOOK, playbook)
	_, err := ioutil.ReadDir(path)
	if err != nil {
		return err, 0
	}

	var playBookPath string

	if strings.HasSuffix(path, "/") {
		playBookPath = path + playbook
	} else {
		playBookPath = path + "/" + playbook
	}

	if _, err := os.Stat(playBookPath); os.IsNotExist(err) {
		return err, 0
	} else {
		ctx.logger.Printf(util.LOG_STARTING_PLAYBOOK, playBookPath)
	}

	calls := make([]string, 0)
	calls = append(calls, playbook)
	if extraVars.Bool {
		ctx.logger.Printf("launched with extra vars :%s", extraVars.String())
		calls = append(calls, "--extra-vars")
		calls = append(calls, extraVars.String())
	}

	if module != "" {
		ctx.logger.Printf("launched with module :%s", module)
		calls = append(calls, "--module-path")
		calls = append(calls, module)
	}

	if inventories != "" {
		ctx.logger.Printf("launched with inventories :%s", inventories)
		strs := strings.Split(inventories, " ")
		for _, v := range strs {
			calls = append(calls, v)
		}
	}

	cmd := exec.Command("ansible-playbook", calls...)

	cmd.Env = os.Environ()
	for k, v := range envars.Content {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	cmd.Dir = path

	errReader, err := cmd.StderrPipe()
	if err != nil {
		return err, 0
	}
	logPipe(errReader, ctx.logger)

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return err, 0
	}
	logPipe(outReader, ctx.logger)

	err = cmd.Start()
	if err != nil {
		return err, 0
	}

	err = cmd.Wait()
	if err != nil {
		e, ok := err.(*exec.ExitError)
		if ok {
			s := e.Sys().(syscall.WaitStatus)
			code := s.ExitStatus()
			ctx.logger.Printf("Ansible returned error code : %v\n", ReturnedError(code))
			return err, code
		} else {
			return err, 0
		}
	}
	return nil, 0
}

// logPipe logs the given pipe, reader/closer on the given logger
func logPipe(rc io.ReadCloser, l *log.Logger) {
	s := bufio.NewScanner(rc)
	go func() {
		for s.Scan() {
			l.Printf("%s\n", s.Text())
		}
	}()
}

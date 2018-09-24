package ansible

import (
	"bufio"
	"github.com/lagoon-platform/engine/component"
	"github.com/lagoon-platform/engine/util"
	"github.com/lagoon-platform/model"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type AnsibleManager interface {
	Execute(component model.Component, playbook string, extraVars ExtraVars, envars EnvVars) int
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

func (ctx context) Execute(component model.Component, playbook string, extraVars ExtraVars, envars EnvVars) int {
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
		ctx.logger.Fatal(err)
	}

	var playBookPath string

	if strings.HasSuffix(path, "/") {
		playBookPath = path + playbook
	} else {
		playBookPath = path + "/" + playbook
	}

	if _, err := os.Stat(playBookPath); os.IsNotExist(err) {
		ctx.logger.Fatalf(util.ERROR_MISSING, playBookPath)
	} else {
		ctx.logger.Printf(util.LOG_STARTING_PLAYBOOK, playBookPath)
	}

	var cmd *exec.Cmd
	if extraVars.Bool {
		ctx.logger.Printf("launched with extra vars :%s", extraVars.String())
		if module == "" {
			ctx.logger.Printf("launched without module ")
			cmd = exec.Command("ansible-playbook", playbook, "--extra-vars", extraVars.String())
		} else {
			ctx.logger.Printf("launched with module :%s", module)
			cmd = exec.Command("ansible-playbook", playbook, "--extra-vars", extraVars.String(), "--module-path", module)
		}
	} else {
		ctx.logger.Printf("launched without extra vars")
		if module == "" {
			ctx.logger.Printf("launched without module ")
			cmd = exec.Command("ansible-playbook", playbook)
		} else {
			ctx.logger.Printf("launched with module :%s", module)
			cmd = exec.Command("ansible-playbook", playbook, "--module-path", module)
		}
	}

	cmd.Env = os.Environ()
	for k, v := range envars.Content {
		//ctx.logger.Printf("setting environment variable: %v=%v ", k, v)
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	cmd.Dir = path

	errReader, err := cmd.StderrPipe()
	if err != nil {
		ctx.logger.Fatal(err)
	}
	logPipe(errReader, ctx.logger)

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		ctx.logger.Fatal(err)
	}
	logPipe(outReader, ctx.logger)

	err = cmd.Start()
	if err != nil {
		ctx.logger.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		e, ok := err.(*exec.ExitError)
		if ok {
			s := e.Sys().(syscall.WaitStatus)
			code := s.ExitStatus()
			ctx.logger.Printf("---> ANSIBLE RETURNED ERROR : %v\n", ReturnedError(code))
			return code
			// TODO write report here
		} else {
			ctx.logger.Fatal(err)
		}
	}
	return 0
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

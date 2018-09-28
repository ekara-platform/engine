package ansible

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/lagoon-platform/engine"
	"github.com/lagoon-platform/model"
)

// LaunchPlayBook launches the playbook located into the given folder with the extraVars
//
//
//	Parameters:
//		manager: the component manager in charge of the localization of the component
//		component: the component holding the playbook to launch
//		playbook: the playbook to launch
//		extraVars: the extra vars passed to the playbook
//		envars: the environment variables set on the command launching the playbook
//		logger: the logger
func LaunchPlayBook(manager engine.ComponentManager, component model.Component, playbook string, extraVars engine.ExtraVars, envars engine.EnvVars, inventories string, logger log.Logger) (error, int) {
	// Path of the component where the plabook is supposed to be located
	path := manager.ComponentPath(component.Id)

	// We check if the component contains mmodules
	module := engine.JoinPaths(path, engine.ComponentModuleFolder)
	if _, err := os.Stat(module); err != nil {
		if os.IsNotExist(err) {
			logger.Printf("No module located in %s", module)
			module = ""
		}
	} else {
		logger.Printf("Module located in %s", module)
	}

	logger.Printf(engine.LOG_LAUNCHING_PLAYBOOK, playbook)
	_, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Println(err.Error())
		return err, 0
	}

	var playBookPath string

	if strings.HasSuffix(path, "/") {
		playBookPath = path + playbook
	} else {
		playBookPath = path + "/" + playbook
	}

	if _, err := os.Stat(playBookPath); os.IsNotExist(err) {
		e := fmt.Errorf(engine.ERROR_MISSING, playBookPath)
		logger.Println(e.Error())
		return e, 0

	} else {
		logger.Printf(engine.LOG_STARTING_PLAYBOOK, playBookPath)
	}

	calls := make([]string, 0)
	calls = append(calls, playbook)
	if extraVars.Bool {
		logger.Printf("launched with extra vars :%s", extraVars.String())
		calls = append(calls, "--extra-vars")
		calls = append(calls, extraVars.String())
	}

	if module != "" {
		logger.Printf("launched with module :%s", module)
		calls = append(calls, "--module-path")
		calls = append(calls, module)
	}

	if inventories != "" {
		logger.Printf("launched with inventories :%s", inventories)
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
		logger.Println(err.Error())
		return err, 0
	}
	logPipe(errReader, logger)

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		logger.Println(err.Error())
		return err, 0
	}
	logPipe(outReader, logger)

	err = cmd.Start()
	if err != nil {
		logger.Println(err.Error())
		return err, 0
	}

	err = cmd.Wait()
	if err != nil {
		e, ok := err.(*exec.ExitError)
		if ok {
			s := e.Sys().(syscall.WaitStatus)
			code := s.ExitStatus()
			logger.Printf("Ansible returned error code : %v\n", ReturnedError(code))
			return err, code
		} else {
			return err, 0
		}
	}
	return nil, 0
}

// logPipe logs the given pipe, reader/closer on the given logger
func logPipe(rc io.ReadCloser, l log.Logger) {
	s := bufio.NewScanner(rc)
	go func() {
		for s.Scan() {
			l.Printf("%s\n", s.Text())
		}
	}()
}

// BuildEquals converts a map[string]string into a succession
// of equalities of type "map key=map value" separated by a space
//
//	Example of returned value :
//		"key1=val1 key2=val2 key3=val3"
func BuildEquals(m map[string]string) string {
	var r string
	for k, v := range m {
		r = r + k + "=" + v + " "
	}
	return r
}

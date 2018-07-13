package engine

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/lagoon-platform/model"
)

// LaunchPlayBook launches the playbook located into the given folder with the extraVars
func LaunchPlayBook(manager ComponentManager, component model.Component, playbook string, extraVars ExtraVars, envars EnvVars, logger log.Logger) {
	// Path of the component where the plabook is supposed to be located
	path := manager.ComponentPath(component.Id)

	// We check if the component contains mmodules
	module := JoinPaths(path, ComponentModuleFolder)
	if _, err := os.Stat(module); err != nil {
		if os.IsNotExist(err) {
			logger.Printf("No module located in %s", module)
			module = ""
		}
	} else {
		logger.Printf("Module located in %s", module)
	}

	logger.Printf(LOG_LAUNCHING_PLAYBOOK, playbook)
	_, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Fatal(err)
	}

	var playBookPath string

	if strings.HasSuffix(path, "/") {
		playBookPath = path + playbook
	} else {
		playBookPath = path + "/" + playbook
	}

	if _, err := os.Stat(playBookPath); os.IsNotExist(err) {
		logger.Fatalf(ERROR_MISSING, playBookPath)
	} else {
		logger.Printf(LOG_STARTING_PLAYBOOK, playBookPath)
	}

	var cmd *exec.Cmd
	if extraVars.Bool {
		logger.Printf("launched with extra vars :%s", extraVars.String())
		if module == "" {
			logger.Printf("launched without module ")
			cmd = exec.Command("ansible-playbook", playbook, "--extra-vars", extraVars.String())
		} else {
			logger.Printf("launched with module :%s", module)
			cmd = exec.Command("ansible-playbook", playbook, "--extra-vars", extraVars.String(), "--module-path", module)
		}
	} else {
		logger.Printf("launched without extra vars")
		if module == "" {
			logger.Printf("launched without module ")
			cmd = exec.Command("ansible-playbook", playbook)
		} else {
			logger.Printf("launched with module :%s", module)
			cmd = exec.Command("ansible-playbook", playbook, "--module-path", module)
		}
	}

	cmd.Env = os.Environ()
	for k, v := range envars.content {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	cmd.Dir = path

	errReader, err := cmd.StderrPipe()
	if err != nil {
		logger.Fatal(err)
	}
	logPipe(errReader, logger)

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		logger.Fatal(err)
	}
	logPipe(outReader, logger)

	err = cmd.Start()
	if err != nil {
		logger.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		logger.Fatal(err)
	}
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

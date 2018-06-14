package engine

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// LaunchPlayBook launches the playbook located into the given folder with the extraVars
func LaunchPlayBook(folder string, playbook string, extraVars string, logger log.Logger, envars ...string) {
	logger.Printf(LOG_LAUNCHING_PLAYBOOK, playbook)
	_, err := ioutil.ReadDir(folder)
	if err != nil {
		logger.Fatal(err)
	}

	var playBookPath string

	if strings.HasSuffix(folder, "/") {
		playBookPath = folder + playbook
	} else {
		playBookPath = folder + "/" + playbook
	}

	if _, err := os.Stat(playBookPath); os.IsNotExist(err) {
		logger.Fatalf(ERROR_MISSING, playBookPath)
	}

	logger.Printf(LOG_STARTING_PLAYBOOK, playBookPath)

	var cmd *exec.Cmd
	if len(extraVars) > 0 {
		cmd = exec.Command("ansible-playbook", playbook, "--extra-vars", extraVars, "--module-path", "/opt/lagoon/ansible/core/modules")
	} else {
		cmd = exec.Command("ansible-playbook", playbook, "--module-path", "/opt/lagoon/ansible/core/modules")
	}
	cmd.Env = os.Environ()
	for _, v := range envars {
		cmd.Env = append(cmd.Env, v)
	}

	cmd.Dir = folder

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

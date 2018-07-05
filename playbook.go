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
func LaunchPlayBook(path string, playbook string, extraVars ExtraVars, module string, logger log.Logger, envars ...string) {
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
	for _, v := range envars {
		cmd.Env = append(cmd.Env, v)
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

type ExtraVars struct {
	Bool bool
	Vals []string
}

func (ev ExtraVars) String() string {
	r := ""
	for _, v := range ev.Vals {
		r = r + " " + v
	}
	return r
}

// BuildEquals converts a map[string]string into a succession
// of equalities of type "map key=map value"
// 	Example:
//		key1=val1 key2=val2 key3=val3
func BuildEquals(m map[string]string) string {
	var r string
	for k, v := range m {
		r = r + k + "=" + v + " "
	}
	return r
}

func BuildExtraVars(extraVars string, inputFolder FolderPath, outputFolder FolderPath) ExtraVars {
	r := ExtraVars{}
	r.Vals = make([]string, 0)

	vars := strings.Trim(extraVars, " ")
	in := strings.Trim(inputFolder.Path(), " ")
	out := strings.Trim(outputFolder.Path(), " ")
	if vars == "" && in == "" && out == "" {
		r.Bool = false
	} else {
		r.Bool = true
		r.Vals = append(r.Vals, "--extra-vars")
		if vars != "" {
			r.Vals = append(r.Vals, vars)
		}

		if in != "" {
			r.Vals = append(r.Vals, "input_dir="+in)
		}

		if out != "" {
			r.Vals = append(r.Vals, "output_dir="+out)
		}
	}
	return r
}

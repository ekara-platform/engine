package ansible

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"syscall"

	"github.com/ekara-platform/model"

	"github.com/ekara-platform/engine/component"
	"github.com/ekara-platform/engine/util"
)

type (
	//Manager executes an ansible playbook
	Manager interface {
		// Execute runs a playbook within a component
		//
		// Parameters:
		//		cr: the component holding the playbook to launch
		//		playbook: the name of the playbook to launch
		//		extraVars: the extra vars passed to the playbook
		//		envars: the environment variables set before launching the playbook
		//		data: the template context required to template a used component
		//
		Execute(cr component.UsableComponent, ctx model.TemplateContext, playbook string, extraVars ExtraVars, envars EnvVars) (int, error)
	}

	manager struct {
		l  *log.Logger
		cF component.Finder
	}
)

//CreateAnsibleManager returns a new AnsibleManager, able to launch playbook
func CreateAnsibleManager(l *log.Logger, cF component.Finder) Manager {
	return &manager{
		l:  l,
		cF: cF,
	}
}

func (aM manager) Execute(uc component.UsableComponent, ctx model.TemplateContext, playbook string, extraVars ExtraVars, envars EnvVars) (int, error) {

	ok, playBookPath := uc.ContainsFile(playbook)

	if !ok {
		return 0, fmt.Errorf("The component \"%s\" does not contains the playbook : %s", uc.Name(), playbook)
	}

	aM.l.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - -")
	aM.l.Println("* * * * * A N S I B L E - - P L A Y B O O K  * * * * ")
	aM.l.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - -")
	aM.l.Printf(util.LOG_STARTING_PLAYBOOK, playBookPath.RelativePath(), playBookPath.Component().RootPath())
	aM.l.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - -")

	var args = []string{playbook}
	modulePaths := aM.cF.ContainsDirectory(util.ComponentModuleFolder, ctx)
	defer modulePaths.Release()
	if modulePaths.Count() > 0 {
		pathsStrings := modulePaths.JoinAbsolutePaths(":")
		aM.l.Printf("Playbook modules directories: %s", pathsStrings)
		args = append(args, "--module-path", pathsStrings)
	} else {
		aM.l.Printf("No playbook module")
	}

	inventoryPaths := aM.cF.ContainsDirectory(util.InventoryModuleFolder, ctx)
	defer inventoryPaths.Release()
	if inventoryPaths.Count() > 0 {
		asArgs := inventoryPaths.PrefixPaths("-i")
		aM.l.Printf("Playbook inventory directories: %s", inventoryPaths.JoinAbsolutePaths(":"))
		for _, v := range asArgs {
			if v == "-i" {
				continue
			}
		}
		args = append(args, asArgs...)
	} else {
		aM.l.Printf("No playbook inventory")
	}

	if extraVars.Bool {
		aM.l.Printf("Playbook extra vars: %s", extraVars.String())
		args = append(args, "--extra-vars", extraVars.String())
	} else {
		aM.l.Printf("No playbook extra-vars")
	}

	cmd := exec.Command("ansible-playbook", args...)
	cmd.Dir = uc.RootPath()
	cmd.Env = []string{}
	for k, v := range envars.Content {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	if len(cmd.Env) > 0 {
		aM.l.Printf("Playbook environment vars: %s", cmd.Env)
	} else {
		aM.l.Printf("No playbook environment vars")
	}

	errReader, err := cmd.StderrPipe()
	if err != nil {
		return 0, err
	}
	logPipe(errReader, aM.l)

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	logPipe(outReader, aM.l)

	err = cmd.Start()
	if err != nil {
		return 0, err
	}

	err = cmd.Wait()
	if err != nil {
		e, ok := err.(*exec.ExitError)
		if ok {
			s := e.Sys().(syscall.WaitStatus)
			code := s.ExitStatus()
			aM.l.Printf("Ansible returned error code : %v\n", ReturnedError(code))
			return code, err
		}
		return 0, err

	}
	return 0, nil
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

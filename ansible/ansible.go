package ansible

import (
	"bufio"
	"fmt"
	"github.com/GroupePSA/componentizer"
	"github.com/ekara-platform/engine/model"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/ekara-platform/engine/util"
)

const (
	taskPrefix           = "TASK ["
	taskSuffix           = "]"
	virtualEnvBinaryPath = "/opt/virtualenvs/%s/bin"
	virtualEnvDefault    = "default"
)

type (
	//Manager executes an ansible playbook
	Manager interface {
		// Play runs a playbook within a component
		//
		// Parameters:
		//		cr: the component holding the playbook to launch
		//      ctx: the context used to template components
		//		playbook: the name of the playbook to launch
		//		extraVars: the extra vars passed to the playbook
		//		envVars: the environment variables set before launching the playbook
		//		fN: feedback notifier
		//
		Play(cr componentizer.UsableComponent, ctx componentizer.TemplateContext, playbook string, extraVars ExtraVars) (int, error)
		// Inventory returns the current inventory of environment nodes
		Inventory(ctx componentizer.TemplateContext) (Inventory, error)
	}

	manager struct {
		lC util.LaunchContext
		cM componentizer.ComponentManager
	}

	execChan struct {
		out    chan string
		err    chan string
		status chan int
	}
)

//CreateAnsibleManager returns a new AnsibleManager, providing managed execution of
//Ansible commands
func CreateAnsibleManager(lC util.LaunchContext, cM componentizer.ComponentManager) Manager {
	return &manager{
		lC: lC,
		cM: cM,
	}
}

func (aM manager) Play(uc componentizer.UsableComponent, ctx componentizer.TemplateContext, playbook string, extraVars ExtraVars) (int, error) {
	ok, playBookPath := uc.ContainsFile(playbook)
	if !ok {
		return 0, fmt.Errorf("component \"%s\" does not contain playbook: %s", uc.Id(), playbook)
	}
	aM.lC.Log().Printf("Executing playbook %s from component %s", playBookPath.RelativePath(), playBookPath.Owner().Id())

	var args = []string{playbook}

	// Discovered modules
	modulePaths := aM.findModulePaths(ctx)
	defer modulePaths.Release()
	args = append(args, aM.buildModuleArgs(modulePaths)...)

	// Discovered inventory sources
	inventoryPaths := aM.findInventoryPaths(ctx)
	defer inventoryPaths.Release()
	args = append(args, aM.buildInventoryArgs(inventoryPaths)...)

	// Component(s) env vars
	allComps := []componentizer.ComponentRef{uc.Source()}
	for _, mp := range modulePaths.Paths {
		allComps = append(allComps, mp.Owner().Source())
	}
	for _, mp := range inventoryPaths.Paths {
		allComps = append(allComps, mp.Owner().Source())
	}
	env := aM.buildEnvVars(allComps...)

	// Extra vars
	etxvs, err := aM.buildExtraVarsArgs(extraVars)
	if err != nil {
		return 0, err
	}
	args = append(args, etxvs...)

	// Verbosity = 2+
	v := aM.lC.Verbosity()
	if v > 1 {
		st := "-"
		for i := 1; i < v; i++ {
			st = st + "v"
		}
		args = append(args, st)
	}

	// SSH private key
	args = append(args, "--private-key="+aM.lC.SSHPrivateKey())

	// Execution
	log.Printf("Running the command \"ansible-playbook\" with arguments: %v", args)
	eC, err := aM.exec(uc.RootPath(), "ansible-playbook", args, env)
	if err != nil {
		return 0, err
	}

	storedLines := make([]string, 0)
	// Read the logs as they come until a status code is returned
	for {
		select {
		case errLine := <-eC.err:
			if aM.lC.Verbosity() > 0 {
				// log stderr directly
				aM.lC.Log().Println(errLine)
			} else {
				// keep stderr for later if playbook ends with error
				storedLines = append(storedLines, errLine)
			}
		case outLine := <-eC.out:
			// Detect tasks to show progression
			sTrim := strings.TrimSpace(outLine)
			if strings.Index(sTrim, "TASK [") == 0 {
				aM.lC.Feedback().Detail(sTrim[len(taskPrefix):strings.LastIndex(sTrim, taskSuffix)])
			}
			if aM.lC.Verbosity() > 0 {
				aM.lC.Log().Println(outLine)
			} else {
				// keep stdout for later if playbook ends with error
				storedLines = append(storedLines, outLine)
			}
		case status := <-eC.status:
			aM.lC.Log().Printf("Playbook finished (%d)", status)
			if status != 0 {
				aM.lC.Log().Printf("Failed playbook output below")
				for _, storeLine := range storedLines {
					aM.lC.Log().Println(storeLine)
				}
				return status, fmt.Errorf("playbook did not complete successfully (%d), check the logs for details", status)
			} else {
				return status, nil
			}
		}
	}
}

func (aM manager) Inventory(ctx componentizer.TemplateContext) (Inventory, error) {
	res := Inventory{}
	args := []string{"--list"}

	// Discovered inventory sources
	inventoryPaths := aM.findInventoryPaths(ctx)
	defer inventoryPaths.Release()
	args = append(args, aM.buildInventoryArgs(inventoryPaths)...)

	// Component(s) env vars
	var allComps []componentizer.ComponentRef
	for _, mp := range inventoryPaths.Paths {
		allComps = append(allComps, mp.Owner().Source())
	}
	env := aM.buildEnvVars(allComps...)

	log.Printf("Running the command \"ansible-inventory\" with arguments: %v", args)
	eC, err := aM.exec(os.TempDir(), "ansible-inventory", args, env)
	if err != nil {
		return res, err
	}

	// Read the output until a status code is returned
	sb := strings.Builder{}
	var finished bool
	for !finished {
		select {
		case errLine := <-eC.err:
			if aM.lC.Verbosity() > 0 {
				// log stderr directly
				aM.lC.Log().Println(errLine)
			}
		case outLine := <-eC.out:
			if aM.lC.Verbosity() > 0 {
				// log stdout directly
				aM.lC.Log().Println(outLine)
			}
			sb.WriteString(outLine)
		case status := <-eC.status:
			aM.lC.Log().Printf("Inventory done (%d)", status)
			finished = true
		}
	}

	// Parse the ansible output
	err = res.UnmarshalAnsibleInventory([]byte(sb.String()))
	if err != nil {
		return res, err
	}

	return res, nil
}

func (aM manager) buildModuleArgs(modulePaths componentizer.MatchingPaths) []string {
	var args []string
	if modulePaths.Count() > 0 {
		pathsStrings := modulePaths.JoinAbsolutePaths(":")
		aM.lC.Log().Printf("Ansible modules directories): %s", pathsStrings)
		args = append(args, "--module-path", pathsStrings)
	} else {
		aM.lC.Log().Printf("No Ansible module directory")
	}
	return args
}

func (aM manager) findModulePaths(ctx componentizer.TemplateContext) componentizer.MatchingPaths {
	return aM.cM.ContainsDirectory(util.ComponentModuleFolder, ctx)
}

func (aM manager) findInventoryPaths(ctx componentizer.TemplateContext) componentizer.MatchingPaths {
	return aM.cM.ContainsDirectory(util.InventoryModuleFolder, ctx)
}

func (aM manager) buildEnvVars(comps ...componentizer.ComponentRef) envVars {
	env := createEnvVars()
	env.addDefaultOsVars()
	// TODO: currently only the default virtual env is supported, later add logic to obtain a specific virtual env from the executed component
	env.prependToVar("PATH", fmt.Sprintf(virtualEnvBinaryPath, virtualEnvDefault))
	env.addProxy(aM.lC.Proxy())
	for _, c := range comps {
		if o, ok := c.(model.EnvVarsAware); ok {
			for envK, envV := range o.EnvVars() {
				env.add(envK, envV)
			}
		}
	}
	aM.lC.Log().Printf("Ansible environment variables: %s", env.String())
	return env
}

func (aM manager) buildInventoryArgs(inventoryPaths componentizer.MatchingPaths) []string {
	var args []string
	if inventoryPaths.Count() > 0 {
		asArgs := inventoryPaths.PrefixPaths("-i")
		aM.lC.Log().Printf("Ansible inventory directories: %s", inventoryPaths.JoinAbsolutePaths(":"))
		for _, v := range asArgs {
			if v == "-i" {
				continue
			}
		}
		args = append(args, asArgs...)
	} else {
		aM.lC.Log().Printf("No Ansible inventory directory")
	}
	return args
}

func (aM manager) buildExtraVarsArgs(extraVars ExtraVars) ([]string, error) {
	var args []string
	if !extraVars.Empty() {
		s, e := extraVars.String()
		if e != nil {
			return args, e
		}
		aM.lC.Log().Printf("Ansible extra vars: %s", s)
		args = append(args, "--extra-vars", s)
	} else {
		aM.lC.Log().Printf("No Ansible extra var")
	}
	return args, nil
}

func (aM manager) exec(dir string, ex string, args []string, envVars envVars) (execChan, error) {
	eC := execChan{
		out:    make(chan string),
		err:    make(chan string),
		status: make(chan int),
	}

	// TODO: currently only the default virtual env is supported, later add logic to obtain a specific virtual env from the executed component
	cmd := exec.Command(fmt.Sprintf(virtualEnvBinaryPath+"/%s", virtualEnvDefault, ex), args...)
	cmd.Dir = dir
	cmd.Env = []string{}
	for k, v := range envVars.Content {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	errReader, err := cmd.StderrPipe()
	if err != nil {
		return eC, err
	}
	logPipe(errReader, eC.err)

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return eC, err
	}
	logPipe(outReader, eC.out)

	err = cmd.Start()
	if err != nil {
		return eC, err
	}

	go func() {
		err = cmd.Wait()
		if err != nil {
			e, ok := err.(*exec.ExitError)
			if ok {
				s := e.Sys().(syscall.WaitStatus)
				eC.status <- s.ExitStatus()
			}
		} else {
			eC.status <- 0
		}
	}()

	return eC, nil
}

// logPipe logs the given pipe, reader/closer on the given logger
func logPipe(rc io.ReadCloser, ch chan string) {
	s := bufio.NewScanner(rc)
	go func() {
		for s.Scan() {
			ch <- s.Text()
		}
	}()
}

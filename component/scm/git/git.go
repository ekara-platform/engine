package git

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	"github.com/ekara-platform/model"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

const defaultGitRemoteName = "origin"

//GitScmHandler Represents the scm connector allowing to fecth GIT repositories.
//
//It implements "github.com/ekara-platform/engine/component/scm.scmHandler
type GitScmHandler struct {
	Logger *log.Logger
}

//Matches implements "github.com/ekara-platform/engine/component/scm.scmHandler
func (gitScm GitScmHandler) Matches(repository model.Repository, path string) bool {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return false
	}

	config, err := repo.Config()
	if err != nil {
		return false
	}

	for _, remoteConfig := range config.Remotes {
		if len(remoteConfig.URLs) > 0 && remoteConfig.URLs[0] == repository.Url.String() {
			return true
		}
	}

	return false
}

//Fetch implements "github.com/ekara-platform/engine/component/scm.scmHandler
func (gitScm GitScmHandler) Fetch(repository model.Repository, path string, auth model.Parameters) error {
	source := repository.Url.String()
	options := git.CloneOptions{
		RemoteName:   defaultGitRemoteName,
		URL:          source,
		SingleBranch: false,
		Tags:         git.AllTags}
	authMethod, e := buildAuthMethod(auth)
	if e != nil {
		return errors.New("error cloning git repository " + source + ": " + e.Error())
	}
	if authMethod != nil {
		options.Auth = authMethod
	}

	gitScm.Logger.Println("cloning GIT repository " + source)
	_, err := git.PlainClone(path, false, &options)
	if err != nil {
		return errors.New("unable to clone git repository " + source + ": " + err.Error())
	}
	return nil
}

//Update implements "github.com/ekara-platform/engine/component/scm.scmHandler
func (gitScm GitScmHandler) Update(path string, auth model.Parameters) error {
	options := git.FetchOptions{
		Tags: git.AllTags}
	authMethod, e := buildAuthMethod(auth)
	if e != nil {
		return errors.New("error updating git repository " + path + ": " + e.Error())
	}
	options.Auth = authMethod

	repo, err := git.PlainOpen(path)
	if err != nil {
		return errors.New("unable to open git repository " + path + ": " + err.Error())
	}
	config, err := repo.Config()
	if err != nil {
		return errors.New("unable to access config of git repository " + path + ": " + err.Error())
	}
	gitScm.Logger.Println("fetching latest data from " + config.Remotes[defaultGitRemoteName].URLs[0])
	err = repo.Fetch(&options)
	if err == git.NoErrAlreadyUpToDate {
		gitScm.Logger.Println("already up-to-date")
		return nil
	}
	return errors.New("unable to fetch latest data for git repository " + path + ": " + err.Error())
}

//Switch implements "github.com/ekara-platform/engine/component/scm.scmHandler
func (gitScm GitScmHandler) Switch(path string, ref string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return errors.New("unable to open git repository " + path + ": " + err.Error())
	}
	tree, err := repo.Worktree()
	if err != nil {
		return errors.New("unable to access work tree of git repository " + path + ": " + err.Error())
	}
	if ref != "" {
		// Only do a checkout if a ref is specified
		if strings.HasPrefix(ref, "refs/") {
			// Raw refs are checked out as-is
			gitScm.Logger.Println("checking out " + ref)
			err = checkout(tree, ref)
		} else {
			// Otherwise try tag first then branch
			gitScm.Logger.Println("checking out tag " + ref)
			err = checkout(tree, fmt.Sprintf("refs/tags/%s", ref))
			if err != nil {
				gitScm.Logger.Println("no tag named " + ref + " checking out branch instead")
				err = checkout(tree, fmt.Sprintf("refs/remotes/%s/%s", defaultGitRemoteName, ref))
			}
		}
	}
	if err != nil {
		return errors.New("unable to checkout " + ref + " in git repository " + path + ": " + err.Error())
	}
	return nil
}

func checkout(tree *git.Worktree, ref string) error {
	return tree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(ref),
		Force:  true})
}

func buildAuthMethod(auth model.Parameters) (transport.AuthMethod, error) {
	if authMethod, ok := auth["method"]; ok {
		switch authMethod {
		case "basic":
			return &http.BasicAuth{Username: auth["user"].(string), Password: auth["password"].(string)}, nil
		case "password":
			return &ssh.Password{User: auth["user"].(string), Password: auth["password"].(string)}, nil
		case "token":
			return &http.TokenAuth{Token: auth["token"].(string)}, nil
		default:
			return nil, errors.New("unknown GIT authentication method")
		}
	}
	return nil, nil
}

package component

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	"github.com/ekara-platform/model"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

const defaultGitRemoteName = "origin"

type GitScmHandler struct {
	ScmHandler
	logger *log.Logger
}

func (gitScm GitScmHandler) Matches(source *url.URL, path string) bool {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return false
	}

	config, err := repo.Config()
	if err != nil {
		return false
	}

	for _, remoteConfig := range config.Remotes {
		if len(remoteConfig.URLs) > 0 && remoteConfig.URLs[0] == source.String() {
			return true
		}
	}

	return false
}

func (gitScm GitScmHandler) Fetch(source *url.URL, path string, auth model.Parameters) error {
	options := git.CloneOptions{
		RemoteName:   defaultGitRemoteName,
		URL:          source.String(),
		SingleBranch: false,
		Tags:         git.AllTags}
	authMethod, e := buildAuthMethod(auth)
	if e != nil {
		return errors.New("error cloning git repository " + source.String() + ": " + e.Error())
	} else {
		if authMethod != nil {
			options.Auth = authMethod
		}
	}
	gitScm.logger.Println("cloning GIT repository " + source.String())
	_, err := git.PlainClone(path, false, &options)
	if err != nil {
		return errors.New("unable to clone git repository " + source.String() + ": " + err.Error())
	}
	return nil
}

func (gitScm GitScmHandler) Update(path string, auth model.Parameters) error {
	options := git.FetchOptions{
		Tags: git.AllTags}
	authMethod, e := buildAuthMethod(auth)
	if e != nil {
		return errors.New("error updating git repository " + path + ": " + e.Error())
	} else {
		options.Auth = authMethod
	}
	repo, err := git.PlainOpen(path)
	if err != nil {
		return errors.New("unable to open git repository " + path + ": " + err.Error())
	}
	config, err := repo.Config()
	if err != nil {
		return errors.New("unable to access config of git repository " + path + ": " + err.Error())
	}
	gitScm.logger.Println("fetching latest data from " + config.Remotes[defaultGitRemoteName].URLs[0])
	err = repo.Fetch(&options)
	if err == git.NoErrAlreadyUpToDate {
		gitScm.logger.Println("already up-to-date")
		return nil
	}
	return errors.New("unable to fetch latest data for git repository " + path + ": " + err.Error())
}

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
			gitScm.logger.Println("checking out " + ref)
			err = checkout(tree, ref)
		} else {
			// Otherwise try tag first then branch
			gitScm.logger.Println("checking out tag " + ref)
			err = checkout(tree, fmt.Sprintf("refs/tags/%s", ref))
			if err != nil {
				gitScm.logger.Println("no tag named " + ref + " checking out branch instead")
				err = checkout(tree, fmt.Sprintf("refs/remotes/%s/%s", defaultGitRemoteName, ref))
			}
		}
	}
	if err != nil {
		return errors.New("unable to checkout " + ref + " in git repository " + path + ": " + err.Error())
	} else {
		return nil
	}
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

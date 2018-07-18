package engine

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"log"
	"net/url"
	"strings"
	"fmt"
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

func (gitScm GitScmHandler) Fetch(source *url.URL, path string) error {
	gitScm.logger.Println("cloning GIT repository " + source.String())
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		RemoteName:   defaultGitRemoteName,
		URL:          source.String(),
		SingleBranch: false,
		Tags:         git.AllTags})
	if err != nil {
		return err
	}
	return nil
}

func (gitScm GitScmHandler) Update(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	config, err := repo.Config()
	if err != nil {
		return err
	}
	gitScm.logger.Println("fetching latest data from " + config.Remotes["origin"].URLs[0])
	err = repo.Fetch(&git.FetchOptions{
		Tags: git.AllTags})
	if err == git.NoErrAlreadyUpToDate {
		gitScm.logger.Println("already up-to-date")
		return nil
	}
	return err
}

func (gitScm GitScmHandler) Switch(path string, ref string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	tree, err := repo.Worktree()
	if err != nil {
		return err
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
	return err
}

func checkout(tree *git.Worktree, ref string) error {
	return tree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(ref),
		Force:  true})
}

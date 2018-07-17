package engine

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"log"
	"net/url"
	"strings"
)

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
	_, err := git.PlainClone(path, false, &git.CloneOptions{URL: source.String()})
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
	err = repo.Fetch(&git.FetchOptions{Tags: git.AllTags})
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
	if ref == "" {
		ref = "#master"
	}
	if strings.HasPrefix(ref, "#") {
		gitScm.logger.Println("checking out branch " + ref[1:])
		return tree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName("refs/heads/" + ref[1:]),
			Force:  true})
	} else {
		gitScm.logger.Println("checking out tag " + ref)
		return tree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName("refs/tags/" + ref),
			Force:  true})
	}
}

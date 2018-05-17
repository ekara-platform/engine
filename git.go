package engine

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"log"
	"net/url"
	"os"
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

func (gitScm GitScmHandler) Fetch(source *url.URL, dest string) error {
	gitScm.logger.Println("cloning GIT repository " + source.String())
	_, err := git.PlainClone(dest, false, &git.CloneOptions{
		Progress: os.Stdout,
		URL:      source.String()})
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
	return repo.Fetch(&git.FetchOptions{
		Progress: os.Stdout,
		Tags:     git.AllTags})
}

func (gitScm GitScmHandler) Switch(path string, tag string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	tree, err := repo.Worktree()
	if err != nil {
		return err
	}
	gitScm.logger.Println("checking out tag " + tag)
	return tree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/tags/" + tag),
		Force:  true})
}

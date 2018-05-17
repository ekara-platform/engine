package engine

import (
	"github.com/lagoon-platform/model"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

type ComponentManager interface {
	RegisterComponent(component model.Component)
	ComponentPath(id string) string
	ComponentsPaths() map[string]string

	Fetch(repository string, version string) (string, error)
	Ensure() error
}

type componentManager struct {
	logger     *log.Logger
	directory  string
	components map[string]model.Component
}

func createComponentManager(ctx *context) (cm ComponentManager, err error) {
	cm = &componentManager{
		logger:     ctx.logger,
		directory:  filepath.Join(ctx.baseDir, "components"),
		components: make(map[string]model.Component)}
	return
}

func (cm *componentManager) RegisterComponent(c model.Component) {
	cm.components[c.Id] = c
}

func (cm *componentManager) ComponentPath(id string) string {
	return filepath.Join(cm.directory, id)
}

func (cm *componentManager) ComponentsPaths() map[string]string {
	panic("implement me")
}

func (cm *componentManager) Fetch(location string, version string) (path string, err error) {
	return cm.gitFetch(location, version)
}

func (cm *componentManager) Ensure() error {
	panic("implement me")
}

func (cm *componentManager) gitFetch(location string, version string) (path string, err error) {
	cId, cUrl, err := model.ResolveRepositoryInfo(&url.URL{}, location)
	if err != nil {
		return
	}

	if !cm.isGitRepoAlreadyCloned(cId, cUrl) {
		cm.cleanPath(cId)
		path, err = cm.gitCloneRepo(cId, cUrl)
		if err != nil {
			return
		}
	}

	cm.gitFetchRepo(cId)

	cm.gitCheckoutTag(cId, "v"+version)

	return
}

func (cm *componentManager) isGitRepoAlreadyCloned(cId string, cUrl *url.URL) bool {
	path := filepath.Join(cm.directory, cId)
	repo, err := git.PlainOpen(path)
	if err != nil {
		return false
	}

	config, err := repo.Config()
	if err != nil {
		return false
	}

	for _, remoteConfig := range config.Remotes {
		if len(remoteConfig.URLs) > 0 && remoteConfig.URLs[0] == cUrl.String() {
			return true
		}
	}

	return false
}

func (cm *componentManager) gitCheckoutTag(cId string, tag string) error {
	path := filepath.Join(cm.directory, cId)
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	tree, err := repo.Worktree()
	if err != nil {
		return err
	}
	cm.logger.Println(cId + ": checking out tag " + tag)
	err = tree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/tags/" + tag),
		Force:  true})
	if err != nil {
		return err
	}
	return nil
}

func (cm *componentManager) gitFetchRepo(cId string) error {
	repo, err := git.PlainOpen(filepath.Join(cm.directory, cId))
	if err != nil {
		return err
	}
	cm.logger.Println(cId + ": fetching latest data from origin")
	return repo.Fetch(&git.FetchOptions{
		Progress: os.Stdout,
		Tags:     git.AllTags})
}

func (cm *componentManager) gitCloneRepo(cId string, cUrl *url.URL) (string, error) {
	path := filepath.Join(cm.directory, cId)
	cm.logger.Println(cId + ": cloning GIT repository " + cUrl.String())
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		Progress: os.Stdout,
		URL:      cUrl.String()})
	if err != nil {
		cm.logger.Println(cId + ": an error occurred during GIT clone, cleaning up")
		err2 := cm.cleanPath(cId)
		if err2 != nil {
			cm.logger.Println(cId + ": cleanup fail (" + err2.Error() + ")")
		}
		return "", err
	}
	cm.logger.Println(cId + ": successfully cloned GIT repository " + cUrl.String())
	return path, nil
}

func (cm *componentManager) cleanPath(cId string) error {
	path := filepath.Join(cm.directory, cId)
	if _, err := os.Stat(path); err == nil {
		cm.logger.Println(cId + ": cleaning up path " + path)
		return os.RemoveAll(path)
	}
	return nil
}

package component

import (
	"strings"

	"path/filepath"
)

type (
	//MatchingPath represents the matching path of the searched content
	MatchingPath interface {
		//Component gives the usable component wherin the searched content has been located
		Component() UsableComponent
		//RelativePath specifies the relatives path of the searched content into the usable component
		RelativePath() string
		//AbsolutePath specifies the absolute path of the searched content into the usable component
		AbsolutePath() string
	}

	mPath struct {
		comp         UsableComponent
		relativePath string
	}

	//MatchingPaths represents the matching paths of the searched content
	MatchingPaths struct {
		//Paths holds the searched  results
		Paths []MatchingPath
	}
)

func (p mPath) Component() UsableComponent {
	return p.comp
}

func (p mPath) RelativePath() string {
	return p.relativePath
}

func (p mPath) AbsolutePath() string {
	return filepath.Join(p.Component().RootPath(), p.RelativePath())
}

//Release deletes, if any, the templated paths returned
func (mp MatchingPaths) Release() {
	for _, v := range mp.Paths {
		v.Component().Release()
	}
}

//Count returns the number of matching paths
func (mp MatchingPaths) Count() int {
	return len(mp.Paths)
}

//JoinAbsolutePaths joins all the matching paths using the given separator
func (mp MatchingPaths) JoinAbsolutePaths(separator string) string {
	paths := make([]string, 0, 0)
	for _, v := range mp.Paths {
		paths = append(paths, filepath.Join(v.Component().RootPath(), v.RelativePath()))
	}
	return strings.Join(paths, separator)
}

//PrefixPaths returns the absolute mathcing paths prefixed with the given prefix
func (mp MatchingPaths) PrefixPaths(prefix string) []string {
	l := len(mp.Paths)
	res := make([]string, 0, 0)
	for i := 0; i < l; i++ {
		res = append(res, prefix)
		res = append(res, mp.Paths[i].AbsolutePath())
	}
	return res
}

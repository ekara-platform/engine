package component

import (
	"strings"

	"path/filepath"
)

type (
	MatchingPath interface {
		Component() UsableComponent
		RelativePath() string
	}

	mPath struct {
		comp         UsableComponent
		relativePath string
	}

	MatchingPaths struct {
		Paths []MatchingPath
	}
)

func (p mPath) Component() UsableComponent {
	return p.comp
}

func (p mPath) RelativePath() string {
	return p.relativePath
}

func (mp MatchingPaths) Release() {
	for _, v := range mp.Paths {
		v.Component().Release()
	}
}

func (mp MatchingPaths) Count() int {
	return len(mp.Paths)
}

func (mp MatchingPaths) JoinAbsolutePaths(separator string) string {
	paths := make([]string, 0, 0)
	for _, v := range mp.Paths {
		paths = append(paths, filepath.Join(v.Component().RootPath()), v.RelativePath())
	}
	return strings.Join(paths, separator)
}

package file

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ekara-platform/model"
)

//FileScmHandler Represents the scm connector allowing to fecth local repositories.
//
//It implements "github.com/ekara-platform/engine/component/scm.scmHandler
type FileScmHandler struct {
	Logger *log.Logger
}

//Matches implements "github.com/ekara-platform/engine/component/scm.scmHandler
func (fileScm FileScmHandler) Matches(repository model.Repository, path string) bool {
	// We always return false to force the repository to be fetched again
	return false
}

//Fetch implements "github.com/ekara-platform/engine/component/scm.scmHandler
func (fileScm FileScmHandler) Fetch(repository model.Repository, path string, auth model.Parameters) error {
	return copyDir(repository.Url.AsFilePath(), path)
}

//Update implements "github.com/ekara-platform/engine/component/scm.scmHandler
func (fileScm FileScmHandler) Update(path string, auth model.Parameters) error {
	// Doing nothing here and it's okay because Matches returns false
	// then the repo will be fetched/copied from scratch and never updated
	return nil
}

//Switch implements "github.com/ekara-platform/engine/component/scm.scmHandler
func (fileScm FileScmHandler) Switch(path string, ref string) error {
	// Doing nothing here and it's okay because we are dealing with
	// physical files then there  is nothing to switch...
	return nil
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func copyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = copyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

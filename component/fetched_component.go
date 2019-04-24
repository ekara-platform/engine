package component

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ekara-platform/model"
)

type (

	//FetchedComponent represents a component fetched locally
	FetchedComponent struct {
		id            string
		localPath     string
		descriptor    string
		descriptorUrl model.EkUrl
		localUrl      model.EkUrl
	}
)

func fetchThroughSccm(cm *context, c model.Component) (FetchedComponent, error) {
	fc := FetchedComponent{
		id:         c.Id,
		descriptor: c.Repository.DescriptorName,
	}

	// TODO dynamically select proper handler using a switch
	scm := GitScmHandler{logger: cm.logger}

	if c.Repository.Url.UpperScheme() == model.SchemeFile {
		cPath := filepath.Join(cm.directory, c.Id)
		err := CopyDir(c.Repository.Url.AsFilePath(), cPath)
		if err != nil {
			return fc, err
		}
		fc.localPath = cPath
	} else {
		cPath := filepath.Join(cm.directory, c.Id)
		if _, err := os.Stat(cPath); err == nil {
			if scm.Matches(c.Repository, cPath) {
				err := scm.Update(cPath, c.Authentication)
				if err != nil {
					return fc, err
				}
			} else {
				cm.logger.Println("directory " + cPath + " already exists but doesn't match component source, deleting it")
				err := os.RemoveAll(cPath)
				if err != nil {
					return fc, err
				}
				err = scm.Fetch(c.Repository, cPath, c.Authentication)
				if err != nil {
					return fc, err
				}
			}
		} else {
			err := scm.Fetch(c.Repository, cPath, c.Authentication)
			if err != nil {
				return fc, err
			}
		}
		err := scm.Switch(cPath, c.Ref)
		if err != nil {
			return fc, err
		}
		fc.localPath = cPath
	}
	u, err := model.CreateUrl(fc.localPath)
	if err != nil {
		return fc, err
	}
	fc.localUrl = u
	du, err := model.CreateUrl(filepath.Join(fc.localPath, fc.descriptor))
	if err != nil {
		return fc, err
	}
	fc.descriptorUrl = du
	return fc, nil
}

func (fc FetchedComponent) String() string {
	return fmt.Sprintf("id: %s, localPath: %s, descriptor: %s", fc.id, fc.localPath, fc.descriptor)
}

func (fc FetchedComponent) hasDescriptor() bool {
	if _, err := os.Stat(fc.descriptorUrl.AsFilePath()); err == nil {
		return true
	}
	return false
}

func CopyFile(src, dst string) (err error) {
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
func CopyDir(src string, dst string) (err error) {
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
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

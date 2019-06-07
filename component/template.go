package component

import (
	"html/template"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ekara-platform/model"
	gl "github.com/gobwas/glob"
	"github.com/oklog/ulid"
)

func runTemplate(ctx model.TemplateContext, componentPath string, patterns model.Patterns, resolver model.ComponentReferencer) (string, error) {
	if len(patterns.Content) > 0 {
		globs := make([]gl.Glob, 0, 0)
		files := make([]string, 0, 0)
		for _, p := range patterns.Content {
			pa := filepath.Join(componentPath, p)
			// workaround of issue : https://github.com/gobwas/glob/issues/35
			if runtime.GOOS == "windows" {
				pa = strings.Replace(pa, "\\", "/", -1)
			}
			globs = append(globs, gl.MustCompile(pa, '/'))
		}
		err := filepath.Walk(componentPath, func(path string, info os.FileInfo, err error) error {
			for _, p := range patterns.Content {
				pa := filepath.Join(componentPath, p)
				if path == pa {
					files = append(files, path)
					return nil
				}
			}

			for _, g := range globs {
				pa := path
				// workaround of issue : https://github.com/gobwas/glob/issues/35
				if runtime.GOOS == "windows" {
					pa = strings.Replace(pa, "\\", "/", -1)
				}
				if g.Match(pa) {
					files = append(files, path)
					return nil
				}
			}
			return nil
		})

		if err != nil {
			return "", err
		}

		var uuid string
		var tmpPath string

		// No matching files encountered then it won't be templated
		if len(files) == 0 {
			return "", nil
		} else {
			uuid = genUlid()
			tmpPath = componentPath + "_" + uuid
			err = copyFolder(componentPath, tmpPath)
			if err != nil {
				return "", err
			}
		}

		for _, f := range files {
			f = strings.Replace(f, componentPath, tmpPath, -1)

			t, err := template.ParseFiles(f)
			if err != nil {
				os.RemoveAll(tmpPath)
				return "", err
			}

			t = template.Must(t.Option("missingkey=error"), err)

			file, err := os.Create(f)
			defer file.Close()
			if err != nil {
				os.RemoveAll(tmpPath)
				return "", err
			}

			err = t.Execute(file, ctx)
			if err != nil {
				os.RemoveAll(tmpPath)
				return "", err
			}
		}
		return tmpPath, nil
	}
	return "", nil
}

func copyFolder(source string, dest string) (err error) {

	info, err := os.Stat(source)
	if err != nil {
		return
	}

	err = os.MkdirAll(dest, info.Mode())
	if err != nil {
		return
	}

	directory, err := os.Open(source)
	if err != nil {
		return
	}
	files, err := directory.Readdir(-1)
	if err != nil {
		return
	}
	for _, obj := range files {

		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			err = copyFolder(sourcefilepointer, destinationfilepointer)
			if err != nil {
				return
			}
		} else {
			err = copyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				return
			}
		}
	}
	return
}

func copyFile(source string, dest string) (err error) {
	s, err := os.Open(source)
	if err != nil {
		return
	}

	defer s.Close()

	d, err := os.Create(dest)
	if err != nil {
		return
	}

	defer d.Close()

	_, err = io.Copy(d, s)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}
	return
}

func genUlid() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}

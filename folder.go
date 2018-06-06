package engine

import (
	_ "log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type FolderPath interface {
	Path() string
	AdaptedPath() string
	Exixts() bool
}

type ExchangeFolder struct {
	Location FolderPath
	Input    FolderPath
	Output   FolderPath
}

func (f ExchangeFolder) Create() error {
	if _, err := os.Stat(f.Location.Path()); os.IsNotExist(err) {
		e := os.Mkdir(f.Location.Path(), 0700)
		if e != nil {
			return e
		}
	}
	if _, err := os.Stat(f.Input.Path()); os.IsNotExist(err) {
		e := os.Mkdir(f.Input.Path(), 0700)
		if e != nil {
			return e
		}
	}
	if _, err := os.Stat(f.Output.Path()); os.IsNotExist(err) {
		e := os.Mkdir(f.Output.Path(), 0700)
		if e != nil {
			return e
		}
	}
	return nil
}

func (f ExchangeFolder) CleanAllFiles() error {
	e := os.RemoveAll(f.Location.Path())
	if e != nil {
		return e
	}
	f.Create()
	return nil
}

// Adapt from c:\Users\e518546\goPathRoot\src\github.com\lagoon-platform\cli
// to /c/Users/e518546/goPathRoot/src/github.com/lagoon-platform/cli
func (f Folder) AdaptedPath() string {
	return AdaptPath(f.path)
}

func (f Folder) Path() string {
	return f.path
}

func (f Folder) Exixts() bool {
	_, err := os.Stat(f.path)
	return os.IsNotExist(err)

}

type Folder struct {
	path string
}

func ClientExchangeFolder(location string, clientName string) (ExchangeFolder, error) {
	outputDir, err := filepath.Abs(location)
	r := ExchangeFolder{}
	if err != nil {
		return r, err
	}
	r.Location = Folder{path: JoinPaths(outputDir, clientName)}
	r.Output = Folder{path: JoinPaths(outputDir, clientName, "output")}
	r.Input = Folder{path: JoinPaths(outputDir, clientName, "input")}
	return r, nil
}

func JoinPaths(paths ...string) string {
	var r string
	for _, v := range paths {
		r = filepath.Join(r, v)
	}
	return r
}

func AdaptPath(path string) string {
	s := path
	if runtime.GOOS == "windows" {
		if strings.Index(s, "c:\\") == 0 {
			s = "/c" + s[2:]
		}
		s = strings.Replace(s, "\\", "/", -1)
	}
	return s
}

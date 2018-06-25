package engine

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type FolderPath interface {
	Path() string
	AdaptedPath() string
	Exixts() bool
	Contains(s string) bool
}

// The ExchangeFolder is used to pass content to a container.
type ExchangeFolder struct {
	// Root location of the exchange folder
	Location FolderPath
	// Input data passed to the container
	Input FolderPath
	// Output data produced by the container
	Output FolderPath
}

// Create creates phisically the folder.
// If the folder already exist  it will remain untouched
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

// CleanAll cleans all the folder content
func (f ExchangeFolder) CleanAll() error {
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

func (f Folder) Contains(s string) bool {
	if _, err := os.Stat(JoinPaths(f.Path(), s)); os.IsNotExist(err) {
		return false
	}
	return true
}

type Folder struct {
	path string
}

func CreateExchangeFolder(location string, folderName string) (ExchangeFolder, error) {
	outputDir, err := filepath.Abs(location)
	r := ExchangeFolder{}
	if err != nil {
		return r, err
	}
	r.Location = Folder{path: JoinPaths(outputDir, folderName)}
	r.Output = Folder{path: JoinPaths(outputDir, folderName, "output")}
	r.Input = Folder{path: JoinPaths(outputDir, folderName, "input")}
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
	log.Printf("AdaptPath runtime %s", runtime.GOOS)
	if runtime.GOOS == "windows" {
		if strings.Index(s, "c:\\") == 0 {
			s = "/c" + s[2:]
		}
		s = strings.Replace(s, "\\", "/", -1)
	}
	log.Printf("AdaptPath func %s", s)
	return s
}

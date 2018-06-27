package engine

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// FolderPath represents a folder used into the ExchangeFolder structure
type FolderPath interface {
	//Path returns the path to this folder
	Path() string

	//AdaptedPath returns the path of this folder adapted depending
	// on the OS we are running on
	// For example on Windows it convert c:\Users\blablabla
	// to /c/Users/blablabla
	AdaptedPath() string

	// Exixts returns true if this folder exixts physically
	Exixts() bool

	// Contains returns true if this folder contains the given content
	Contains(s string) bool

	// ContainsParamYaml returns true if this folder contains a file named ParamYamlFileName
	ContainsParamYaml() (ok bool, b []byte, e error)

	// ContainsExtraVarsYaml returns true if this folder contains a file named ExtraVarYamlFileName
	ContainsExtraVarsYaml() (ok bool, b []byte, e error)

	// ContainsEnvYaml returns true if this folder contains a file named EnvYamlFileName
	ContainsEnvYaml() (ok bool, b []byte, e error)

	// Copy copies, if it exists, the content to the destination foler
	Copy(content string, destination FolderPath) error
}

// ExchangeFolder represents a folder structure used to pass and get content to container.
type ExchangeFolder struct {
	// Root location of the exchange folder
	Location FolderPath
	// Input data passed to the container
	Input FolderPath
	// Output data produced by the container
	Output FolderPath
}

// Create creates physically the ExchangeFolder.
// If the ExchangeFolder already exist  it will remain untouched
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

// CleanAll cleans all the ExchangeFolder content
func (f ExchangeFolder) CleanAll() error {
	e := os.RemoveAll(f.Location.Path())
	if e != nil {
		return e
	}
	f.Create()
	return nil
}

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

func (f Folder) ContainsParamYaml() (ok bool, b []byte, e error) {
	ok, b, e = containsAndRead(f, ParamYamlFileName)
	return
}

func (f Folder) ContainsExtraVarsYaml() (ok bool, b []byte, e error) {
	ok, b, e = containsAndRead(f, ExtraVarYamlFileName)
	return
}

func (f Folder) ContainsEnvYaml() (ok bool, b []byte, e error) {
	ok, b, e = containsAndRead(f, EnvYamlFileName)
	return
}

func (f Folder) Copy(content string, destination FolderPath) error {
	if f.Contains(content) {
		in, err := os.Open(JoinPaths(f.Path(), content))
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(JoinPaths(destination.Path(), content))
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
		return out.Close()
	} else {
		return fmt.Errorf("The content %s doesn't exist", content)
	}
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

func containsAndRead(f Folder, s string) (bool, []byte, error) {
	if ok := f.Contains(s); ok {
		b, err := ioutil.ReadFile(path.Join(f.Path(), s))
		if err != nil {
			return false, nil, err
		}
		return ok, b, nil
	}
	return false, nil, nil

}

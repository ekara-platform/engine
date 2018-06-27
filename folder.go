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
	e := f.Location.create()
	if e != nil {
		return e
	}
	e = f.Input.create()
	if e != nil {
		return e
	}
	e = f.Output.create()
	if e != nil {
		return e
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

// Delete deletes the ExchangeFolder
func (f ExchangeFolder) Delete() error {
	e := os.RemoveAll(f.Location.Path())
	if e != nil {
		return e
	}
	return nil
}

//AdaptedPath returns the path of this folder adapted depending
// on the OS we are running on
// For example on Windows it convert c:\Users\blablabla
// to /c/Users/blablabla
func (f FolderPath) AdaptedPath() string {
	return AdaptPath(f.path)
}

//Path returns the path to this folder
func (f FolderPath) Path() string {
	return f.path
}

// Exixts returns true if this folder exixts physically
func (f FolderPath) Exixts() (bool, error) {
	_, err := os.Stat(f.path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Contains returns true if this folder contains the given content
func (f FolderPath) Contains(s string) bool {
	if _, err := os.Stat(JoinPaths(f.Path(), s)); os.IsNotExist(err) {
		return false
	}
	return true
}

// ContainsParamYaml returns true if this folder contains a file named ParamYamlFileName
func (f FolderPath) ContainsParamYaml() (ok bool, b []byte, e error) {
	ok, b, e = containsAndRead(f, ParamYamlFileName)
	return
}

// ContainsExtraVarsYaml returns true if this folder contains a file named ExtraVarYamlFileName
func (f FolderPath) ContainsExtraVarsYaml() (ok bool, b []byte, e error) {
	ok, b, e = containsAndRead(f, ExtraVarYamlFileName)
	return
}

// ContainsEnvYaml returns true if this folder contains a file named EnvYamlFileName
func (f FolderPath) ContainsEnvYaml() (ok bool, b []byte, e error) {
	ok, b, e = containsAndRead(f, EnvYamlFileName)
	return
}

// Copy copies, if it exists, the content to the destination foler
func (f FolderPath) Copy(content string, destination FolderPath) error {
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

// AddChildExchangeFolder creates a child of type ExchangeFolder having
// this folder as root location
func (f *FolderPath) AddChildExchangeFolder(childName string) (ExchangeFolder, error) {
	child, err := CreateExchangeFolder(f.Path(), childName)
	if err != nil {
		return ExchangeFolder{}, err
	}
	f.Children[childName] = child
	return child, nil
}

// Write writes the given bytes into file within this folder
func (f FolderPath) Write(bytes []byte, fileName string) error {
	if ok, e := f.Exixts(); ok {
		if e != nil {
			return e
		}
		e = ioutil.WriteFile(JoinPaths(f.Path(), fileName), bytes, 0644)
		if e != nil {
			return e
		}
	} else {
		return fmt.Errorf("The directory is not created")
	}
	return nil
}

// AddChildExchangeFolder creates a child of type ExchangeFolder having
// this folder as root location
func (f FolderPath) create() error {
	if _, err := os.Stat(f.Path()); os.IsNotExist(err) {
		e := os.Mkdir(f.Path(), 0700)
		if e != nil {
			return e
		}
	}

	for _, v := range f.Children {
		e := v.Create()
		if e != nil {
			return e
		}
	}
	return nil
}

type FolderPath struct {
	path     string
	Children map[string]ExchangeFolder
}

func CreateExchangeFolder(location string, folderName string) (ExchangeFolder, error) {
	outputDir, err := filepath.Abs(location)
	r := ExchangeFolder{}
	if err != nil {
		return r, err
	}
	r.Location = FolderPath{path: JoinPaths(outputDir, folderName), Children: make(map[string]ExchangeFolder)}
	r.Output = FolderPath{path: JoinPaths(outputDir, folderName, "output"), Children: make(map[string]ExchangeFolder)}
	r.Input = FolderPath{path: JoinPaths(outputDir, folderName, "input"), Children: make(map[string]ExchangeFolder)}

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

func containsAndRead(f FolderPath, s string) (bool, []byte, error) {
	if ok := f.Contains(s); ok {
		b, err := ioutil.ReadFile(path.Join(f.Path(), s))
		if err != nil {
			return false, nil, err
		}
		return ok, b, nil
	}
	return false, nil, nil

}

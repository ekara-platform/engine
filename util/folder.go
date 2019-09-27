package util

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
	RootFolder string
	// Root location of the exchange folder
	Location FolderPath
	// Input data passed to the container
	Input FolderPath
	// Output data produced by the container
	Output FolderPath
}

// FolderPath holds the location of a folder and its eventual children
type FolderPath struct {
	path     string
	Children map[string]ExchangeFolder
}

// Create creates physically the ExchangeFolder.
// If the ExchangeFolder already exist  it will remain untouched
func (f ExchangeFolder) Create() error {
	if _, err := os.Stat(f.RootFolder); os.IsNotExist(err) {
		log.Printf("Creating exchange folder root %s", f.RootFolder)
		e := os.Mkdir(f.RootFolder, 0777)
		if e != nil {
			return e
		}
	}
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

// ContainsInventoryYamlFileName returns true if this folder contains a file named InventoryYamlFileName
func (f FolderPath) ContainsInventoryYamlFileName() (ok bool, b []byte, e error) {
	ok, b, e = containsAndRead(f, InventoryYamlFileName)
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
	}
	return fmt.Errorf("The content %s doesn't exist", content)

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
		e := os.Mkdir(f.Path(), 0777)
		if e != nil {
			return e
		}
	}
	if len(f.Children) > 0 {
		for _, v := range f.Children {

			e := v.Create()
			if e != nil {
				return e
			}
		}
	}

	return nil
}

//CreateFolderPath returns a FolderPath matching the given path
func CreateFolderPath(path string) FolderPath {
	return FolderPath{path: path}
}

//CreateExchangeFolder creates an ExchangeFolder at the provided location
// and with the give name
func CreateExchangeFolder(location string, folderName string) (ExchangeFolder, error) {
	r := ExchangeFolder{}
	var err error
	r.RootFolder, err = filepath.Abs(location)
	if err != nil {
		return r, err
	}

	r.Location = FolderPath{path: JoinPaths(r.RootFolder, folderName), Children: make(map[string]ExchangeFolder)}
	r.Output = FolderPath{path: JoinPaths(r.RootFolder, folderName, "output"), Children: make(map[string]ExchangeFolder)}
	r.Input = FolderPath{path: JoinPaths(r.RootFolder, folderName, "input"), Children: make(map[string]ExchangeFolder)}
	return r, nil
}

//JoinPaths joins the provided paths as a string
func JoinPaths(paths ...string) string {
	var r string
	for _, v := range paths {
		r = filepath.Join(r, v)
	}
	return r
}

//AdaptPath tunes the received path base on the runtime OS
//
//Example for Windows: "c:\\" will be replaced by "/c", and
//  "\\" by "/"
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

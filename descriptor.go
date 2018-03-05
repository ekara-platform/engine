package descriptor

import (
	"log"
	"gopkg.in/yaml.v2"
	"strings"
	"fmt"
	"io/ioutil"
	"path"
	"os"
	"path/filepath"
	"github.com/imdario/mergo"
	"net/http"
)

type Lagoon interface {
	GetDescriptor() Descriptor
}

type holder struct {
	// Global state
	logger *log.Logger

	// Descriptor info
	location string
	desc     *Descriptor
}

func (h *holder) GetDescriptor() (Descriptor) {
	return *h.desc
}

func Parse(logger *log.Logger, location string) (lagoon Lagoon, err error) {
	h := holder{logger: logger, location: location}
	desc, err := parseDescriptor(h.location)
	if err != nil {
		return nil, err
	}
	h.desc = &desc
	return &h, nil
}

func parseDescriptor(location string) (desc Descriptor, err error) {
	var content []byte

	desc.BaseLocation, content, err = readDescriptor(location)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(content, &desc)
	if err != nil {
		return
	}

	err = processImports(&desc)
	if err != nil {
		return
	}

	return
}

func readDescriptor(location string) (base string, content []byte, err error) {
	if strings.Index(location, "http") == 0 {
		fmt.Println("Loading URL", location)
		var response *http.Response
		response, err = http.Get(location)
		if err != nil {
			return
		}
		defer response.Body.Close()
		content, err = ioutil.ReadAll(response.Body)
		base = path.Dir(location) + "/"
	} else {
		fmt.Println("Loading file", location)
		var file *os.File
		file, err = os.Open(location)
		if err != nil {
			return
		}
		defer file.Close()
		content, err = ioutil.ReadAll(file)
		if err != nil {
			return
		}
		var absLocation string
		absLocation, err = filepath.Abs(location)
		if err != nil {
			return
		}
		base = filepath.Dir(absLocation) + string(filepath.Separator)
	}
	return
}

func processImports(desc *Descriptor) error {
	if len(desc.Imports) > 0 {
		fmt.Println("Processing imports", desc.Imports)
		for _, val := range desc.Imports {
			importedDesc, err := parseDescriptor(desc.BaseLocation + val)
			if err != nil {
				return err
			}
			mergo.Merge(desc, importedDesc)
		}
		desc.Imports = nil
	} else {
		fmt.Println("No import to process")
	}
	return nil
}

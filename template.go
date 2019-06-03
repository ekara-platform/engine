package engine

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ekara-platform/model"
	gl "github.com/gobwas/glob"
)

func runTemplate(lC LaunchContext, sCs *StepResults, sc *StepResult, patterns model.Patterns, resolver model.ComponentReferencer) error {
	if len(patterns.Content) > 0 {
		globs := make([]gl.Glob, 0, 0)
		files := make([]string, 0, 0)
		componentPath := lC.Ekara().ComponentManager().ComponentPath(resolver.ComponentName())
		for _, p := range patterns.Content {
			pa := filepath.Join(componentPath, p)
			globs = append(globs, gl.MustCompile(pa, '/'))
		}
		err := filepath.Walk(componentPath, func(path string, info os.FileInfo, err error) error {
			for _, p := range patterns.Content {
				pa := filepath.Join(componentPath, p)
				if path == pa {
					log.Printf("match found on pattern \"%s\" \n", path)
					files = append(files, path)
					return nil
				}
			}

			for _, g := range globs {
				if g.Match(path) {
					log.Printf("match found on glob \"%s\" \n", path)
					files = append(files, path)
					return nil
				}
			}
			return nil

		})

		if err != nil {
			FailsOnCode(sc, err, fmt.Sprintf("An error occurred walking into component : %s", resolver.ComponentName()), nil)
			sCs.Add(*sc)
			return nil
		}

		for _, f := range files {
			lC.Log().Printf("Applying template on file  \"%s\"", f)

			dat, err := ioutil.ReadFile(f)
			if err != nil {
				FailsOnCode(sc, err, fmt.Sprintf("An error occurred reading the file : %s", f), nil)
				sCs.Add(*sc)
				continue
			}

			no := f + ".origin"
			err = os.Rename(f, no)
			if err != nil {
				FailsOnCode(sc, err, fmt.Sprintf("An error occurred copying the file : %s", no), nil)
				sCs.Add(*sc)
				continue
			}

			tmpl := template.Must(template.New("").Option("missingkey=error").Parse(string(dat)))
			if err != nil {
				FailsOnCode(sc, err, fmt.Sprintf("An error occurred parsing the template content : %s", no), nil)
				sCs.Add(*sc)
				continue
			}
			fi, err := os.Create(f)
			if err != nil {
				FailsOnCode(sc, err, fmt.Sprintf("An error occurred creating the template target : %s", f), nil)
				sCs.Add(*sc)
				continue
			}
			defer fi.Close()
			w := bufio.NewWriter(fi)

			err = tmpl.Execute(w, lC.TemplateContext())
			if err != nil {
				FailsOnCode(sc, err, fmt.Sprintf("An error occurred writting the templated target : %s", f), nil)
				sCs.Add(*sc)
				continue
			}
			w.Flush()
		}
	}
	return nil
}

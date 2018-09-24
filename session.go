package engine

import (
	"encoding/json"
	"github.com/lagoon-platform/engine/util"
	"log"
	"os"
	"path"
)

type EngineSession struct {
	CreationSession *CreationSession
	File            string
}

type CreationSession struct {
	Client string
	Uids   map[string]string
}

//Add adds or replace the unique id for the given name
func (s *CreationSession) Add(name string, uid string) {
	s.Uids[name] = uid
}

// Content returns the json representation of the session
func (s *CreationSession) Content() (b []byte, e error) {
	b, e = json.Marshal(s)
	return
}

// TODO implement this
func (s *EngineSession) DeleteTodo() (b []byte, e error) {
	b, e = json.Marshal(s)
	return
}

func HasCreationSession(ef util.ExchangeFolder) (logged bool, session *EngineSession) {
	var s CreationSession

	file := path.Join(ef.Location.Path(), util.CreationSessionFileName)
	if _, err := os.Stat(file); err == nil {
		if data, err := os.Open(file); err == nil {

			defer data.Close()
			err = json.NewDecoder(data).Decode(&s)
			if err != nil {
				log.Fatal(err.Error())
			}
			logged = true
			session = &EngineSession{CreationSession: &s, File: file}
			return
		} else {
			log.Fatal(err.Error())
		}
	}
	logged = false
	session = &EngineSession{}
	return
}

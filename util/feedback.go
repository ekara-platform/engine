package util

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	FeedbackPrefix = "FEEDBACK: "
)

type (
	FeedbackNotifier interface {
		Info(message string, v ...interface{})
		Error(message string, v ...interface{})
		Progress(key string, message string, v ...interface{})
		ProgressG(key string, goal int, message string, v ...interface{})
		Detail(message string, v ...interface{})
	}

	FeedbackUpdate struct {
		Type    string `json:"t,omitempty"`
		Key     string `json:"k,omitempty"`
		Message string `json:"m,omitempty"`
		Goal    int    `json:"g,omitempty"`
	}

	logFeedbackNotifier struct {
		logger *log.Logger
	}
)

func (r logFeedbackNotifier) Info(message string, v ...interface{}) {
	r.marshal("I", "", 0, message, v)
}

func (r logFeedbackNotifier) Error(message string, v ...interface{}) {
	r.marshal("E", "", 0, message, v)
}

func (r logFeedbackNotifier) Progress(key string, message string, v ...interface{}) {
	r.marshal("P", key, 1, message, v)
}

func (r logFeedbackNotifier) ProgressG(key string, goal int, message string, v ...interface{}) {
	r.marshal("P", key, goal, message, v)
}

func (r logFeedbackNotifier) Detail(message string, v ...interface{}) {
	r.marshal("D", "", 0, message, v)
}

func (r logFeedbackNotifier) marshal(t string, key string, goal int, message string, v []interface{}) {
	b, err := json.Marshal(FeedbackUpdate{Type: t, Key: key, Goal: goal, Message: fmt.Sprintf(message, v...)})
	if err != nil {
		r.logger.Printf("Unable to marshall progress message: %s\n", message)
	} else {
		r.logger.Println(FeedbackPrefix + string(b))
	}
}

func CreateLoggingProgressNotifier(logger *log.Logger) FeedbackNotifier {
	return logFeedbackNotifier{logger: logger}
}

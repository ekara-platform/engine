package util

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	ProgressPrefix = "PROGRESS: "
)

type (
	ProgressNotifier interface {
		Detail(message string, v ...interface{})
		Notify(key string, message string, v ...interface{})
		NotifyWithGoal(key string, goal int, message string, v ...interface{})
	}

	progressUpdate struct {
		Key     string `json:"key,omitempty"`
		Message string `json:"msg,omitempty"`
		Count   int    `json:"c,omitempty"`
	}

	loggingProgressNotifier struct {
		logger *log.Logger
	}
)

func CreateProgressNotifier(logger *log.Logger) ProgressNotifier {
	return loggingProgressNotifier{logger: logger}
}

func (r loggingProgressNotifier) Detail(message string, v ...interface{}) {
	r.Notify("", message, v...)
}

func (r loggingProgressNotifier) Notify(key string, message string, v ...interface{}) {
	r.NotifyWithGoal(key, 0, message, v...)
}

func (r loggingProgressNotifier) NotifyWithGoal(key string, count int, message string, v ...interface{}) {
	b, err := json.Marshal(progressUpdate{Key: key, Count: count, Message: fmt.Sprintf(message, v...)})
	if err != nil {
		r.logger.Printf("Unable to marshall progress message: %s\n", message)
	} else {
		r.logger.Println(ProgressPrefix + string(b))
	}
}

package core

import (
	"fmt"
	"log"
)

type Logger interface {
	Info(messages ...string)
	Error(messages ...string)
	Warn(messages ...string)
}

type DefaultLogger struct {
	signature string
}

func NewDefaultLogger(signature string) *DefaultLogger {
	return &DefaultLogger{
		signature: signature,
	}
}

func (dl *DefaultLogger) Error(messages ...string) {
	dl.log("ERROR", messages...)
}

func (dl *DefaultLogger) Warn(messages ...string) {
	dl.log("WARN", messages...)
}

func (dl *DefaultLogger) Info(messages ...string) {
	dl.log("INFO", messages...)
}

func (dl *DefaultLogger) log(prefix string, messages ...string) {
	logString := fmt.Sprintf("%s [%s]: ", dl.signature, prefix)

	for index, message := range messages {
		logString += fmt.Sprintf("%s", message)

		if index != len(messages)-1 {
			logString += ": "
		}
	}

	log.Printf(logString)
}

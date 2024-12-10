package core

import (
	"fmt"
	"log"
)

type Logger interface {
	Log(messages ...string)
}

type DefaultLogger struct {
	signature string
}

func NewDefaultLogger(signature string) *DefaultLogger {
	return &DefaultLogger{
		signature: signature,
	}
}

func (dl *DefaultLogger) Log(messages ...string) {
	logString := fmt.Sprintf("%s: ", dl.signature)

	for index, message := range messages {
		logString += fmt.Sprintf("%s", message)

		if index != len(messages)-1 {
			logString += ": "
		}
	}
	log.Printf(logString)
}

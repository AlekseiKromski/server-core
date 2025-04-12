package core

import (
	"fmt"
	"log"
)

type Logger interface {
	Info(incoming any)
	Error(incoming any)
	Warn(incoming any)
}

type DefaultLogger struct {
	signature string
}

func NewDefaultLogger(signature string) *DefaultLogger {
	return &DefaultLogger{
		signature: signature,
	}
}

func (dl *DefaultLogger) Error(incoming any) {
	dl.log("ERROR", incoming)
}

func (dl *DefaultLogger) Warn(incoming any) {
	dl.log("WARN", incoming)
}

func (dl *DefaultLogger) Info(incoming any) {
	dl.log("INFO", incoming)
}

func (dl *DefaultLogger) log(prefix string, incoming any) {
	messages := []string{}

	switch v := incoming.(type) {
	case string:
		messages = append(messages, v)
	case []string:
		messages = append(messages, v...)
	default:
		messages = append(messages, "CANNOT PARSE INCOMING LOG INFORMATION")
	}

	logString := fmt.Sprintf("%s [%s]: ", dl.signature, prefix)

	for index, message := range messages {
		logString += fmt.Sprintf("%s", message)

		if index != len(messages)-1 {
			logString += ": "
		}
	}

	log.Printf(logString)
}

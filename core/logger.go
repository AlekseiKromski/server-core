package core

import (
	"fmt"
	"log"
	"encoding/json"
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
	case error:
		messages = append(messages, v.Error())
	case interface{}:
		encoded, err := json.Marshal(v)
		if err != nil {
			dl.Error(fmt.Errorf("cannot encode interface{} to json string: %v", err))
			return
		}
		messages = append(messages, string(encoded))
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

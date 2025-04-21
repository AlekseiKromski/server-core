package core

type DefaultLogEntry struct {
	Signature string `json:"signature"`
	Message   string `json:"message"`
}

func NewDefaultLogEntry(signature, message string) *DefaultLogEntry {
	return &DefaultLogEntry{
		Signature: signature,
		Message:   message,
	}
}

type SignedLogger interface {
	SetSignature(signature string)
	Logger
}

type DefaultSignedLogger struct {
	signature string
	logger    Logger
}

func NewDefaultSignedLogger(logger Logger) *DefaultSignedLogger {
	return &DefaultSignedLogger{
		logger: logger,
	}
}

func (dsl *DefaultSignedLogger) SetSignature(signature string) {
	dsl.signature = signature
}

func (dsl *DefaultSignedLogger) Info(incoming any) {
	entry := dsl.prepareDefaultLogEntry(incoming)
	if entry == nil {
		dsl.Error("incoming data for log is not a allowed type")
		return
	}

	dsl.logger.Info(entry)
}

func (dsl *DefaultSignedLogger) Error(incoming any) {
	entry := dsl.prepareDefaultLogEntry(incoming)
	if entry == nil {
		dsl.Error("incoming data for log is not a allowed type")
		return
	}

	dsl.logger.Error(entry)
}

func (dsl *DefaultSignedLogger) Warn(incoming any) {
	entry := dsl.prepareDefaultLogEntry(incoming)
	if entry == nil {
		dsl.Error("incoming data for log is not a allowed type")
		return
	}

	dsl.logger.Warn(entry)
}

func (dsl *DefaultSignedLogger) prepareDefaultLogEntry(incoming any) *DefaultLogEntry {
	entry := ""
	switch incm := incoming.(type) {
	case string:
		entry = incm
	case error:
		entry = incm.Error()
	default:
		return nil
	}

	return NewDefaultLogEntry(dsl.signature, entry)
}

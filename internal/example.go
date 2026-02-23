package internal

import "github.com/AlekseiKromski/server-core/core"

// Check module implementation during compile time
var _ core.Module = &Module{}

type Module struct {
	core.SignedLogger
}

func NewModule() *Module {
	m := &Module{}
	logger := core.NewDefaultLogger(m.Signature())
	m.SignedLogger = core.NewDefaultSignedLogger(logger)

	return m
}

func (m *Module) Start(notifyChannel chan struct{}, eventBusSender func(event *core.BusEvent), requirements map[core.Signature]core.Module) {
	notifyChannel <- struct{}{}
}

func (m *Module) Stop() {}

func (m *Module) Require() []core.Signature { return []core.Signature{} }

func (m *Module) Signature() core.Signature { return core.Signature("module") }

func (m *Module) Listen(event *core.BusEvent) {}

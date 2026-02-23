package core

type Module interface {
	Start(notifyChannel chan struct{}, eventBusSender func(event *BusEvent), requirements map[Signature]Module) // Start module
	Stop()                                                                                                      // Stop module
	Signature() Signature                                                                                       // Signature unique name // Log should have log mechanism
	SignedLogger
}

type Require interface {
	Require() []Signature // Require list of required modules
}

type Listener interface {
	Listen(event *BusEvent) // An entrypoint for BusEvent event
}

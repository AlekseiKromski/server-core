package core

type Module interface {
	Start(
		notifyChannel chan struct{},
		eventBusSender func(event *BusEvent),
		requirements map[Signature]Module,
	) // Start module
	Stop()                // Stop module
	Signature() Signature // Signature unique name
	SignedLogger          // Log should have log mechanism
}

type Require interface {
	Require() []Signature // Require list of required modules
}

type Listener interface {
	Listen(event *BusEvent)
}

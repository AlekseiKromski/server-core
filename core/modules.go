package core

type Module interface {
	Start(notifyChannel chan struct{}, eventBusSender func(event *BusEvent), requirements map[Signature]Module) // Start module
	Stop()                                                                                                      // Stop module
	Require() []Signature                                                                                       // Require list of required modules
	Signature() Signature                                                                                       // Signature unique name // Log should have log mechanism
	Listen(event *BusEvent)                                                                                     // An entrypoint for BusEvent event
	SignedLogger
}

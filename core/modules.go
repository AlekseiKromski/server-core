package core

type Module interface {
	Start(notifyChannel chan struct{}, busEventChannel chan BusEvent, requirements map[string]Module) // Start module
	Stop()                                                                                            // Stop module
	Require() []string                                                                                // Require list of required modules
	Signature() string                                                                                // Signature unique name // Log should have log mechanism
	Logger                                                                                            // Logger custom or core.DefaultLogger logger
}

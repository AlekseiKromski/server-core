package core

import (
	"log"
	"os"
	"os/signal"
	"time"
)

type Core struct {
	systemChannel   chan os.Signal
	notifyChannel   chan struct{}
	eventBusChannel chan BusEvent
}

func NewCore() *Core {
	return &Core{
		systemChannel:   make(chan os.Signal, 1),
		notifyChannel:   make(chan struct{}, 1),
		eventBusChannel: make(chan BusEvent, 1),
	}
}

func (c *Core) Init(modules []Module) {
	// Subscribe to Interrupt syscall
	signal.Notify(c.systemChannel, os.Interrupt)

	// Create requirements tree
	log.Println("Core: create requirements tree")
	modules, requirements := c.createRequireTree(modules)
	log.Println("Core: requirements tree created")

	c.startModules(modules, requirements)

	// Block main thread until stop
	<-c.systemChannel

	c.stopModules(modules)

	// Exit from application
	os.Exit(0)
}

func (c *Core) startModules(modules []Module, requirements map[string]Module) {
	defer close(c.notifyChannel)

	startTime := time.Now()
	log.Printf("Core: Start %d modules", len(modules))
	for _, module := range modules {
		mReqs := map[string]Module{}

		for _, requirement := range module.Require() {
			mReqs[requirement] = requirements[requirement]
		}

		go module.Start(c.notifyChannel, c.eventBusChannel, mReqs)

		//Wait until module start
		<-c.notifyChannel
	}
	log.Printf("Core: All modules started in %f seconds", time.Now().Sub(startTime).Seconds())
}

// createRequireTree return a list of sorted modules by requirements
func (c *Core) createRequireTree(modules []Module) ([]Module, map[string]Module) {
	requirements := map[string]Module{}
	var sortedModules []Module

	index := 0
Main:
	for {
		// Fallback to start if index more, that count of  original modules
		if len(modules) == index {
			index = 0
		}

		m := modules[index]
		signature := m.Signature()

		//Skip already existed modules
		if requirements[signature] != nil {
			index++
			continue
		}

		mReqs := m.Require()
		for _, req := range mReqs {
			// If no requirement, let's process next
			if requirements[req] == nil {
				index++
				continue Main
			}
		}

		// Save only not existed
		requirements[signature] = m
		sortedModules = append(sortedModules, m)

		// Exit if len of origin modules equals len on sorted
		if len(modules) == len(sortedModules) {
			break
		}

		index++
	}

	return sortedModules, requirements
}

func (c *Core) stopModules(modules []Module) {
	startTime := time.Now()
	log.Printf("Core: will stop %d modules", len(modules))

	for _, module := range modules {
		module.Stop()
	}

	log.Printf("Core: All modules stopped in %f seconds", time.Now().Sub(startTime).Seconds())
}

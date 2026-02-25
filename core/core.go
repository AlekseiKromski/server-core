package core

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Core struct {
	notifyChannel  chan struct{}
	eventBusSender *eventBus
}

func NewCore() *Core {
	return &Core{
		notifyChannel:  make(chan struct{}, 1),
		eventBusSender: newEventBus(),
	}
}

func (c *Core) Init(modules []Module) {
	// Create a context for graceful shutdown for core modules only
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Subscribe to Interrupt syscall
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // Kubernetes/Docker termination
		syscall.SIGQUIT, // Ctrl+\
	)

	// Create requirements tree
	log.Println("Core: create requirements tree")
	modules, requirements := c.createRequireTree(modules)
	log.Println("Core: requirements tree created")

	c.startModules(modules, requirements)

	// All modules started, let's start Listener
	c.eventBusSender.listen(ctx, modules)

	// Block main thread until stop
	<-sigChan

	c.stopModules(modules)

	// Exit from application
	os.Exit(0)
}

func (c *Core) startModules(modules []Module, requirements map[Signature]Module) {
	defer close(c.notifyChannel)

	startTime := time.Now()
	log.Printf("Core: Start %d modules", len(modules))
	for _, module := range modules {
		mReqs := map[Signature]Module{}

		// If we have some requirements, load them
		if r, ok := module.(Require); ok {
			for _, requirement := range r.Require() {
				mReqs[requirement] = requirements[requirement]
			}
		}

		go module.Start(c.notifyChannel, c.eventBusSender.send, mReqs)

		//Wait until module start
		<-c.notifyChannel
	}
	log.Printf("Core: All modules started in %f seconds", time.Now().Sub(startTime).Seconds())
}

// createRequireTree return a list of sorted modules by requirements
func (c *Core) createRequireTree(modules []Module) ([]Module, map[Signature]Module) {
	requirements := map[Signature]Module{}
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

		// If we implement Require interface, we can get requirements
		// Otherwise empty array
		mReqs := []Signature{}
		if m, ok := m.(Require); ok {
			mReqs = m.Require()
		}

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

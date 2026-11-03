package core

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
)

type Core struct {
	// Notify all external listeners, that application started by closing channel
	// There will be no `struct{}` events - this avoids the issue where an event
	// might be sent to only a single listener. Since closing a channel affects
	// all listeners, the signal will be received universally.
	ApplicationStarted chan struct{}

	eventBusSender *eventBus
}

func NewCore() *Core {
	return &Core{
		ApplicationStarted: make(chan struct{}, 1),
		eventBusSender:     newEventBus(),
	}
}

func (c *Core) Init(modules []Module) {
	c.init(context.Background(), modules)
}

// InitManual - if you want manually handle context
// You don't need strictly use it, it's useful for testing packages.
// If you are not sure, please use Init() instead
func (c *Core) InitManual(ctx context.Context, modules []Module) {
	c.init(ctx, modules)
}

func (c *Core) init(ctx context.Context, modules []Module) {
	// Subscribe to Interrupt syscall
	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // Kubernetes/Docker termination
		syscall.SIGQUIT, // Ctrl+\
	)
	defer cancel()

	// Create requirements tree
	log.Println("Core: create requirements tree")
	modules, requirements := c.createRequireTree(modules)
	log.Println("Core: requirements tree created")

	modulesStartedChannel := make(chan struct{}, 1)
	go c.startModules(ctx, modulesStartedChannel, modules, requirements)

	// All modules started, let's start Listener
	go c.eventBusSender.listen(ctx, modules)

	// Notify externals, that application started by closing channel
	close(c.ApplicationStarted)

	// Block main thread until stop
	<-ctx.Done()

	// Wait until we finish with core gorutines
	<-modulesStartedChannel
	c.eventBusSender.runnedListeners.Wait()

	// Stop modules
	c.stopModules(modules)
}

func (c *Core) startModules(ctx context.Context, modulesStartedChannel chan struct{}, modules []Module, requirements map[Signature]Module) {
	defer close(modulesStartedChannel)

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

		go module.Start(modulesStartedChannel, c.eventBusSender.send, mReqs)

		//Wait until module start or application stopped
		select {
		case <-modulesStartedChannel:
			continue
		case <-ctx.Done():
			log.Printf("Core: cannot start %s module, force exit", module.Signature())
			return
		}
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

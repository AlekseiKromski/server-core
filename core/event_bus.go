package core

import (
	"context"
	"sync"
)

type BusEvent struct {
	Receiver Signature
	Payload  interface{}
}

func NewBusEvent(r Signature, p interface{}) *BusEvent {
	return &BusEvent{
		Receiver: r,
		Payload:  p,
	}
}

type eventBus struct {
	c               chan *BusEvent
	runnedListeners sync.WaitGroup
	mu              sync.RWMutex
	closed          bool
}

func newEventBus() *eventBus {
	return &eventBus{
		c:               make(chan *BusEvent, 1),
		runnedListeners: sync.WaitGroup{},
	}
}

func (eb *eventBus) listen(ctx context.Context, modules []Module) {

	// Convert modules list to modules map for faster processing
	modulesMap := map[Signature]Listener{}
	for _, m := range modules {
		// Skip non-listeners
		if listener, ok := m.(Listener); ok {
			modulesMap[m.Signature()] = listener
		}
	}

	// Start context monitoring gorutine
	go func() {
		<-ctx.Done()

		// Prevent to send new events by locking
		eb.mu.Lock()

		eb.closed = true
		close(eb.c)

		// Unlock everything
		eb.mu.Unlock()
	}()

	// Read from channel until context canceled
	for e := range eb.c {
		if m, ok := modulesMap[e.Receiver]; ok {
			eb.runnedListeners.Add(1)

			go func(l Listener, ev *BusEvent) {
				defer eb.runnedListeners.Done()

				l.Listen(ev)
			}(m, e)
		}
	}

	eb.runnedListeners.Wait()
}

func (eb *eventBus) send(event *BusEvent) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if eb.closed {
		return
	}

	eb.c <- event
}

func (eb *eventBus) wait() {
	// Wait until all listeners complete their work
	eb.runnedListeners.Wait()
}

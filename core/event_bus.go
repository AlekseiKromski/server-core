package core

import (
	"context"
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
	c      chan *BusEvent
	ctx    context.Context
	cancel context.CancelFunc
}

func newEventBus() *eventBus {
	ctx, cancel := context.WithCancel(context.Background())
	return &eventBus{
		c:      make(chan *BusEvent, 1),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (eb *eventBus) listen(modules []Module) {

	// Convert modules list to modules map for faster processing
	modulesMap := map[Signature]Module{}
	for _, m := range modules {
		modulesMap[m.Signature()] = m
	}

	for {
		select {
		case e := <-eb.c:
			// Forward event directly to module
			if m := modulesMap[e.Receiver]; m != nil {
				m.Listen(e)
				continue
			}
		case <-eb.ctx.Done():
			return
		}
	}
}

func (eb *eventBus) send(event *BusEvent) {
	// Send event to channel
	eb.c <- event
}

func (eb *eventBus) stop() {
	// Stop all event bus listening proccess
	eb.cancel()
}

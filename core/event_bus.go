package core

import "context"

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
	c chan *BusEvent
}

func newEventBus() *eventBus {
	return &eventBus{
		c: make(chan *BusEvent, 1),
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

	for {
		select {
		case e := <-eb.c:
			// Forward event directly to module
			if m := modulesMap[e.Receiver]; m != nil {
				m.Listen(e)
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}

func (eb *eventBus) send(event *BusEvent) {
	// Send event to channel
	eb.c <- event
}

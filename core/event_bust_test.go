package core

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Check module implementation during compile time
var _ Module = &A{}
var _ Listener = &A{}
var _ Require = &A{}

type A struct {
	ebs            func(e *BusEvent)
	listenerAction func(ebs func(e *BusEvent), e *BusEvent)
	signature      Signature
	SignedLogger
}

func NewA(logger Logger, sig Signature, listenerAction func(ebs func(e *BusEvent), e *BusEvent)) *A {
	a := &A{
		listenerAction: listenerAction,
		signature:      sig,
	}

	a.SignedLogger = NewDefaultSignedLogger(logger)
	a.SignedLogger.SetSignature(a.Signature())

	return a
}

func (a *A) Start(notifyChannel chan struct{}, eventBusSender func(event *BusEvent), requirements map[Signature]Module) {
	a.ebs = eventBusSender

	notifyChannel <- struct{}{}
}

func (a *A) Stop() {

}

func (a *A) Require() []Signature {
	return []Signature{}
}

func (a *A) Signature() Signature {
	return Signature(a.signature)
}

func (a *A) Listen(event *BusEvent) {
	a.listenerAction(a.ebs, event)
}

func TestEventBusDeliveryToOneModule(t *testing.T) {
	logger := NewDefaultLogger(Signature("simple-logger"))
	nc := make(chan struct{}, 1)
	eb := newEventBus()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	a := NewA(
		logger,
		Signature("A"),
		func(_ func(*BusEvent), e *BusEvent) {
			// Check received event and unblock waitgroup
			if m, ok := e.Payload.(string); ok {
				assert.Equal(t, "This is my bus event to A", m)
				wg.Done()
				return
			}
		},
	)

	a.Start(nc, eb.send, map[Signature]Module{})
	<-nc

	go func() {
		eb.listen(
			[]Module{
				a,
			},
		)
	}()

	eb.send(NewBusEvent(Signature("A"), "This is my bus event to A"))

	logger.Info("waiting for event delivery")
	wg.Wait()
	logger.Info("all events delivered")
}

func TestEventBusDeliveryToTwoModules(t *testing.T) {
	logger := NewDefaultLogger(Signature("simple-logger"))
	nc := make(chan struct{}, 1)
	eb := newEventBus()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	a := NewA(
		logger,
		Signature("A"),
		func(_ func(*BusEvent), e *BusEvent) {
			// Check received event and unblock waitgroup
			if m, ok := e.Payload.(string); ok {
				assert.Equal(t, "This is my bus event to A", m)
				wg.Done()
				return
			}
		},
	)
	b := NewA(
		logger,
		Signature("B"),
		func(_ func(*BusEvent), e *BusEvent) {
			// Check received event and unblock waitgroup
			if m, ok := e.Payload.(string); ok {
				assert.Equal(t, "This is my bus event to B", m)
				wg.Done()
				return
			}
		},
	)

	a.Start(nc, eb.send, map[Signature]Module{})
	<-nc

	b.Start(nc, eb.send, map[Signature]Module{})
	<-nc

	go func() {
		eb.listen(
			[]Module{
				a, b,
			},
		)
	}()

	eb.send(NewBusEvent(Signature("A"), "This is my bus event to A"))
	eb.send(NewBusEvent(Signature("B"), "This is my bus event to B"))

	logger.Info("waiting for event delivery")
	wg.Wait()
	logger.Info("all events delivered")
}

func TestEventBusDeliveryFromOneAToB(t *testing.T) {
	logger := NewDefaultLogger(Signature("simple-logger"))
	nc := make(chan struct{}, 1)
	eb := newEventBus()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	a := NewA(
		logger,
		Signature("A"),
		func(ebs func(*BusEvent), e *BusEvent) {
			defer wg.Done()

			// Emit event for B
			ebs(NewBusEvent(Signature("B"), "FROM A TO B"))
		},
	)
	b := NewA(
		logger,
		Signature("B"),
		func(_ func(*BusEvent), e *BusEvent) {
			// Check received event and unblock waitgroup
			if m, ok := e.Payload.(string); ok {
				assert.Equal(t, "FROM A TO B", m)
				wg.Done()
				return
			}
		},
	)

	a.Start(nc, eb.send, map[Signature]Module{})
	<-nc

	b.Start(nc, eb.send, map[Signature]Module{})
	<-nc

	go func() {
		eb.listen(
			[]Module{
				a, b,
			},
		)
	}()

	eb.send(NewBusEvent(Signature("A"), "Trigger A to emit B"))

	logger.Info("waiting for event delivery")
	wg.Wait()
	logger.Info("all events delivered")
}

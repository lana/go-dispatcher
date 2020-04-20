package dispatcher

import (
	"context"
	"fmt"
	"regexp"
	"sync"
)

// EventType defines the kind of the dispatched event.
type EventType string

// String returns the converted name of the event.
func (et EventType) String() string {
	return string(et)
}

// Event defines the type and the data of an event.
type Event interface {
	Type() EventType
	Data() interface{}
}

// ListenerFunc is a function that can receive events.
type ListenerFunc func(context.Context, Event)

type Dispatcher interface {
	// On registers an event listener to events of a given type.
	On(EventType, ListenerFunc)

	// Dispatch fires an event of a given type.
	Dispatch(context.Context, Event)
}

// listeners is the internal representation of a list of listeners.
type listeners map[EventType][]*ListenerFunc

// Get retrieves all listeners from a given EventType.
func (l *listeners) Get(et EventType) []*ListenerFunc {
	var res []*ListenerFunc
	for i, fns := range *l {
		if ok, _ := regexp.MatchString(fmt.Sprintf("^%s$", i.String()), et.String()); ok {
			res = append(res, fns...)
		}
	}

	return res
}

// Add groups listeners by EventType.
func (l listeners) Add(et EventType, fn *ListenerFunc) {
	fns := l.Get(et)
	fns = append(fns, fn)
	l[et] = fns
}

// dispatcher is the internal implementation of Dispatcher.
type dispatcher struct {
	listeners *listeners
}

// New creates a new Dispatcher instance.
func New() Dispatcher {
	return &dispatcher{
		listeners: &listeners{},
	}
}

// On registers an event listener to events of a given type.
func (d *dispatcher) On(et EventType, l ListenerFunc) {
	d.listeners.Add(et, &l)
}

// Dispatch fires an event of a given type.
func (d *dispatcher) Dispatch(ctx context.Context, e Event) {
	fns := d.listeners.Get(e.Type())

	var wg sync.WaitGroup
	wg.Add(len(fns))
	for _, fn := range fns {
		x := *fn
		go func() {
			x(ctx, e)
			wg.Done()
		}()
	}
	wg.Wait()
}

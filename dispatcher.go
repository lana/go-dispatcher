package dispatcher

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrListenerNotFound = errors.New("listener not found")
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
type listeners map[EventType][]ListenerFunc

// Get retrieves all listeners from a given EventType.
func (l listeners) Get(et EventType) ([]ListenerFunc, error) {
	if fns, ok := l[et]; ok {
		return fns, nil
	}

	return nil, ErrListenerNotFound
}

// Add groups listeners by EventType.
func (l listeners) Add(et EventType, fn ListenerFunc) {
	fns, err := l.Get(et)
	if err != nil {
		fns = make([]ListenerFunc, 0)
	}

	fns = append(fns, fn)
	l[et] = fns
}

// dispatcher is the internal implementation of Dispatcher.
type dispatcher struct {
	listeners listeners
}

// New creates a new Dispatcher instance.
func New() Dispatcher {
	return &dispatcher{
		listeners: listeners{},
	}
}

// On registers an event listener to events of a given type.
func (d *dispatcher) On(et EventType, l ListenerFunc) {
	d.listeners.Add(et, l)
}

// Dispatch fires an event of a given type.
func (d *dispatcher) Dispatch(ctx context.Context, e Event) {
	fns, err := d.listeners.Get(e.Type())
	if err != nil {
		return // ignore events without listeners
	}

	var wg sync.WaitGroup
	wg.Add(len(fns))
	for _, fn := range fns {
		go func() {
			defer wg.Done()

			fn(ctx, e)
		}()
	}
	wg.Wait()
}

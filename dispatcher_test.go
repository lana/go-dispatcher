package dispatcher

import (
	"context"
	"log"
	"testing"
)

var foo EventType = "foo-event"

type fooEvent struct{}

func (e fooEvent) Type() EventType {
	return foo
}

func (e fooEvent) Data() interface{} {
	return []string{"foo", "bar", "baz"}
}

func l(ctx context.Context, e Event) {
	log.Printf("%v", e.Data())
}

func TestDispatcher(t *testing.T) {
	dispatcher := New()
	dispatcher.On(foo, l)

	e := fooEvent{}
	ctx := context.Background()

	dispatcher.Dispatch(ctx, e)
}

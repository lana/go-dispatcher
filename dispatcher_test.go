package dispatcher

import (
	"context"
	"log"
	"testing"
)

type fooEvent struct{}

func (e fooEvent) Type() EventType {
	return "event.foo.1"
}

func (e fooEvent) Data() interface{} {
	return []string{"foo", "bar", "baz"}
}

type fooListener struct {
	c int
}

func (l *fooListener) listen(ctx context.Context, e Event) {
	l.c++
	log.Printf("%v", e.Data())
}

func TestDispatcher(t *testing.T) {
	l := &fooListener{}
	dispatcher := New()
	dispatcher.On("event.*", l.listen)

	e := fooEvent{}
	ctx := context.Background()

	dispatcher.Dispatch(ctx, e)
	if l.c != 1 {
		t.Error("Event listener must be called once")
	}
}

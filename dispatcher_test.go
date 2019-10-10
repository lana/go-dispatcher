package dispatcher

import (
	"context"
	"testing"
)

var count = 0

type fooEvent struct{}

func (e fooEvent) Type() EventType {
	return "event.foo.1"
}

func (e fooEvent) Data() interface{} {
	return []string{"foo", "bar", "baz"}
}

type fooListener struct {
}

func (l *fooListener) listen(ctx context.Context, e Event) {
	count++
}

func l(ctx context.Context, e Event) {
	// do nothing
}

func TestDispatcher(t *testing.T) {
	listener := &fooListener{}
	dispatcher := New()
	dispatcher.On("*", l)
	dispatcher.On("event.*", listener.listen)

	e := fooEvent{}
	ctx := context.Background()
	runs := 10
	expected := 10

	for i := 0; i < runs; i++ {
		dispatcher.Dispatch(ctx, e)
	}

	if count != runs {
		t.Errorf("Count should be \"%d\", got \"%d\"", expected, count)
	}
}

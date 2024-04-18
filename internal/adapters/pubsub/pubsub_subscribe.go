package pubsub_adapters

import (
	"context"
	"fmt"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"
)

type handler[T pubsub.EventType, E pubsub.Event[T]] chan<- E

func (h handler[T, E]) Type() T {
	var e E
	return e.Type()
}

func (h handler[T, E]) Handle(event pubsub.Event[T]) {
	h <- event.(E)
}

func Subscribe[T pubsub.EventType, E pubsub.Event[T]](
	subs pubsub.SubscriptionsManager[T],
	preStopper module.PreStopper,
) <-chan E {
	channel := make(chan E)
	h := handler[T, E](channel)
	unSubscribe := subs.AddHandler(h)
	preStopper.PreStop(module.NewHook(
		fmt.Sprintf("event_handler_%v", h.Type()),
		func(_ context.Context) error {
			unSubscribe()
			close(channel)
			return nil
		},
	))
	return channel
}

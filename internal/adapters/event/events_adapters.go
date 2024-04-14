package event_adapters

import "github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"

type HandlerMux[T pubsub.EventType, E pubsub.Event[T]] struct {
	su
}

func Handle()

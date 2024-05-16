package pubsub_adapters

import "github.com/x0k/veterinary-clinic-backend/internal/lib/pubsub"

type NullPublisher[T pubsub.EventType] struct{}

func NewNullPublisher[T pubsub.EventType]() *NullPublisher[T] {
	return &NullPublisher[T]{}
}

func (NullPublisher[T]) Publish(_ pubsub.Event[T]) error {
	return nil
}

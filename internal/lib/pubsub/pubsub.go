package pubsub

import (
	"errors"
	"sync"
)

var ErrUnknownHandler = errors.New("unknown handler")

type EventType interface {
	~int | ~string
}

type Event[T EventType] interface {
	Type() T
}

type Handler[T EventType] interface {
	Type() T
	Handle(event Event[T])
}

type SubscriptionsManager[T EventType] interface {
	AddHandler(h Handler[T]) func()
}

type Publisher[T EventType] interface {
	Publish(event Event[T]) error
}

type PubSub[T EventType] struct {
	handlersMu sync.RWMutex
	handlers   map[T][]Handler[T]
}

func New[T EventType]() *PubSub[T] {
	return &PubSub[T]{
		handlers: map[T][]Handler[T]{},
	}
}

func (p *PubSub[T]) removeHandler(eventType T, h Handler[T]) {
	p.handlersMu.Lock()
	defer p.handlersMu.Unlock()
	for i, v := range p.handlers[eventType] {
		if v == h {
			p.handlers[eventType] = append(p.handlers[eventType][:i], p.handlers[eventType][i+1:]...)
			return
		}
	}
}

func (p *PubSub[T]) AddHandler(h Handler[T]) func() {
	p.handlersMu.Lock()
	defer p.handlersMu.Unlock()
	p.handlers[h.Type()] = append(p.handlers[h.Type()], h)
	return func() {
		p.removeHandler(h.Type(), h)
	}
}

func (p *PubSub[T]) Publish(event Event[T]) error {
	p.handlersMu.RLock()
	defer p.handlersMu.RUnlock()
	for _, h := range p.handlers[event.Type()] {
		h.Handle(event)
	}
	return nil
}

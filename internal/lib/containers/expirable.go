package containers

import (
	"context"
	"sync"
	"time"
)

type Expiable[T any] struct {
	mu     sync.RWMutex
	actual bool
	val    T
	ttl    time.Duration
	timer  *time.Timer
}

func NewExpiable[T any](ttl time.Duration) *Expiable[T] {
	return &Expiable[T]{
		ttl: ttl,
	}
}

func (e *Expiable[T]) markAsExpired() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.actual = false
}

func (e *Expiable[T]) Start(ctx context.Context) {
	e.timer = time.NewTimer(e.ttl)
	defer e.timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-e.timer.C:
			e.markAsExpired()
		}
	}
}

func (e *Expiable[T]) Get() (T, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.val, e.actual
}

func (e *Expiable[T]) Set(val T) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.val = val
	e.actual = true
	if !e.timer.Stop() {
		<-e.timer.C
	}
	e.timer.Reset(e.ttl)
}

func (e *Expiable[T]) Load(loader func() (T, error)) (T, error) {
	cached, ok := e.Get()
	if ok {
		return cached, nil
	}
	loaded, err := loader()
	if err != nil {
		return cached, err
	}
	e.Set(loaded)
	return loaded, nil
}

package adapters

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/containers"
)

type ExpirableStateContainer[S any] struct {
	name      string
	ttl       time.Duration
	mu        sync.RWMutex
	lastIndex uint64
	keys      *containers.ExpirationQueue[StateId]
	values    map[StateId]S
}

func NewExpirableStateContainer[S any](name string, seed uint64, ttl time.Duration) *ExpirableStateContainer[S] {
	return &ExpirableStateContainer[S]{
		name:      name,
		ttl:       ttl,
		lastIndex: seed,
		keys:      containers.NewExpirationQueue[StateId](),
		values:    make(map[StateId]S),
	}
}

func (c *ExpirableStateContainer[S]) Name() string {
	return c.name
}

func (c *ExpirableStateContainer[S]) flush(now time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.keys.RemoveExpired(now.Add(-c.ttl), func(key StateId) {
		delete(c.values, key)
	})
}

func (c *ExpirableStateContainer[S]) Start(ctx context.Context) error {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case now := <-ticker.C:
			c.flush(now)
		}
	}
}

func (c *ExpirableStateContainer[S]) Load(key StateId) (S, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.values[key]
	return value, ok
}

func (c *ExpirableStateContainer[S]) Save(value S) StateId {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := StateId(strconv.FormatUint(c.lastIndex, 10))
	c.lastIndex++
	c.keys.Push(key)
	c.values[key] = value
	return key
}

func (c *ExpirableStateContainer[S]) SaveByKey(key StateId, value S) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.keys.Push(key)
	c.values[key] = value
}

func (c *ExpirableStateContainer[S]) Pop(key StateId) (S, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.values[key]
	if ok {
		c.keys.Remove(key)
		delete(c.values, key)
	}
	return value, ok
}

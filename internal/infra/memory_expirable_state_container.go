package infra

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/containers"
)

type MemoryExpirableStateContainer[S any] struct {
	ttl       time.Duration
	mu        sync.RWMutex
	lastIndex uint64
	keys      *containers.ExpirationQueue[adapters.StateId]
	values    map[adapters.StateId]S
}

func NewMemoryExpirableStateContainer[S any](seed uint64, ttl time.Duration) *MemoryExpirableStateContainer[S] {
	return &MemoryExpirableStateContainer[S]{
		ttl:       ttl,
		lastIndex: seed,
		keys:      containers.NewExpirationQueue[adapters.StateId](),
		values:    make(map[adapters.StateId]S),
	}
}

func (c *MemoryExpirableStateContainer[S]) flush(now time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.keys.RemoveExpired(now.Add(-c.ttl), func(key adapters.StateId) {
		delete(c.values, key)
	})
}

func (c *MemoryExpirableStateContainer[S]) Start(ctx context.Context) error {
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

func (c *MemoryExpirableStateContainer[S]) Load(key adapters.StateId) (S, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.values[key]
	return value, ok
}

func (c *MemoryExpirableStateContainer[S]) Save(value S) adapters.StateId {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := adapters.StateId(strconv.FormatUint(c.lastIndex, 10))
	c.lastIndex++
	c.keys.Push(key)
	c.values[key] = value
	return key
}

func (c *MemoryExpirableStateContainer[S]) SaveByKey(key adapters.StateId, value S) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.keys.Push(key)
	c.values[key] = value
}

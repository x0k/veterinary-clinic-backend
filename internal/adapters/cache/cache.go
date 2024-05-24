package cache_adapters

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/cache"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
)

func StartSimpleExpirableCache[T any](
	m *module.Module,
	name string,
	c *cache.SimpleExpirableCache[T],
) *cache.SimpleExpirableCache[T] {
	m.Append(module.NewService(name, func(ctx context.Context) error {
		c.Start(ctx)
		return nil
	}))
	return c
}

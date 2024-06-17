package cache_adapters

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/cache/memory"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
)

func StartSimpleExpirableCache[T any](
	m *module.Module,
	name string,
	c *memory.SimpleExpirableCache[T],
) *memory.SimpleExpirableCache[T] {
	m.Append(module.NewService(name, func(ctx context.Context) error {
		c.Start(ctx)
		return nil
	}))
	return c
}

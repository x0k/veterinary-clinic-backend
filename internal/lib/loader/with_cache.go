package loader

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/cache"
)

func WithCache[T any](cache cache.Simple[T], loader Simple[T]) Simple[T] {
	return func(ctx context.Context) (T, error) {
		cached, ok := cache.Get(ctx)
		if ok {
			return cached, nil
		}
		loaded, err := loader(ctx)
		if err != nil {
			return loaded, err
		}
		cache.Add(ctx, loaded)
		return loaded, nil
	}
}

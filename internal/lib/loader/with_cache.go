package loader

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/cache"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

func WithCache[T any](log *logger.Logger, loader Simple[T], cache cache.Simple[T]) Simple[T] {
	return func(ctx context.Context) (T, error) {
		cached, ok := cache.Get(ctx)
		if ok {
			return cached, nil
		}
		loaded, err := loader(ctx)
		if err != nil {
			return loaded, err
		}
		if err := cache.Add(ctx, loaded); err != nil {
			log.Error(ctx, "failed to add to cache", sl.Err(err))
		}
		return loaded, nil
	}
}

package telegram_adapters

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
)

func NewController(name string, controller func(context.Context) error) module.Hook {
	return module.NewHook(name, func(ctx context.Context) error {
		return controller(ctx)
	})
}

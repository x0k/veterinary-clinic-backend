package adapters_telegram

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
)

func NewController(name string, controller func(context.Context) error) module.Service {
	return module.NewService(name, func(ctx context.Context) error {
		if err := controller(ctx); err != nil {
			return err
		}
		<-ctx.Done()
		return nil
	})
}

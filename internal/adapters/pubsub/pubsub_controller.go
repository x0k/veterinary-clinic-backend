package pubsub_adapters

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/module"
)

func NewController(name string, controller func(context.Context) error) module.Service {
	return module.NewService(name, controller)
}

package appointment

import (
	"github.com/x0k/veterinary-clinic-backend/internal/infra/module"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
)

func NewModule(log *logger.Logger) *module.Module {
	module := module.New(log, "appointment")

	return module
}

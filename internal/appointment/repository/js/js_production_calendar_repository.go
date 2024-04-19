//go:build js && wasm

package appointment_js_repository

import (
	"context"
	"syscall/js"

	"github.com/norunners/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

type ProductionCalendarRepositoryConfig struct {
	ProductionCalendar js.Value `js:"loadProductionCalendar"`
}

type ProductionCalendarRepository struct {
	cfg ProductionCalendarRepositoryConfig
}

func NewProductionCalendarRepository(
	cfg ProductionCalendarRepositoryConfig,
) *ProductionCalendarRepository {
	return &ProductionCalendarRepository{
		cfg: cfg,
	}
}

func (r *ProductionCalendarRepository) ProductionCalendar(
	ctx context.Context,
) (appointment.ProductionCalendar, error) {
	promise := r.cfg.ProductionCalendar.Invoke()
	productionCalendarJsValue, err := js_adapters.Await(ctx, promise)
	if err != nil {
		return nil, err
	}
	productionCalendar := appointment.NewProductionCalendar()
	if err := vert.ValueOf(productionCalendarJsValue).AssignTo(&productionCalendar); err != nil {
		return nil, err
	}
	return productionCalendar, nil
}

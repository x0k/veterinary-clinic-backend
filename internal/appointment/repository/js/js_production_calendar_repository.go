//go:build js && wasm

package appointment_js_repository

import (
	"context"
	"syscall/js"

	"github.com/x0k/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/js"
)

type ProductionCalendarRepositoryConfig struct {
	ProductionCalendar *js.Value `js:"loadProductionCalendar"`
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
	productionCalendarDataJsValue, err := js_adapters.Await(ctx, promise)
	if err != nil {
		return appointment.ProductionCalendar{}, err
	}
	productionCalendarDataDTO := make(
		appointment_js_adapters.ProductionCalendarDataDTO,
		js_adapters.ObjectConstructor.Call("keys", productionCalendarDataJsValue).Length(),
	)
	if err := vert.Assign(productionCalendarDataJsValue, &productionCalendarDataDTO); err != nil {
		return appointment.ProductionCalendar{}, err
	}
	productionCalendarData, err := appointment.NewProductionCalendarData(productionCalendarDataDTO)
	if err != nil {
		return appointment.ProductionCalendar{}, err
	}
	return appointment.NewProductionCalendar(productionCalendarData), nil
}

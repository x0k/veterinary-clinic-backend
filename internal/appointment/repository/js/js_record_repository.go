//go:build js && wasm

package appointment_js_repository

import (
	"context"
	"syscall/js"
	"time"

	"github.com/norunners/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/js"
)

type RecordRepositoryConfig struct {
	CreateRecord              js.Value `js:"createRecord"`
	BusyPeriods               js.Value `js:"loadBusyPeriods"`
	CustomerActiveAppointment js.Value `js:"loadCustomerActiveAppointment"`
	RemoveRecord              js.Value `js:"removeRecord"`
}

type RecordRepository struct {
	cfg RecordRepositoryConfig
}

func NewRecordRepository(
	cfg RecordRepositoryConfig,
) *RecordRepository {
	return &RecordRepository{
		cfg: cfg,
	}
}

func (r *RecordRepository) CreateRecord(
	ctx context.Context,
	rec *appointment.RecordEntity,
) error {
	promise := r.cfg.CreateRecord.Invoke(
		vert.ValueOf(appointment_js_adapters.RecordToDTO(*rec)),
	)
	recordId, err := js_adapters.Await(ctx, promise)
	if err != nil {
		return err
	}
	rec.SetId(appointment.RecordId(recordId.String()))
	return nil
}

func (r *RecordRepository) BusyPeriods(
	ctx context.Context,
	now time.Time,
) (appointment.BusyPeriods, error) {
	promise := r.cfg.BusyPeriods.Invoke(
		now.Format(time.RFC3339),
	)
	jsValue, err := js_adapters.Await(ctx, promise)
	if err != nil {
		return nil, err
	}
	busyPeriodsDTO := make([]appointment_js_adapters.TimePeriodDTO, 0)
	if err := vert.ValueOf(jsValue).AssignTo(&busyPeriodsDTO); err != nil {
		return nil, err
	}
	busyPeriods := make(appointment.BusyPeriods, len(busyPeriodsDTO))
	for i, busyPeriodDTO := range busyPeriodsDTO {
		busyPeriods[i] = appointment_js_adapters.TimePeriodFromDTO(busyPeriodDTO)
	}
	return busyPeriods, nil
}

func (r *RecordRepository) CustomerActiveAppointment(
	ctx context.Context,
	customerId appointment.CustomerId,
) (appointment.RecordEntity, error) {
	promise := r.cfg.CustomerActiveAppointment.Invoke(
		vert.ValueOf(string(customerId)),
	)
	jsValue, err := js_adapters.Await(ctx, promise)
	if err != nil {
		return appointment.RecordEntity{}, err
	}
	dto := appointment_js_adapters.RecordDTO{}
	if err := vert.ValueOf(jsValue).AssignTo(&dto); err != nil {
		return appointment.RecordEntity{}, err
	}
	return appointment_js_adapters.RecordFromDTO(dto)
}

func (r *RecordRepository) RemoveRecord(
	ctx context.Context,
	recordId appointment.RecordId,
) error {
	promise := r.cfg.RemoveRecord.Invoke(
		vert.ValueOf(string(recordId)),
	)
	_, err := js_adapters.Await(ctx, promise)
	return err
}

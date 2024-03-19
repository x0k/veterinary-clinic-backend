package make_appointment

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type datePickerPresenter[R any] interface {
	RenderDatePicker(serviceId entity.ServiceId, schedule entity.Schedule) (R, error)
}

type DatePickerUseCase[R any] struct {
	productionCalendarRepo usecase.ProductionCalendarLoader
	openingHoursRepo       usecase.OpeningHoursLoader
	busyPeriodsRepo        usecase.BusyPeriodsLoader
	workBreaksRepo         usecase.WorkBreaksLoader
	presenter              datePickerPresenter[R]
}

func NewDatePickerUseCase[R any](
	productionCalendarRepo usecase.ProductionCalendarLoader,
	openingHoursRepo usecase.OpeningHoursLoader,
	busyPeriodsRepo usecase.BusyPeriodsLoader,
	workBreaksRepo usecase.WorkBreaksLoader,
	presenter datePickerPresenter[R],
) *DatePickerUseCase[R] {
	return &DatePickerUseCase[R]{
		productionCalendarRepo: productionCalendarRepo,
		openingHoursRepo:       openingHoursRepo,
		busyPeriodsRepo:        busyPeriodsRepo,
		workBreaksRepo:         workBreaksRepo,
		presenter:              presenter,
	}
}

func (u *DatePickerUseCase[R]) DatePicker(
	ctx context.Context,
	serviceId entity.ServiceId,
	now time.Time,
	preferredDate time.Time,
) (R, error) {
	schedule, err := usecase.FetchAndCalculateSchedule(
		ctx,
		now,
		preferredDate,
		u.productionCalendarRepo,
		u.openingHoursRepo,
		u.busyPeriodsRepo,
		u.workBreaksRepo,
	)
	if err != nil {
		return *new(R), err
	}
	return u.presenter.RenderDatePicker(serviceId, schedule)
}

package usecase

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type clinicSchedulePresenter[R any] interface {
	RenderSchedule(schedule entity.Schedule) (R, error)
}

type ClinicScheduleUseCase[R any] struct {
	productionCalendarRepo ProductionCalendarLoader
	openingHoursRepo       OpeningHoursLoader
	busyPeriodsRepo        BusyPeriodsLoader
	workBreaksRepo         WorkBreaksLoader
	presenter              clinicSchedulePresenter[R]
}

func NewClinicScheduleUseCase[R any](
	productionCalendarRepo ProductionCalendarLoader,
	openingHoursRepo OpeningHoursLoader,
	busyPeriodsRepo BusyPeriodsLoader,
	workBreaksRepo WorkBreaksLoader,
	schedulePresenter clinicSchedulePresenter[R],
) *ClinicScheduleUseCase[R] {
	return &ClinicScheduleUseCase[R]{
		productionCalendarRepo: productionCalendarRepo,
		openingHoursRepo:       openingHoursRepo,
		busyPeriodsRepo:        busyPeriodsRepo,
		workBreaksRepo:         workBreaksRepo,
		presenter:              schedulePresenter,
	}
}

func (u *ClinicScheduleUseCase[R]) Schedule(ctx context.Context, now, preferredDate time.Time) (R, error) {
	schedule, err := FetchAndCalculateSchedule(
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
	return u.presenter.RenderSchedule(schedule)
}

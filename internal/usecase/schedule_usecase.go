package usecase

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type schedulePresenter[R any] interface {
	RenderSchedule(schedule shared.Schedule) (R, error)
}

type ScheduleUseCase[R any] struct {
	productionCalendarRepo ProductionCalendarLoader
	openingHoursRepo       OpeningHoursLoader
	busyPeriodsRepo        BusyPeriodsLoader
	workBreaksRepo         WorkBreaksLoader
	presenter              schedulePresenter[R]
}

func NewScheduleUseCase[R any](
	productionCalendarRepo ProductionCalendarLoader,
	openingHoursRepo OpeningHoursLoader,
	busyPeriodsRepo BusyPeriodsLoader,
	workBreaksRepo WorkBreaksLoader,
	schedulePresenter schedulePresenter[R],
) *ScheduleUseCase[R] {
	return &ScheduleUseCase[R]{
		productionCalendarRepo: productionCalendarRepo,
		openingHoursRepo:       openingHoursRepo,
		busyPeriodsRepo:        busyPeriodsRepo,
		workBreaksRepo:         workBreaksRepo,
		presenter:              schedulePresenter,
	}
}

func (u *ScheduleUseCase[R]) Schedule(ctx context.Context, now, preferredDate time.Time) (R, error) {
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

package usecase

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type ProductionCalendarRepo interface {
	ProductionCalendar(ctx context.Context) (entity.ProductionCalendar, error)
}

type OpeningHoursRepo interface {
	OpeningHours(ctx context.Context) (entity.OpeningHours, error)
}

type BusyPeriodsRepo interface {
	BusyPeriods(ctx context.Context, t time.Time) ([]entity.TimePeriod, error)
}

type WorkBreaksRepo interface {
	WorkBreaks(ctx context.Context) (entity.WorkBreaks, error)
}

type ClinicSchedulePresenter[R any] interface {
	RenderSchedule(schedule entity.Schedule) (R, error)
}

type ClinicScheduleUseCase[R any] struct {
	productionCalendarRepo ProductionCalendarRepo
	openingHoursRepo       OpeningHoursRepo
	busyPeriodsRepo        BusyPeriodsRepo
	workBreaksRepo         WorkBreaksRepo
	schedulePresenter      ClinicSchedulePresenter[R]
}

func NewClinicScheduleUseCase[R any](
	productionCalendarRepo ProductionCalendarRepo,
	openingHoursRepo OpeningHoursRepo,
	busyPeriodsRepo BusyPeriodsRepo,
	workBreaksRepo WorkBreaksRepo,
	schedulePresenter ClinicSchedulePresenter[R],
) *ClinicScheduleUseCase[R] {
	return &ClinicScheduleUseCase[R]{
		productionCalendarRepo: productionCalendarRepo,
		openingHoursRepo:       openingHoursRepo,
		busyPeriodsRepo:        busyPeriodsRepo,
		workBreaksRepo:         workBreaksRepo,
		schedulePresenter:      schedulePresenter,
	}
}

func (u *ClinicScheduleUseCase[R]) Schedule(ctx context.Context, t time.Time) (R, error) {
	productionCalendar, err := u.productionCalendarRepo.ProductionCalendar(ctx)
	if err != nil {
		return *new(R), err
	}
	now := time.Now()
	date := nextAvailableDay(productionCalendar, t)
	openingHours, err := u.openingHoursRepo.OpeningHours(ctx)
	if err != nil {
		return *new(R), err
	}
	freePeriods, err := freePeriods(productionCalendar, openingHours, now, date)
	if err != nil {
		return *new(R), err
	}
	busyPeriods, err := u.busyPeriodsRepo.BusyPeriods(ctx, date)
	if err != nil {
		return *new(R), err
	}
	allWorkBreaks, err := u.workBreaksRepo.WorkBreaks(ctx)
	if err != nil {
		return *new(R), err
	}
	workBreaks, err := workBreaks(allWorkBreaks, date)
	if err != nil {
		return *new(R), err
	}
	return u.schedulePresenter.RenderSchedule(schedule(
		productionCalendar,
		freePeriods,
		busyPeriods,
		workBreaks,
		now,
		date,
	))
}

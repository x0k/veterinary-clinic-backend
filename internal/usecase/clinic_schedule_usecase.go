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

type ClinicSchedulePresenter[R any] interface {
	RenderSchedule(schedule entity.Schedule) (R, error)
}

type ClinicScheduleUseCase[R any] struct {
	productionCalendarRepo ProductionCalendarRepo
	openingHoursRepo       OpeningHoursRepo
	schedulePresenter      ClinicSchedulePresenter[R]
}

func NewClinicScheduleUseCase[R any](
	productionCalendarRepo ProductionCalendarRepo,
	openingHoursRepo OpeningHoursRepo,
	schedulePresenter ClinicSchedulePresenter[R],
) *ClinicScheduleUseCase[R] {
	return &ClinicScheduleUseCase[R]{
		productionCalendarRepo: productionCalendarRepo,
		openingHoursRepo:       openingHoursRepo,
		schedulePresenter:      schedulePresenter,
	}
}

func (u *ClinicScheduleUseCase[R]) Schedule(ctx context.Context) (R, error) {
	productionCalendar, err := u.productionCalendarRepo.ProductionCalendar(ctx)
	if err != nil {
		return *new(R), err
	}
	now := time.Now()
	date := nextAvailableDay(productionCalendar, now)
	openingHours, err := u.openingHoursRepo.OpeningHours(ctx)
	if err != nil {
		return *new(R), err
	}
	freePeriods, err := freePeriods(productionCalendar, openingHours, now, date)
	if err != nil {
		return *new(R), err
	}

	return u.schedulePresenter.RenderSchedule(schedule)
}

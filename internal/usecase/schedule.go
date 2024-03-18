package usecase

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

func fetchAndCalculateSchedule(
	ctx context.Context,
	now time.Time,
	preferredDate time.Time,
	productionCalendarRepo productionCalendarRepo,
	openingHoursRepo openingHoursRepo,
	busyPeriodsRepo busyPeriodsRepo,
	workBreaksRepo workBreaksRepo,
) (entity.Schedule, error) {
	productionCalendar, err := productionCalendarRepo.ProductionCalendar(ctx)
	if err != nil {
		return entity.Schedule{}, err
	}
	date := entity.CalculateNextAvailableDay(productionCalendar, preferredDate)
	openingHours, err := openingHoursRepo.OpeningHours(ctx)
	if err != nil {
		return entity.Schedule{}, err
	}
	freePeriods, err := entity.CalculateFreePeriods(productionCalendar, openingHours, now, date)
	if err != nil {
		return entity.Schedule{}, err
	}
	busyPeriods, err := busyPeriodsRepo.BusyPeriods(ctx, date)
	if err != nil {
		return entity.Schedule{}, err
	}
	allWorkBreaks, err := workBreaksRepo.WorkBreaks(ctx)
	if err != nil {
		return entity.Schedule{}, err
	}
	workBreaks, err := entity.CalculateWorkBreaks(allWorkBreaks, date)
	if err != nil {
		return entity.Schedule{}, err
	}
	return entity.CalculateSchedule(productionCalendar, freePeriods, busyPeriods, workBreaks, now, date), nil
}

package usecase

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

func FetchAndCalculateSchedule(
	ctx context.Context,
	now time.Time,
	preferredDate time.Time,
	productionCalendarRepo ProductionCalendarLoader,
	openingHoursRepo OpeningHoursLoader,
	busyPeriodsRepo BusyPeriodsLoader,
	workBreaksRepo WorkBreaksLoader,
) (shared.Schedule, error) {
	productionCalendar, err := productionCalendarRepo.ProductionCalendar(ctx)
	if err != nil {
		return shared.Schedule{}, err
	}
	date := shared.CalculateNextAvailableDay(productionCalendar, preferredDate)
	openingHours, err := openingHoursRepo.OpeningHours(ctx)
	if err != nil {
		return shared.Schedule{}, err
	}
	freePeriods, err := shared.CalculateFreePeriods(productionCalendar, openingHours, now, date)
	if err != nil {
		return shared.Schedule{}, err
	}
	busyPeriods, err := busyPeriodsRepo.BusyPeriods(ctx, date)
	if err != nil {
		return shared.Schedule{}, err
	}
	allWorkBreaks, err := workBreaksRepo.WorkBreaks(ctx)
	if err != nil {
		return shared.Schedule{}, err
	}
	workBreaks, err := shared.CalculateWorkBreaks(allWorkBreaks, date)
	if err != nil {
		return shared.Schedule{}, err
	}
	return shared.CalculateSchedule(productionCalendar, freePeriods, busyPeriods, workBreaks, now, date), nil
}

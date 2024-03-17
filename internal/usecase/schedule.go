package usecase

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

func nextAvailableDay(productionCalendar entity.ProductionCalendar, from time.Time) time.Time {
	return entity.NewNextAvailableDayCalculator(productionCalendar).Calculate(from)
}

func prevAvailableDay(productionCalendar entity.ProductionCalendar, now time.Time, from time.Time) *time.Time {
	return entity.NewPrevAvailableDayCalculator(productionCalendar, now).Calculate(from)
}

func schedule(
	productionCalendar entity.ProductionCalendar,
	freePeriods entity.FreePeriods,
	busyPeriods entity.BusyPeriods,
	workBreaks entity.WorkBreaks,
	now time.Time,
	date time.Time,
) entity.Schedule {
	schedulePeriods := entity.CalculateSchedulePeriods(freePeriods, busyPeriods, workBreaks)
	next := nextAvailableDay(productionCalendar, date)
	prev := prevAvailableDay(productionCalendar, now, date)
	return entity.NewSchedule(date, schedulePeriods).SetDates(&next, prev)
}

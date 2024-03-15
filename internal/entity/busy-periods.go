package entity

import (
	"slices"
	"time"
)

type BusyPeriods []DateTimePeriod

type BusyPeriodsCalculator struct {
	busyPeriods BusyPeriods
}

func NewBusyPeriodsCalculator(
	busyPeriods BusyPeriods,
) *BusyPeriodsCalculator {
	return &BusyPeriodsCalculator{
		busyPeriods: busyPeriods,
	}
}

func (c *BusyPeriodsCalculator) Calculate(
	t time.Time,
) BusyPeriods {
	dayPeriod := DateToDayTimePeriod(GoTimeToDate(t))
	periods := make(BusyPeriods, 0, len(c.busyPeriods))
	for _, bp := range c.busyPeriods {
		period := DateTimePeriodApi.IntersectPeriods(bp, dayPeriod)
		if DateTimePeriodApi.IsValidPeriod(period) {
			periods = append(periods, period)
		}
	}
	slices.SortFunc(periods, DateTimePeriodApi.ComparePeriods)
	return periods
}

package entity

import (
	"errors"
	"fmt"
	"time"
)

var ErrUnknownDayType = errors.New("unknown day type")

type FreePeriods []TimePeriod

type FreePeriodsCalculator struct {
	openingHours       OpeningHours
	productionCalendar ProductionCalendar
	currentDateTime    DateTime
}

func NewFreePeriodsCalculator(
	openingHours OpeningHours,
	productionCalendar ProductionCalendar,
	currentDateTime DateTime,
) *FreePeriodsCalculator {
	return &FreePeriodsCalculator{
		openingHours:       openingHours,
		productionCalendar: productionCalendar,
		currentDateTime:    currentDateTime,
	}
}

func (c *FreePeriodsCalculator) getOpeningHours(t time.Time) DayTimePeriods {
	periods := make([]TimePeriod, 0, 1)
	if period, ok := c.openingHours[t.Weekday()]; ok {
		periods = append(periods, period)
	}
	return DayTimePeriods{
		Date:    GoTimeToDate(t),
		Periods: periods,
	}
}

func (c *FreePeriodsCalculator) applyCurrentDateTime(
	data DayTimePeriods,
) DayTimePeriods {
	compareResult := CompareDate(data.Date, c.currentDateTime.Date)
	if compareResult < 0 {
		return DayTimePeriods{
			Date:    data.Date,
			Periods: []TimePeriod{},
		}
	}
	if compareResult > 0 {
		return data
	}
	period := TimePeriod{
		Start: Time{
			Hours:   0,
			Minutes: 0,
		},
		End: c.currentDateTime.Time,
	}
	periods := make([]TimePeriod, 0, len(data.Periods))
	for _, p := range data.Periods {
		periods = append(periods, TimePeriodApi.SubtractPeriods(p, period)...)
	}
	return DayTimePeriods{
		Date:    data.Date,
		Periods: TimePeriodApi.SortAndUnitePeriods(periods),
	}
}

func (c *FreePeriodsCalculator) applyProductionCalendar(data DayTimePeriods) (DayTimePeriods, error) {
	dayType, ok := c.productionCalendar[GoTimeToJsonDate(
		DateToGoTime(data.Date),
	)]
	if !ok {
		return data, nil
	}
	switch dayType {
	case Weekend:
		return DayTimePeriods{
			Date:    data.Date,
			Periods: []TimePeriod{},
		}, nil
	case Holiday:
		return DayTimePeriods{
			Date:    data.Date,
			Periods: []TimePeriod{},
		}, nil
	case PreHoliday:
		if len(data.Periods) < 1 {
			return data, nil
		}
		periods := TimePeriodApi.SortAndUnitePeriods(data.Periods)
		minutesToReduce := DurationInMinutes(-60)
		i := len(periods)
		var reducedLastPeriod TimePeriod
		for minutesToReduce < 0 && i > 0 {
			i--
			lastPeriod := periods[i]
			shift := MakeTimeShifter(Time{
				Minutes: int(minutesToReduce),
			})
			reducedLastPeriod = TimePeriod{
				Start: lastPeriod.Start,
				End:   shift(lastPeriod.End),
			}
			minutesToReduce = TimePeriodDurationInMinutes(reducedLastPeriod)
		}
		if minutesToReduce < 0 {
			return DayTimePeriods{
				Date:    data.Date,
				Periods: nil,
			}, nil
		}
		return DayTimePeriods{
			Date:    data.Date,
			Periods: append(periods[:i], reducedLastPeriod),
		}, nil
	}
	return data, fmt.Errorf("%w: %d", ErrUnknownDayType, dayType)
}

func (c *FreePeriodsCalculator) Calculate(
	t time.Time,
) (FreePeriods, error) {
	data, err := c.applyProductionCalendar(
		c.applyCurrentDateTime(
			c.getOpeningHours(t),
		),
	)
	if err != nil {
		return nil, err
	}
	return data.Periods, nil
}

func CalculateFreePeriods(
	productionCalendar ProductionCalendar,
	openingHours OpeningHours,
	now time.Time,
	forDate time.Time,
) (FreePeriods, error) {
	return NewFreePeriodsCalculator(
		openingHours,
		productionCalendar,
		GoTimeToDateTime(now),
	).Calculate(forDate)
}

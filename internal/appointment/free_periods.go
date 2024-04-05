package appointment

import (
	"errors"
	"fmt"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

var ErrUnknownDayType = errors.New("unknown day type")

type FreePeriods []entity.TimePeriod

type FreePeriodsCalculator struct {
	openingHours       WorkingHours
	productionCalendar ProductionCalendar
	currentDateTime    entity.DateTime
}

func NewFreePeriodsCalculator(
	openingHours WorkingHours,
	productionCalendar ProductionCalendar,
	currentDateTime entity.DateTime,
) *FreePeriodsCalculator {
	return &FreePeriodsCalculator{
		openingHours:       openingHours,
		productionCalendar: productionCalendar,
		currentDateTime:    currentDateTime,
	}
}

func (c *FreePeriodsCalculator) getOpeningHours(t time.Time) entity.DayTimePeriods {
	periods := make([]entity.TimePeriod, 0, 1)
	if period, ok := c.openingHours[t.Weekday()]; ok {
		periods = append(periods, period)
	}
	return entity.DayTimePeriods{
		Date:    entity.GoTimeToDate(t),
		Periods: periods,
	}
}

func (c *FreePeriodsCalculator) applyCurrentDateTime(
	data entity.DayTimePeriods,
) entity.DayTimePeriods {
	compareResult := entity.CompareDate(data.Date, c.currentDateTime.Date)
	if compareResult < 0 {
		return entity.DayTimePeriods{
			Date:    data.Date,
			Periods: []entity.TimePeriod{},
		}
	}
	if compareResult > 0 {
		return data
	}
	period := entity.TimePeriod{
		Start: entity.Time{
			Hours:   0,
			Minutes: 0,
		},
		End: c.currentDateTime.Time,
	}
	periods := make([]entity.TimePeriod, 0, len(data.Periods))
	for _, p := range data.Periods {
		periods = append(periods, entity.TimePeriodApi.SubtractPeriods(p, period)...)
	}
	return entity.DayTimePeriods{
		Date:    data.Date,
		Periods: entity.TimePeriodApi.SortAndUnitePeriods(periods),
	}
}

func (c *FreePeriodsCalculator) applyProductionCalendar(data entity.DayTimePeriods) (entity.DayTimePeriods, error) {
	dayType, ok := c.productionCalendar[entity.GoTimeToJsonDate(
		entity.DateToGoTime(data.Date),
	)]
	if !ok {
		return data, nil
	}
	switch dayType {
	case Weekend:
		return entity.DayTimePeriods{
			Date:    data.Date,
			Periods: []entity.TimePeriod{},
		}, nil
	case Holiday:
		return entity.DayTimePeriods{
			Date:    data.Date,
			Periods: []entity.TimePeriod{},
		}, nil
	case PreHoliday:
		if len(data.Periods) < 1 {
			return data, nil
		}
		periods := entity.TimePeriodApi.SortAndUnitePeriods(data.Periods)
		minutesToReduce := entity.DurationInMinutes(-60)
		i := len(periods)
		var reducedLastPeriod entity.TimePeriod
		for minutesToReduce < 0 && i > 0 {
			i--
			lastPeriod := periods[i]
			shift := entity.MakeTimeShifter(entity.Time{
				Minutes: int(minutesToReduce),
			})
			reducedLastPeriod = entity.TimePeriod{
				Start: lastPeriod.Start,
				End:   shift(lastPeriod.End),
			}
			minutesToReduce = entity.TimePeriodDurationInMinutes(reducedLastPeriod)
		}
		if minutesToReduce < 0 {
			return entity.DayTimePeriods{
				Date:    data.Date,
				Periods: nil,
			}, nil
		}
		return entity.DayTimePeriods{
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
	openingHours WorkingHours,
	now time.Time,
	forDate time.Time,
) (FreePeriods, error) {
	return NewFreePeriodsCalculator(
		openingHours,
		productionCalendar,
		entity.GoTimeToDateTime(now),
	).Calculate(forDate)
}

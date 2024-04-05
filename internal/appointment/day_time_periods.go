package appointment

import (
	"errors"
	"fmt"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

var ErrUnknownDayType = errors.New("unknown day type")

// TODO: Convert to a domain object
type DayTimePeriods struct {
	Date    entity.Date
	Periods []entity.TimePeriod
}

func NewDayTimePeriods(
	date entity.Date,
	periods []entity.TimePeriod,
) DayTimePeriods {
	return DayTimePeriods{
		Date:    date,
		Periods: periods,
	}
}

func (data DayTimePeriods) OmitPast(now entity.DateTime) DayTimePeriods {
	compareResult := entity.CompareDate(data.Date, now.Date)
	if compareResult < 0 {
		return NewDayTimePeriods(
			data.Date,
			[]entity.TimePeriod{},
		)
	}
	if compareResult > 0 {
		return data
	}
	period := entity.TimePeriod{
		Start: entity.Time{
			Hours:   0,
			Minutes: 0,
		},
		End: now.Time,
	}
	periods := make([]entity.TimePeriod, 0, len(data.Periods))
	for _, p := range data.Periods {
		periods = append(periods, entity.TimePeriodApi.SubtractPeriods(p, period)...)
	}
	return NewDayTimePeriods(
		data.Date,
		entity.TimePeriodApi.SortAndUnitePeriods(periods),
	)
}

func (data DayTimePeriods) ConsiderProductionCalendar(cal ProductionCalendar) (DayTimePeriods, error) {
	dayType, ok := cal[entity.GoTimeToJsonDate(
		entity.DateToGoTime(data.Date),
	)]
	if !ok {
		return data, nil
	}
	switch dayType {
	case Weekend:
		return NewDayTimePeriods(
			data.Date,
			[]entity.TimePeriod{},
		), nil
	case Holiday:
		return NewDayTimePeriods(
			data.Date,
			[]entity.TimePeriod{},
		), nil
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
			return NewDayTimePeriods(
				data.Date,
				nil,
			), nil
		}
		return NewDayTimePeriods(
			data.Date,
			append(periods[:i], reducedLastPeriod),
		), nil
	}
	return data, fmt.Errorf("%w: %d", ErrUnknownDayType, dayType)
}

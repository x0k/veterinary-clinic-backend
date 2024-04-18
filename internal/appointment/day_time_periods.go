package appointment

import (
	"errors"
	"fmt"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

var ErrUnknownDayType = errors.New("unknown day type")

// TODO: Convert to a domain object
type DayTimePeriods struct {
	Date    shared.Date
	Periods []shared.TimePeriod
}

func NewDayTimePeriods(
	date shared.Date,
	periods []shared.TimePeriod,
) DayTimePeriods {
	return DayTimePeriods{
		Date:    date,
		Periods: periods,
	}
}

func (data DayTimePeriods) OmitPast(now shared.DateTime) DayTimePeriods {
	compareResult := shared.CompareDate(data.Date, now.Date)
	if compareResult < 0 {
		return NewDayTimePeriods(
			data.Date,
			[]shared.TimePeriod{},
		)
	}
	if compareResult > 0 {
		return data
	}
	period := shared.TimePeriod{
		Start: shared.Time{
			Hours:   0,
			Minutes: 0,
		},
		End: now.Time,
	}
	periods := make([]shared.TimePeriod, 0, len(data.Periods))
	for _, p := range data.Periods {
		periods = append(periods, shared.TimePeriodApi.SubtractPeriods(p, period)...)
	}
	return NewDayTimePeriods(
		data.Date,
		shared.TimePeriodApi.SortAndUnitePeriods(periods),
	)
}

func (data DayTimePeriods) ConsiderProductionCalendar(cal ProductionCalendar) (DayTimePeriods, error) {
	dayType, ok := cal[shared.GoTimeToJsonDate(
		shared.DateToGoTime(data.Date),
	)]
	if !ok {
		return data, nil
	}
	switch dayType {
	case Weekend:
		return NewDayTimePeriods(
			data.Date,
			[]shared.TimePeriod{},
		), nil
	case Holiday:
		return NewDayTimePeriods(
			data.Date,
			[]shared.TimePeriod{},
		), nil
	case PreHoliday:
		if len(data.Periods) < 1 {
			return data, nil
		}
		periods := shared.TimePeriodApi.SortAndUnitePeriods(data.Periods)
		minutesToReduce := shared.DurationInMinutes(-60)
		i := len(periods)
		var reducedLastPeriod shared.TimePeriod
		for minutesToReduce < 0 && i > 0 {
			i--
			lastPeriod := periods[i]
			shift := shared.MakeTimeShifter(shared.Time{
				Minutes: int(minutesToReduce),
			})
			reducedLastPeriod = shared.TimePeriod{
				Start: lastPeriod.Start,
				End:   shift(lastPeriod.End),
			}
			minutesToReduce = shared.TimePeriodDurationInMinutes(reducedLastPeriod)
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

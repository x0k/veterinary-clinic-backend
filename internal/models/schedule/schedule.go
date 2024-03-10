package schedule

import (
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/models/date"
)

var ErrUnknownDayType = errors.New("unknown day type")

type DayType int

const (
	Weekend    DayType = 1
	Holiday    DayType = 2
	PreHoliday DayType = 3
)

type ProductionCalendar map[date.JsonDate]DayType

type OpeningHours map[time.Weekday]date.TimePeriod

type WorkBreakId string

type WorkBreak struct {
	Id              WorkBreakId
	Title           string
	MatchExpression string
	Period          date.TimePeriod
	DateFormat      string
}

type WorkBreaks []WorkBreak

type BusyPeriods []date.DateTimePeriod

type DateTimePeriods struct {
	Date    date.Date
	Periods []date.TimePeriod
}

func ProductionCalendarWithoutSaturdayWeekend(
	cal ProductionCalendar,
) ProductionCalendar {
	cloned := maps.Clone(cal)
	for d, dt := range cal {
		if dt != Weekend {
			continue
		}
		t, err := date.JsonDateToGoTime(d)
		if err != nil || t.Weekday() == time.Saturday {
			delete(cloned, d)
		}
	}
	return cloned
}

func GoTimeToDateTimePeriod(t time.Time) date.DateTimePeriod {
	d := date.GoTimeToDate(t)
	return date.DateTimePeriod{
		Start: date.DateTime{
			Date: d,
			Time: date.Time{
				Hours:   0,
				Minutes: 0,
			},
		},
		End: date.DateTime{
			Date: d,
			Time: date.Time{
				Hours:   23,
				Minutes: 59,
			},
		},
	}
}

type FreePeriodsCalculator struct {
	openingHours       OpeningHours
	productionCalendar ProductionCalendar
	currentDateTime    date.DateTime
}

func (c *FreePeriodsCalculator) getOpeningHours(t time.Time) DateTimePeriods {
	periods := make([]date.TimePeriod, 0, 1)
	if period, ok := c.openingHours[t.Weekday()]; ok {
		periods = append(periods, period)
	}
	return DateTimePeriods{
		Date:    date.GoTimeToDate(t),
		Periods: periods,
	}
}

func (c *FreePeriodsCalculator) applyCurrentDateTime(
	data DateTimePeriods,
) DateTimePeriods {
	compareResult := date.CompareDate(data.Date, c.currentDateTime.Date)
	if compareResult < 0 {
		return DateTimePeriods{
			Date:    data.Date,
			Periods: []date.TimePeriod{},
		}
	}
	if compareResult > 0 {
		return data
	}
	period := date.TimePeriod{
		Start: date.Time{
			Hours:   0,
			Minutes: 0,
		},
		End: c.currentDateTime.Time,
	}
	periods := make([]date.TimePeriod, 0, len(data.Periods))
	for _, p := range data.Periods {
		periods = append(periods, date.TimePeriodApi.SubtractPeriods(p, period)...)
	}
	return DateTimePeriods{
		Date:    data.Date,
		Periods: date.TimePeriodApi.SortAndUnitePeriods(periods),
	}
}

func (c *FreePeriodsCalculator) applyProductionCalendar(data DateTimePeriods) (DateTimePeriods, error) {
	dayType, ok := c.productionCalendar[date.GoTimeToJsonDate(
		date.DateToGoTime(data.Date),
	)]
	if !ok {
		return data, nil
	}
	switch dayType {
	case Weekend:
		return DateTimePeriods{
			Date:    data.Date,
			Periods: []date.TimePeriod{},
		}, nil
	case Holiday:
		return DateTimePeriods{
			Date:    data.Date,
			Periods: []date.TimePeriod{},
		}, nil
	case PreHoliday:
		if len(data.Periods) < 1 {
			return data, nil
		}
		periods := date.TimePeriodApi.SortAndUnitePeriods(data.Periods)
		minutesToReduce := -60
		i := len(periods) - 1
		var reducedLastPeriod date.TimePeriod
		for minutesToReduce < 0 && i >= 0 {
			lastPeriod := periods[i]
			i--
			shift := date.MakeTimeShifter(date.Time{
				Minutes: minutesToReduce,
			})
			reducedLastPeriod = date.TimePeriod{
				Start: lastPeriod.Start,
				End:   shift(lastPeriod.End),
			}
			minutesToReduce = date.TimePeriodDurationInMinutes(reducedLastPeriod)
		}
		if minutesToReduce < 0 {
			return DateTimePeriods{
				Date:    data.Date,
				Periods: nil,
			}, nil
		}
		return DateTimePeriods{
			Date:    data.Date,
			Periods: append(periods[:i+1], reducedLastPeriod),
		}, nil
	}
	return data, fmt.Errorf("%w: %d", ErrUnknownDayType, dayType)
}

func (c *FreePeriodsCalculator) Calculate(
	t time.Time,
) ([]date.TimePeriod, error) {
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

type FreeTimePeriodsWithDurationCalculator struct {
	durationInMinutes   int
	durationShift       func(date.Time) date.Time
	sampleRateInMinutes int
	sampleRateShift     func(date.Time) date.Time
}

func NewFreeTimePeriodsWithDurationCalculator(
	durationInMinutes int,
	sampleRateInMinutes int,
) *FreeTimePeriodsWithDurationCalculator {
	durationShift := date.MakeTimeShifter(date.Time{
		Minutes: durationInMinutes,
	})
	sampleRateShift := date.MakeTimeShifter(date.Time{
		Minutes: sampleRateInMinutes,
	})
	return &FreeTimePeriodsWithDurationCalculator{
		durationInMinutes:   durationInMinutes,
		durationShift:       durationShift,
		sampleRateInMinutes: sampleRateInMinutes,
		sampleRateShift:     sampleRateShift,
	}
}

func (c *FreeTimePeriodsWithDurationCalculator) Calculate(period date.TimePeriod) []date.TimePeriod {
	rest := period.Start.Minutes % c.sampleRateInMinutes
	if rest != 0 {
		return c.Calculate(date.TimePeriod{
			Start: date.MakeTimeShifter(date.Time{
				Minutes: c.sampleRateInMinutes - rest,
			})(period.Start),
			End: period.End,
		})
	}
	periodDuration := date.TimePeriodDurationInMinutes(period)
	if periodDuration < c.durationInMinutes {
		return []date.TimePeriod{period}
	}
	periods := make([]date.TimePeriod, 1)
	periods[0] = date.TimePeriod{
		Start: period.Start,
		End:   c.durationShift(period.Start),
	}
	for _, p := range date.TimePeriodApi.SubtractPeriods(period, date.TimePeriod{
		Start: period.Start,
		End:   c.sampleRateShift(period.Start),
	}) {
		periods = append(periods, c.Calculate(p)...)
	}
	return periods
}

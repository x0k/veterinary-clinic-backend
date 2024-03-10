package entity

import (
	"errors"
	"fmt"
	"maps"
	"regexp"
	"time"
)

var ErrUnknownDayType = errors.New("unknown day type")
var ErrUnsupportedDateFormat = errors.New("unsupported date format")
var ErrFailedToCompileMatchExpression = errors.New("failed to compile match expression")

type DayType int

const (
	Weekend    DayType = 1
	Holiday    DayType = 2
	PreHoliday DayType = 3
)

func IsNonWorkingDayType(dayType DayType) bool {
	return dayType == Holiday || dayType == PreHoliday
}

type ProductionCalendar map[JsonDate]DayType

type OpeningHours map[time.Weekday]TimePeriod

type WorkBreakId string

type WorkBreak struct {
	Id              WorkBreakId
	Title           string
	MatchExpression string
	Period          TimePeriod
	DateFormat      string
}

type WorkBreaks []WorkBreak

type BusyPeriods []DateTimePeriod

type DateTimePeriods struct {
	Date    Date
	Periods []TimePeriod
}

func ProductionCalendarWithoutSaturdayWeekend(
	cal ProductionCalendar,
) ProductionCalendar {
	cloned := maps.Clone(cal)
	for d, dt := range cal {
		if dt != Weekend {
			continue
		}
		t, err := JsonDateToGoTime(d)
		if err != nil || t.Weekday() == time.Saturday {
			delete(cloned, d)
		}
	}
	return cloned
}

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

func (c *FreePeriodsCalculator) getOpeningHours(t time.Time) DateTimePeriods {
	periods := make([]TimePeriod, 0, 1)
	if period, ok := c.openingHours[t.Weekday()]; ok {
		periods = append(periods, period)
	}
	return DateTimePeriods{
		Date:    GoTimeToDate(t),
		Periods: periods,
	}
}

func (c *FreePeriodsCalculator) applyCurrentDateTime(
	data DateTimePeriods,
) DateTimePeriods {
	compareResult := CompareDate(data.Date, c.currentDateTime.Date)
	if compareResult < 0 {
		return DateTimePeriods{
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
	return DateTimePeriods{
		Date:    data.Date,
		Periods: TimePeriodApi.SortAndUnitePeriods(periods),
	}
}

func (c *FreePeriodsCalculator) applyProductionCalendar(data DateTimePeriods) (DateTimePeriods, error) {
	dayType, ok := c.productionCalendar[GoTimeToJsonDate(
		DateToGoTime(data.Date),
	)]
	if !ok {
		return data, nil
	}
	switch dayType {
	case Weekend:
		return DateTimePeriods{
			Date:    data.Date,
			Periods: []TimePeriod{},
		}, nil
	case Holiday:
		return DateTimePeriods{
			Date:    data.Date,
			Periods: []TimePeriod{},
		}, nil
	case PreHoliday:
		if len(data.Periods) < 1 {
			return data, nil
		}
		periods := TimePeriodApi.SortAndUnitePeriods(data.Periods)
		minutesToReduce := -60
		i := len(periods)
		var reducedLastPeriod TimePeriod
		for minutesToReduce < 0 && i > 0 {
			i--
			lastPeriod := periods[i]
			shift := MakeTimeShifter(Time{
				Minutes: minutesToReduce,
			})
			reducedLastPeriod = TimePeriod{
				Start: lastPeriod.Start,
				End:   shift(lastPeriod.End),
			}
			minutesToReduce = TimePeriodDurationInMinutes(reducedLastPeriod)
		}
		if minutesToReduce < 0 {
			return DateTimePeriods{
				Date:    data.Date,
				Periods: nil,
			}, nil
		}
		return DateTimePeriods{
			Date:    data.Date,
			Periods: append(periods[:i], reducedLastPeriod),
		}, nil
	}
	return data, fmt.Errorf("%w: %d", ErrUnknownDayType, dayType)
}

func (c *FreePeriodsCalculator) Calculate(
	t time.Time,
) ([]TimePeriod, error) {
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
	durationShift       func(Time) Time
	sampleRateInMinutes int
	sampleRateShift     func(Time) Time
}

func NewFreeTimePeriodsWithDurationCalculator(
	durationInMinutes int,
	sampleRateInMinutes int,
) *FreeTimePeriodsWithDurationCalculator {
	durationShift := MakeTimeShifter(Time{
		Minutes: durationInMinutes,
	})
	sampleRateShift := MakeTimeShifter(Time{
		Minutes: sampleRateInMinutes,
	})
	return &FreeTimePeriodsWithDurationCalculator{
		durationInMinutes:   durationInMinutes,
		durationShift:       durationShift,
		sampleRateInMinutes: sampleRateInMinutes,
		sampleRateShift:     sampleRateShift,
	}
}

func (c *FreeTimePeriodsWithDurationCalculator) Calculate(period TimePeriod) []TimePeriod {
	rest := period.Start.Minutes % c.sampleRateInMinutes
	if rest != 0 {
		return c.Calculate(TimePeriod{
			Start: MakeTimeShifter(Time{
				Minutes: c.sampleRateInMinutes - rest,
			})(period.Start),
			End: period.End,
		})
	}
	periodDuration := TimePeriodDurationInMinutes(period)
	if periodDuration < c.durationInMinutes {
		return []TimePeriod{period}
	}
	periods := make([]TimePeriod, 1)
	periods[0] = TimePeriod{
		Start: period.Start,
		End:   c.durationShift(period.Start),
	}
	for _, p := range TimePeriodApi.SubtractPeriods(period, TimePeriod{
		Start: period.Start,
		End:   c.sampleRateShift(period.Start),
	}) {
		periods = append(periods, c.Calculate(p)...)
	}
	return periods
}

type NextAvailableDayCalculator struct {
	productionCalendar ProductionCalendar
}

func NewNextAvailableDayCalculator(
	productionCalendar ProductionCalendar,
) *NextAvailableDayCalculator {
	return &NextAvailableDayCalculator{
		productionCalendar: productionCalendar,
	}
}

func (c *NextAvailableDayCalculator) Calculate(
	today time.Time,
) time.Time {
	nextDay := today
	for {
		nextDay = nextDay.AddDate(0, 0, 1)
		nextDayJson := GoTimeToJsonDate(nextDay)
		if dayType, ok := c.productionCalendar[nextDayJson]; !ok || !IsNonWorkingDayType(dayType) {
			return nextDay
		}
	}
}

type WorkBreaksCalculator struct {
	workBreaks WorkBreaks
}

func NewWorkBreaksCalculator(
	workBreaks WorkBreaks,
) *WorkBreaksCalculator {
	return &WorkBreaksCalculator{
		workBreaks: workBreaks,
	}
}

const date_format = "2006-01-02T15:04:05"

func (c *WorkBreaksCalculator) Calculate(
	t time.Time,
) (WorkBreaks, error) {
	date := fmt.Sprintf("%d %s", t.Weekday(), t.Format(date_format))
	breaks := make(WorkBreaks, 0, len(c.workBreaks))
	for _, wb := range c.workBreaks {
		if wb.DateFormat != "" {
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedDateFormat, wb.DateFormat)
		}
		expr, err := regexp.Compile(wb.MatchExpression)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrFailedToCompileMatchExpression, wb.MatchExpression)
		}
		if expr.MatchString(date) {
			breaks = append(breaks, wb)
		}
	}
	return breaks, nil
}

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
	dayPeriod := DateToDateTimePeriod(GoTimeToDate(t))
	periods := make(BusyPeriods, 0, len(c.busyPeriods))
	for _, bp := range c.busyPeriods {
		period := DateTimePeriodApi.IntersectPeriods(bp, dayPeriod)
		if DateTimePeriodApi.IsValidPeriod(period) {
			periods = append(periods, period)
		}
	}
	return periods
}

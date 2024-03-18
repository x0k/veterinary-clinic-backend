package entity

import (
	"slices"
	"time"
)

type TimePeriodType int

const (
	FreePeriod TimePeriodType = iota
	BusyPeriod
)

type TitledTimePeriod struct {
	TimePeriod
	Type  TimePeriodType
	Title string
}

type SchedulePeriods []TitledTimePeriod

type Schedule struct {
	Date     time.Time
	Periods  SchedulePeriods
	NextDate *time.Time
	PrevDate *time.Time
}

func NewSchedule(t time.Time, periods []TitledTimePeriod) Schedule {
	return Schedule{
		Date:    t,
		Periods: periods,
	}
}

func (s Schedule) SetDates(next *time.Time, prev *time.Time) Schedule {
	s.NextDate = next
	s.PrevDate = prev
	return s
}

func CalculateSchedulePeriods(
	freePeriods FreePeriods,
	busyPeriods BusyPeriods,
	workBreaks CalculatedWorkBreaks,
) SchedulePeriods {
	allBusyPeriods := make([]TimePeriod, len(busyPeriods), len(busyPeriods)+len(workBreaks))
	copy(allBusyPeriods, busyPeriods)
	for _, wb := range workBreaks {
		allBusyPeriods = append(allBusyPeriods, wb.Period)
	}

	actualFreePeriods := TimePeriodApi.SortAndUnitePeriods(
		TimePeriodApi.SubtractPeriodsFromPeriods(
			freePeriods,
			allBusyPeriods,
		),
	)

	schedule := make([]TitledTimePeriod, 0, len(actualFreePeriods)+len(allBusyPeriods))
	for _, p := range actualFreePeriods {
		schedule = append(schedule, TitledTimePeriod{
			TimePeriod: p,
			Type:       FreePeriod,
			Title:      "Свободно",
		})
	}
	for _, p := range busyPeriods {
		schedule = append(schedule, TitledTimePeriod{
			TimePeriod: p,
			Type:       BusyPeriod,
			Title:      "Занято",
		})
	}
	for _, p := range workBreaks {
		schedule = append(schedule, TitledTimePeriod{
			TimePeriod: p.Period,
			Type:       BusyPeriod,
			Title:      p.Title,
		})
	}
	slices.SortFunc(schedule, func(a, b TitledTimePeriod) int {
		return TimePeriodApi.ComparePeriods(a.TimePeriod, b.TimePeriod)
	})
	return schedule
}

func CalculateNextAvailableDay(productionCalendar ProductionCalendar, from time.Time) time.Time {
	return NewNextAvailableDayCalculator(productionCalendar).Calculate(from)
}

func CalculatePrevAvailableDay(productionCalendar ProductionCalendar, now time.Time, from time.Time) *time.Time {
	return NewPrevAvailableDayCalculator(productionCalendar, now).Calculate(from)
}

func CalculateSchedule(
	productionCalendar ProductionCalendar,
	freePeriods FreePeriods,
	busyPeriods BusyPeriods,
	workBreaks CalculatedWorkBreaks,
	now time.Time,
	date time.Time,
) Schedule {
	schedulePeriods := CalculateSchedulePeriods(freePeriods, busyPeriods, workBreaks)
	next := CalculateNextAvailableDay(productionCalendar, date.AddDate(0, 0, 1))
	prev := CalculatePrevAvailableDay(productionCalendar, now, date.AddDate(0, 0, -1))
	return NewSchedule(date, schedulePeriods).SetDates(&next, prev)
}

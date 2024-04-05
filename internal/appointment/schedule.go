package appointment

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type SchedulePeriods TitledTimePeriods

type Schedule struct {
	Date     time.Time
	Periods  SchedulePeriods
	NextDate time.Time
	PrevDate time.Time
}

func CalculateSchedulePeriods(
	freePeriods FreePeriods,
	busyPeriods BusyPeriods,
	workBreaks CalculatedWorkBreaks,
) SchedulePeriods {
	allBusyPeriods := make([]entity.TimePeriod, len(busyPeriods), len(busyPeriods)+len(workBreaks))
	copy(allBusyPeriods, busyPeriods)
	for _, wb := range workBreaks {
		allBusyPeriods = append(allBusyPeriods, wb.Period)
	}

	actualFreePeriods := entity.TimePeriodApi.SortAndUnitePeriods(
		entity.TimePeriodApi.SubtractPeriodsFromPeriods(
			freePeriods,
			allBusyPeriods,
		),
	)

	schedule := make(TitledTimePeriods, 0, len(actualFreePeriods)+len(allBusyPeriods))
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
	return SchedulePeriods(schedule.SortAndFlat())
}

func NewSchedule(
	productionCalendar ProductionCalendar,
	freePeriods FreePeriods,
	busyPeriods BusyPeriods,
	workBreaks CalculatedWorkBreaks,
	now time.Time,
	date time.Time,
) Schedule {
	schedulePeriods := CalculateSchedulePeriods(freePeriods, busyPeriods, workBreaks)
	next := productionCalendar.NowOrNextWorkingDay(date.AddDate(0, 0, 1))
	prev := productionCalendar.NowOrPrevWorkingDay(date.AddDate(0, 0, -1))
	return Schedule{
		Date:     date,
		Periods:  schedulePeriods,
		NextDate: next,
		PrevDate: prev,
	}
}

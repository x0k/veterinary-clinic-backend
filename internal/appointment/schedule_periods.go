package appointment

import "github.com/x0k/veterinary-clinic-backend/internal/entity"

type SchedulePeriods TitledTimePeriods

func NewSchedulePeriods(
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

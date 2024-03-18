package entity

type FreeTimeSlot TimePeriod

type FreeTimeSlots []TimePeriod

func CalculateFreeTimeSlots(
	freePeriods FreePeriods,
	busyPeriods BusyPeriods,
	workBreaks CalculatedWorkBreaks,
) FreeTimeSlots {
	allBusyPeriods := make([]TimePeriod, len(busyPeriods), len(busyPeriods)+len(workBreaks))
	copy(allBusyPeriods, busyPeriods)
	for _, wb := range workBreaks {
		allBusyPeriods = append(allBusyPeriods, wb.Period)
	}
	return TimePeriodApi.SortAndUnitePeriods(
		TimePeriodApi.SubtractPeriodsFromPeriods(
			freePeriods,
			allBusyPeriods,
		),
	)
}

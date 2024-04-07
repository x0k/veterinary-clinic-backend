package appointment

import "github.com/x0k/veterinary-clinic-backend/internal/entity"

type FreeTimeSlots []entity.TimePeriod

func NewFreeTimeSlots(
	freePeriods FreePeriods,
	busyPeriods BusyPeriods,
	workBreaks WorkBreaks,
) FreeTimeSlots {
	allBusyPeriods := make([]entity.TimePeriod, len(busyPeriods), len(busyPeriods)+len(workBreaks))
	copy(allBusyPeriods, busyPeriods)
	for _, wb := range workBreaks {
		allBusyPeriods = append(allBusyPeriods, wb.Period)
	}
	return entity.TimePeriodApi.SortAndUnitePeriods(
		entity.TimePeriodApi.SubtractPeriodsFromPeriods(
			freePeriods,
			allBusyPeriods,
		),
	)
}

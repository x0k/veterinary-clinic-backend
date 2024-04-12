package appointment

import (
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type FreeTimeSlots []entity.TimePeriod

func NewFreeTimeSlots(
	dayTimePeriods DayTimePeriods,
	busyPeriods BusyPeriods,
	workBreaks DayWorkBreaks,
) (FreeTimeSlots, error) {
	allBusyPeriods := make([]entity.TimePeriod, len(busyPeriods), len(busyPeriods)+len(workBreaks))
	copy(allBusyPeriods, busyPeriods)
	for _, wb := range workBreaks {
		allBusyPeriods = append(allBusyPeriods, wb.Period)
	}
	return entity.TimePeriodApi.SortAndUnitePeriods(
		entity.TimePeriodApi.SubtractPeriodsFromPeriods(
			dayTimePeriods.Periods,
			allBusyPeriods,
		),
	), nil
}

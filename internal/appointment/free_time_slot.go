package appointment

import (
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type FreeTimeSlots []shared.TimePeriod

func NewFreeTimeSlots(
	dayTimePeriods DayTimePeriods,
	busyPeriods BusyPeriods,
	workBreaks DayWorkBreaks,
) (FreeTimeSlots, error) {
	allBusyPeriods := make([]shared.TimePeriod, len(busyPeriods), len(busyPeriods)+len(workBreaks))
	copy(allBusyPeriods, busyPeriods)
	for _, wb := range workBreaks {
		allBusyPeriods = append(allBusyPeriods, wb.Period)
	}
	return shared.TimePeriodApi.SortAndUnitePeriods(
		shared.TimePeriodApi.SubtractPeriodsFromPeriods(
			dayTimePeriods.Periods,
			allBusyPeriods,
		),
	), nil
}

func (slots FreeTimeSlots) Includes(period shared.TimePeriod) bool {
	for _, s := range slots {
		if shared.TimePeriodApi.Includes(s, period) {
			return true
		}
	}
	return false
}

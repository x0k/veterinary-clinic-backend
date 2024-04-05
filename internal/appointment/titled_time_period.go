package appointment

import (
	"slices"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type TimePeriodType int

const (
	FreePeriod TimePeriodType = iota
	BusyPeriod
)

type TitledTimePeriod struct {
	entity.TimePeriod
	Type  TimePeriodType
	Title string
}

type TitledTimePeriods []TitledTimePeriod

func (periods TitledTimePeriods) SortAndFlat() TitledTimePeriods {
	if len(periods) < 2 {
		return periods
	}
	flat := slices.Clone(periods)
	slices.SortFunc(flat, func(a, b TitledTimePeriod) int {
		return entity.TimePeriodApi.ComparePeriods(a.TimePeriod, b.TimePeriod)
	})
	nextIndex := 1
	for i := 1; i < len(flat); i++ {
		prevPeriod := flat[nextIndex-1]
		currentPeriod := flat[i]
		if entity.TimePeriodApi.IsValidPeriod(
			entity.TimePeriodApi.IntersectPeriods(prevPeriod.TimePeriod, currentPeriod.TimePeriod),
		) {
			if prevPeriod.Type == BusyPeriod || currentPeriod.Type == FreePeriod {
				diff := entity.TimePeriodApi.SubtractPeriods(currentPeriod.TimePeriod, prevPeriod.TimePeriod)
				if len(diff) == 0 {
					continue
				}
				flat[nextIndex] = TitledTimePeriod{
					TimePeriod: diff[0],
					Type:       currentPeriod.Type,
					Title:      currentPeriod.Title,
				}
			} else {
				diff := entity.TimePeriodApi.SubtractPeriods(prevPeriod.TimePeriod, currentPeriod.TimePeriod)
				if len(diff) == 0 {
					flat[nextIndex-1] = currentPeriod
					continue
				}
				flat[nextIndex-1] = TitledTimePeriod{
					TimePeriod: diff[0],
					Type:       prevPeriod.Type,
					Title:      prevPeriod.Title,
				}
				flat[nextIndex] = currentPeriod
			}
		} else {
			flat[nextIndex] = currentPeriod
		}
		nextIndex++
	}
	return flat[:nextIndex]
}

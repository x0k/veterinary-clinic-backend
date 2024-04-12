package appointment

import (
	"slices"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type ScheduleEntryType int

const (
	FreePeriod ScheduleEntryType = iota
	BusyPeriod
)

type ScheduleEntry struct {
	entity.TimePeriod
	Type  ScheduleEntryType
	Title string
}

type ScheduleEntries []ScheduleEntry

// Performs a mutation
//
// Returns a modified slice
func (periods ScheduleEntries) OmitPast(now time.Time) ScheduleEntries {
	shift := 0
	for i := 0; i < len(periods); i++ {
		end := periods[i].End
		if end.Hours < now.Hour() ||
			end.Hours == now.Hour() && end.Minutes < now.Minute() {
			shift++
		} else {
			periods[i-shift] = periods[i]
		}
	}
	return periods[:len(periods)-shift]
}

// Performs a mutation
//
// Returns a modified slice
func (periods ScheduleEntries) SortAndFlat() ScheduleEntries {
	if len(periods) < 2 {
		return periods
	}
	slices.SortFunc(periods, func(a, b ScheduleEntry) int {
		return entity.TimePeriodApi.ComparePeriods(a.TimePeriod, b.TimePeriod)
	})
	nextIndex := 1
	for i := 1; i < len(periods); i++ {
		prevPeriod := periods[nextIndex-1]
		currentPeriod := periods[i]
		if entity.TimePeriodApi.IsValidPeriod(
			entity.TimePeriodApi.IntersectPeriods(prevPeriod.TimePeriod, currentPeriod.TimePeriod),
		) {
			if prevPeriod.Type == BusyPeriod || currentPeriod.Type == FreePeriod {
				diff := entity.TimePeriodApi.SubtractPeriods(currentPeriod.TimePeriod, prevPeriod.TimePeriod)
				if len(diff) == 0 {
					continue
				}
				periods[nextIndex] = ScheduleEntry{
					TimePeriod: diff[0],
					Type:       currentPeriod.Type,
					Title:      currentPeriod.Title,
				}
			} else {
				diff := entity.TimePeriodApi.SubtractPeriods(prevPeriod.TimePeriod, currentPeriod.TimePeriod)
				if len(diff) == 0 {
					periods[nextIndex-1] = currentPeriod
					continue
				}
				periods[nextIndex-1] = ScheduleEntry{
					TimePeriod: diff[0],
					Type:       prevPeriod.Type,
					Title:      prevPeriod.Title,
				}
				periods[nextIndex] = currentPeriod
			}
		} else {
			periods[nextIndex] = currentPeriod
		}
		nextIndex++
	}
	return periods[:nextIndex]
}

func newScheduleEntries(
	freeTimeSlots FreeTimeSlots,
	busyPeriods BusyPeriods,
	workBreaks DayWorkBreaks,
) ScheduleEntries {
	periods := make(ScheduleEntries, 0, len(freeTimeSlots)+len(busyPeriods)+len(workBreaks))
	for _, p := range freeTimeSlots {
		periods = append(periods, ScheduleEntry{
			TimePeriod: p,
			Type:       FreePeriod,
			Title:      "Свободно",
		})
	}
	for _, p := range busyPeriods {
		periods = append(periods, ScheduleEntry{
			TimePeriod: p,
			Type:       BusyPeriod,
			Title:      "Занято",
		})
	}
	for _, p := range workBreaks {
		periods = append(periods, ScheduleEntry{
			TimePeriod: p.Period,
			Type:       BusyPeriod,
			Title:      p.Title,
		})
	}
	return periods.SortAndFlat()
}

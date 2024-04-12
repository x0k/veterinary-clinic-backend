package appointment

import (
	"slices"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type ScheduleEntryType int

const (
	FreePeriod ScheduleEntryType = iota
	BusyPeriod
)

type scheduleEntry struct {
	entity.DateTimePeriod
	Type  ScheduleEntryType
	Title string
}

type scheduleEntries []scheduleEntry

// Performs a mutation
//
// Returns a modified slice
func (periods scheduleEntries) SortAndFlat() scheduleEntries {
	if len(periods) < 2 {
		return periods
	}
	slices.SortFunc(periods, func(a, b scheduleEntry) int {
		return entity.DateTimePeriodApi.ComparePeriods(a.DateTimePeriod, b.DateTimePeriod)
	})
	nextIndex := 1
	for i := 1; i < len(periods); i++ {
		prevPeriod := periods[nextIndex-1]
		currentPeriod := periods[i]
		if entity.DateTimePeriodApi.IsValidPeriod(
			entity.DateTimePeriodApi.IntersectPeriods(prevPeriod.DateTimePeriod, currentPeriod.DateTimePeriod),
		) {
			if prevPeriod.Type == BusyPeriod || currentPeriod.Type == FreePeriod {
				diff := entity.DateTimePeriodApi.SubtractPeriods(currentPeriod.DateTimePeriod, prevPeriod.DateTimePeriod)
				if len(diff) == 0 {
					continue
				}
				periods[nextIndex] = scheduleEntry{
					DateTimePeriod: diff[0],
					Type:           currentPeriod.Type,
					Title:          currentPeriod.Title,
				}
			} else {
				diff := entity.DateTimePeriodApi.SubtractPeriods(prevPeriod.DateTimePeriod, currentPeriod.DateTimePeriod)
				if len(diff) == 0 {
					periods[nextIndex-1] = currentPeriod
					continue
				}
				periods[nextIndex-1] = scheduleEntry{
					DateTimePeriod: diff[0],
					Type:           prevPeriod.Type,
					Title:          prevPeriod.Title,
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

// Performs a mutation
//
// Returns a modified slice
func (periods scheduleEntries) OmitPast(now entity.DateTime) scheduleEntries {
	shift := 0
	for i := 0; i < len(periods); i++ {
		end := periods[i].End
		if entity.CompareDateTime(end, now) < 0 {
			shift++
		} else {
			periods[i-shift] = periods[i]
		}
	}
	return periods[:len(periods)-shift]
}

func newScheduleEntries(
	appointmentDate entity.Date,
	freeTimeSlots FreeTimeSlots,
	busyPeriods BusyPeriods,
	workBreaks DayWorkBreaks,
) scheduleEntries {
	periods := make(scheduleEntries, 0, len(freeTimeSlots)+len(busyPeriods)+len(workBreaks))
	for _, p := range freeTimeSlots {
		periods = append(periods, scheduleEntry{
			DateTimePeriod: entity.DateTimePeriod{
				Start: entity.DateTime{
					Date: appointmentDate,
					Time: p.Start,
				},
				End: entity.DateTime{
					Date: appointmentDate,
					Time: p.End,
				},
			},
			Type:  FreePeriod,
			Title: "Свободно",
		})
	}
	for _, p := range busyPeriods {
		periods = append(periods, scheduleEntry{
			DateTimePeriod: entity.DateTimePeriod{
				Start: entity.DateTime{
					Date: appointmentDate,
					Time: p.Start,
				},
				End: entity.DateTime{
					Date: appointmentDate,
					Time: p.End,
				},
			},
			Type:  BusyPeriod,
			Title: "Занято",
		})
	}
	for _, p := range workBreaks {
		periods = append(periods, scheduleEntry{
			DateTimePeriod: entity.DateTimePeriod{
				Start: entity.DateTime{
					Date: appointmentDate,
					Time: p.Period.Start,
				},
				End: entity.DateTime{
					Date: appointmentDate,
					Time: p.Period.End,
				},
			},
			Type:  BusyPeriod,
			Title: p.Title,
		})
	}
	return periods.SortAndFlat()
}

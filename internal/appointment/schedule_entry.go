package appointment

import (
	"slices"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type ScheduleEntryType int

func (t ScheduleEntryType) Int() int {
	return int(t)
}

const (
	FreePeriod ScheduleEntryType = iota
	BusyPeriod
)

type ScheduleEntry struct {
	shared.DateTimePeriod
	Type  ScheduleEntryType
	Title string
}

type ScheduleEntries []ScheduleEntry

// Performs a mutation
//
// Returns a modified slice
func (periods ScheduleEntries) SortAndFlat() ScheduleEntries {
	if len(periods) < 2 {
		return periods
	}
	slices.SortFunc(periods, func(a, b ScheduleEntry) int {
		return shared.DateTimePeriodApi.ComparePeriods(a.DateTimePeriod, b.DateTimePeriod)
	})
	nextIndex := 1
	for i := 1; i < len(periods); i++ {
		prevPeriod := periods[nextIndex-1]
		currentPeriod := periods[i]
		if shared.DateTimePeriodApi.IsValidPeriod(
			shared.DateTimePeriodApi.IntersectPeriods(prevPeriod.DateTimePeriod, currentPeriod.DateTimePeriod),
		) {
			if prevPeriod.Type == BusyPeriod || currentPeriod.Type == FreePeriod {
				diff := shared.DateTimePeriodApi.SubtractPeriods(currentPeriod.DateTimePeriod, prevPeriod.DateTimePeriod)
				if len(diff) == 0 {
					continue
				}
				periods[nextIndex] = ScheduleEntry{
					DateTimePeriod: diff[0],
					Type:           currentPeriod.Type,
					Title:          currentPeriod.Title,
				}
			} else {
				diff := shared.DateTimePeriodApi.SubtractPeriods(prevPeriod.DateTimePeriod, currentPeriod.DateTimePeriod)
				if len(diff) == 0 {
					periods[nextIndex-1] = currentPeriod
					continue
				}
				periods[nextIndex-1] = ScheduleEntry{
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
func (periods ScheduleEntries) OmitPast(now shared.DateTime) ScheduleEntries {
	shift := 0
	for i := 0; i < len(periods); i++ {
		end := periods[i].End
		if shared.CompareDateTime(end, now) < 0 {
			shift++
		} else {
			periods[i-shift] = periods[i]
		}
	}
	return periods[:len(periods)-shift]
}

func newScheduleEntries(
	appointmentDate shared.Date,
	freeTimeSlots FreeTimeSlots,
	busyPeriods BusyPeriods,
	workBreaks DayWorkBreaks,
) ScheduleEntries {
	periods := make(ScheduleEntries, 0, len(freeTimeSlots)+len(busyPeriods)+len(workBreaks))
	for _, p := range freeTimeSlots {
		periods = append(periods, ScheduleEntry{
			DateTimePeriod: shared.DateTimePeriod{
				Start: shared.DateTime{
					Date: appointmentDate,
					Time: p.Start,
				},
				End: shared.DateTime{
					Date: appointmentDate,
					Time: p.End,
				},
			},
			Type:  FreePeriod,
			Title: "Свободно",
		})
	}
	for _, p := range busyPeriods {
		periods = append(periods, ScheduleEntry{
			DateTimePeriod: shared.DateTimePeriod{
				Start: shared.DateTime{
					Date: appointmentDate,
					Time: p.Start,
				},
				End: shared.DateTime{
					Date: appointmentDate,
					Time: p.End,
				},
			},
			Type:  BusyPeriod,
			Title: "Занято",
		})
	}
	for _, p := range workBreaks {
		periods = append(periods, ScheduleEntry{
			DateTimePeriod: shared.DateTimePeriod{
				Start: shared.DateTime{
					Date: appointmentDate,
					Time: p.Period.Start,
				},
				End: shared.DateTime{
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

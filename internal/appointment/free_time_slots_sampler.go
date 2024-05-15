package appointment

import "github.com/x0k/veterinary-clinic-backend/internal/shared"

type SampledFreeTimeSlots []shared.TimePeriod

type freeTimeSlotsSampler struct {
	durationInMinutes   shared.DurationInMinutes
	durationShift       func(shared.Time) shared.Time
	sampleRateInMinutes SampleRateInMinutes
	sampleRateShift     func(shared.Time) shared.Time
}

func newFreeTimeSlotsSampler(
	durationInMinutes shared.DurationInMinutes,
	sampleRateInMinutes SampleRateInMinutes,
) *freeTimeSlotsSampler {
	durationShift := shared.MakeTimeShifter(shared.Time{
		Minutes: durationInMinutes.Int(),
	})
	sampleRateShift := shared.MakeTimeShifter(shared.Time{
		Minutes: sampleRateInMinutes.Minutes(),
	})
	return &freeTimeSlotsSampler{
		durationInMinutes:   durationInMinutes,
		durationShift:       durationShift,
		sampleRateInMinutes: sampleRateInMinutes,
		sampleRateShift:     sampleRateShift,
	}
}

func (c *freeTimeSlotsSampler) Sample(period shared.TimePeriod) SampledFreeTimeSlots {
	rest := period.Start.Minutes % c.sampleRateInMinutes.Minutes()
	if rest != 0 {
		return c.Sample(shared.TimePeriod{
			Start: shared.MakeTimeShifter(shared.Time{
				Minutes: c.sampleRateInMinutes.Minutes() - rest,
			})(period.Start),
			End: period.End,
		})
	}
	periodDuration := shared.TimePeriodDurationInMinutes(period)
	if periodDuration < c.durationInMinutes {
		return []shared.TimePeriod{period}
	}
	periods := make([]shared.TimePeriod, 1)
	periods[0] = shared.TimePeriod{
		Start: period.Start,
		End:   c.durationShift(period.Start),
	}
	for _, p := range shared.TimePeriodApi.SubtractPeriods(period, shared.TimePeriod{
		Start: period.Start,
		End:   c.sampleRateShift(period.Start),
	}) {
		periods = append(periods, c.Sample(p)...)
	}
	return periods
}

func NewSampleFreeTimeSlots(
	durationInMinutes shared.DurationInMinutes,
	sampleRateInMinutes SampleRateInMinutes,
	periods FreeTimeSlots,
) SampledFreeTimeSlots {
	calculator := newFreeTimeSlotsSampler(
		durationInMinutes,
		sampleRateInMinutes,
	)
	sampledPeriods := make(SampledFreeTimeSlots, 0, len(periods))
	for _, p := range periods {
		sampledPeriod := calculator.Sample(p)
		sampledPeriods = append(sampledPeriods, sampledPeriod...)
	}
	return sampledPeriods
}

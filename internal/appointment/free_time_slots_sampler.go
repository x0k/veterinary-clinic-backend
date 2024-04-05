package appointment

import "github.com/x0k/veterinary-clinic-backend/internal/entity"

type SampledFreeTimeSlots []entity.TimePeriod

type freeTimeSlotsSampler struct {
	durationInMinutes   entity.DurationInMinutes
	durationShift       func(entity.Time) entity.Time
	sampleRateInMinutes SampleRateInMinutes
	sampleRateShift     func(entity.Time) entity.Time
}

func newFreeTimeSlotsSampler(
	durationInMinutes entity.DurationInMinutes,
	sampleRateInMinutes SampleRateInMinutes,
) *freeTimeSlotsSampler {
	durationShift := entity.MakeTimeShifter(entity.Time{
		Minutes: int(durationInMinutes),
	})
	sampleRateShift := entity.MakeTimeShifter(entity.Time{
		Minutes: int(sampleRateInMinutes),
	})
	return &freeTimeSlotsSampler{
		durationInMinutes:   durationInMinutes,
		durationShift:       durationShift,
		sampleRateInMinutes: sampleRateInMinutes,
		sampleRateShift:     sampleRateShift,
	}
}

func (c *freeTimeSlotsSampler) Sample(period entity.TimePeriod) SampledFreeTimeSlots {
	rest := period.Start.Minutes % int(c.sampleRateInMinutes)
	if rest != 0 {
		return c.Sample(entity.TimePeriod{
			Start: entity.MakeTimeShifter(entity.Time{
				Minutes: int(c.sampleRateInMinutes) - rest,
			})(period.Start),
			End: period.End,
		})
	}
	periodDuration := entity.TimePeriodDurationInMinutes(period)
	if periodDuration < c.durationInMinutes {
		return []entity.TimePeriod{period}
	}
	periods := make([]entity.TimePeriod, 1)
	periods[0] = entity.TimePeriod{
		Start: period.Start,
		End:   c.durationShift(period.Start),
	}
	for _, p := range entity.TimePeriodApi.SubtractPeriods(period, entity.TimePeriod{
		Start: period.Start,
		End:   c.sampleRateShift(period.Start),
	}) {
		periods = append(periods, c.Sample(p)...)
	}
	return periods
}

func NewSampleFreeTimeSlots(
	durationInMinutes entity.DurationInMinutes,
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

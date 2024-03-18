package entity

type SampledFreeTimeSlots []TimePeriod

type FreeTimeSlotsSampler struct {
	durationInMinutes   DurationInMinutes
	durationShift       func(Time) Time
	sampleRateInMinutes SampleRateInMinutes
	sampleRateShift     func(Time) Time
}

func NewFreeTimeSlotsSampler(
	durationInMinutes DurationInMinutes,
	sampleRateInMinutes SampleRateInMinutes,
) *FreeTimeSlotsSampler {
	durationShift := MakeTimeShifter(Time{
		Minutes: int(durationInMinutes),
	})
	sampleRateShift := MakeTimeShifter(Time{
		Minutes: int(sampleRateInMinutes),
	})
	return &FreeTimeSlotsSampler{
		durationInMinutes:   durationInMinutes,
		durationShift:       durationShift,
		sampleRateInMinutes: sampleRateInMinutes,
		sampleRateShift:     sampleRateShift,
	}
}

func (c *FreeTimeSlotsSampler) Sample(period TimePeriod) SampledFreeTimeSlots {
	rest := period.Start.Minutes % int(c.sampleRateInMinutes)
	if rest != 0 {
		return c.Sample(TimePeriod{
			Start: MakeTimeShifter(Time{
				Minutes: int(c.sampleRateInMinutes) - rest,
			})(period.Start),
			End: period.End,
		})
	}
	periodDuration := TimePeriodDurationInMinutes(period)
	if periodDuration < c.durationInMinutes {
		return []TimePeriod{period}
	}
	periods := make([]TimePeriod, 1)
	periods[0] = TimePeriod{
		Start: period.Start,
		End:   c.durationShift(period.Start),
	}
	for _, p := range TimePeriodApi.SubtractPeriods(period, TimePeriod{
		Start: period.Start,
		End:   c.sampleRateShift(period.Start),
	}) {
		periods = append(periods, c.Sample(p)...)
	}
	return periods
}

func SampleFreeTimeSlots(
	durationInMinutes DurationInMinutes,
	sampleRateInMinutes SampleRateInMinutes,
	periods FreeTimeSlots,
) SampledFreeTimeSlots {
	calculator := NewFreeTimeSlotsSampler(
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

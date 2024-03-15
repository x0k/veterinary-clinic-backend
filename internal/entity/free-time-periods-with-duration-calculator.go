package entity

type FreeTimePeriodsWithDurationCalculator struct {
	durationInMinutes   int
	durationShift       func(Time) Time
	sampleRateInMinutes int
	sampleRateShift     func(Time) Time
}

func NewFreeTimePeriodsWithDurationCalculator(
	durationInMinutes int,
	sampleRateInMinutes int,
) *FreeTimePeriodsWithDurationCalculator {
	durationShift := MakeTimeShifter(Time{
		Minutes: durationInMinutes,
	})
	sampleRateShift := MakeTimeShifter(Time{
		Minutes: sampleRateInMinutes,
	})
	return &FreeTimePeriodsWithDurationCalculator{
		durationInMinutes:   durationInMinutes,
		durationShift:       durationShift,
		sampleRateInMinutes: sampleRateInMinutes,
		sampleRateShift:     sampleRateShift,
	}
}

func (c *FreeTimePeriodsWithDurationCalculator) Calculate(period TimePeriod) []TimePeriod {
	rest := period.Start.Minutes % c.sampleRateInMinutes
	if rest != 0 {
		return c.Calculate(TimePeriod{
			Start: MakeTimeShifter(Time{
				Minutes: c.sampleRateInMinutes - rest,
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
		periods = append(periods, c.Calculate(p)...)
	}
	return periods
}

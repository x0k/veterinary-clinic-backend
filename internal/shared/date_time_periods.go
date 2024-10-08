package shared

import "github.com/x0k/veterinary-clinic-backend/internal/lib/period"

type TimePeriod = period.Period[Time]

type DateTimePeriod = period.Period[DateTime]

var TimePeriodApi = period.NewApi(CompareTime)

var DateTimePeriodApi = period.NewApi(CompareDateTime)

func TimePeriodDurationInMinutes(period TimePeriod) DurationInMinutes {
	return DurationInMinutes(
		(period.End.Hours-period.Start.Hours)*60 + (period.End.Minutes - period.Start.Minutes),
	)
}

func DateToDayTimePeriod(d Date) DateTimePeriod {
	return DateTimePeriod{
		Start: DateTime{
			Date: d,
			Time: Time{
				Hours:   0,
				Minutes: 0,
			},
		},
		End: DateTime{
			Date: d,
			Time: Time{
				Hours:   23,
				Minutes: 59,
			},
		},
	}
}

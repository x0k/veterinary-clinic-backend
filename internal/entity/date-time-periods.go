package entity

import "github.com/x0k/veterinary-clinic-backend/internal/lib/period"

type TimePeriod = period.Period[Time]

type DateTimePeriod = period.Period[DateTime]

var TimePeriodApi = period.NewApi(CompareTime)

var DateTimePeriodApi = period.NewApi(CompareDateTime)

func TimePeriodDurationInMinutes(period TimePeriod) int {
	return (period.End.Hours-period.Start.Hours)*60 + (period.End.Minutes - period.Start.Minutes)
}

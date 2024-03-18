package entity

type DayTimePeriods struct {
	Date    Date
	Periods []TimePeriod
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

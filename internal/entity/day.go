package entity

type DayType int

const (
	Weekend    DayType = 1
	Holiday    DayType = 2
	PreHoliday DayType = 3
)

func IsNonWorkingDayType(dayType DayType) bool {
	return dayType == Holiday || dayType == Weekend
}

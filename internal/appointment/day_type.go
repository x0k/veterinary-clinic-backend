package appointment

import "errors"

var ErrInvalidDayType = errors.New("invalid day type")

type DayType int

const (
	Weekend    DayType = 1
	Holiday    DayType = 2
	PreHoliday DayType = 3
)

func NewDayType(dayType int) (DayType, error) {
	if dayType < 1 || dayType > 3 {
		return 0, ErrInvalidDayType
	}
	return DayType(dayType), nil
}

func (d DayType) Int() int {
	return int(d)
}

func IsNonWorkingDayType(dayType DayType) bool {
	return dayType == Holiday || dayType == Weekend
}

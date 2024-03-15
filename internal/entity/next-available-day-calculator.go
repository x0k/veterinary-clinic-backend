package entity

import "time"

type NextAvailableDayCalculator struct {
	productionCalendar ProductionCalendar
}

func NewNextAvailableDayCalculator(
	productionCalendar ProductionCalendar,
) *NextAvailableDayCalculator {
	return &NextAvailableDayCalculator{
		productionCalendar: productionCalendar,
	}
}

func (c *NextAvailableDayCalculator) Calculate(
	today time.Time,
) time.Time {
	nextDay := today
	for {
		nextDay = nextDay.AddDate(0, 0, 1)
		nextDayJson := GoTimeToJsonDate(nextDay)
		if dayType, ok := c.productionCalendar[nextDayJson]; !ok || !IsNonWorkingDayType(dayType) {
			return nextDay
		}
	}
}

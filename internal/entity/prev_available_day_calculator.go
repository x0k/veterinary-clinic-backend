package entity

import "time"

type PrevAvailableDayCalculator struct {
	productionCalendar ProductionCalendar
	nowDate            time.Time
}

func NewPrevAvailableDayCalculator(
	productionCalendar ProductionCalendar,
	now time.Time,
) *PrevAvailableDayCalculator {
	return &PrevAvailableDayCalculator{
		productionCalendar: productionCalendar,
		nowDate:            time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
	}
}

func (c *PrevAvailableDayCalculator) Calculate(
	fromDate time.Time,
) *time.Time {
	prevDay := fromDate
	for prevDay.Sub(c.nowDate) > 24*time.Hour {
		prevDayJson := GoTimeToJsonDate(prevDay)
		if dayType, ok := c.productionCalendar[prevDayJson]; !ok || !IsNonWorkingDayType(dayType) {
			return &prevDay
		}
		prevDay = prevDay.AddDate(0, 0, -1)
	}
	return nil
}

package usecase

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

func freePeriods(
	productionCalendar entity.ProductionCalendar,
	openingHours entity.OpeningHours,
	now time.Time,
	forDate time.Time,
) ([]entity.TimePeriod, error) {
	return entity.NewFreePeriodsCalculator(
		openingHours,
		productionCalendar,
		entity.GoTimeToDateTime(now),
	).Calculate(forDate)
}

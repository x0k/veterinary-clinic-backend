package usecase

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

func nextAvailableDay(productionCalendar entity.ProductionCalendar, from time.Time) time.Time {
	return entity.NewNextAvailableDayCalculator(productionCalendar).Calculate(from)
}

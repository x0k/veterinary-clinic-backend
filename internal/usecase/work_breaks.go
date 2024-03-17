package usecase

import (
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

func workBreaks(workBreaks entity.WorkBreaks, date time.Time) (entity.WorkBreaks, error) {
	return entity.NewWorkBreaksCalculator(workBreaks).Calculate(date)
}

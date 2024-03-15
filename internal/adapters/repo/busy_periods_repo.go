package repo

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type BusyPeriods struct {
}

func NewBusyPeriods() *BusyPeriods {
	return &BusyPeriods{}
}

func (s *BusyPeriods) BusyPeriods(ctx context.Context, t time.Time) ([]entity.TimePeriod, error) {
	records, err := s.recordsRepo.Records(ctx)
	if err != nil {
		return nil, err
	}
	allBusyPeriods := make([]entity.DateTimePeriod, 0, len(records))
	for _, record := range records {
		allBusyPeriods = append(allBusyPeriods, record.DateTimePeriod)
	}
	busyPeriodsCalculator := entity.NewBusyPeriodsCalculator(allBusyPeriods)
	actualBusyPeriods := busyPeriodsCalculator.Calculate(t)
	timePeriods := make([]entity.TimePeriod, 0, len(actualBusyPeriods))
	for _, busyPeriod := range actualBusyPeriods {
		timePeriods = append(timePeriods, entity.TimePeriod{
			Start: busyPeriod.Start.Time,
			End:   busyPeriod.End.Time,
		})
	}
	return timePeriods, nil
}

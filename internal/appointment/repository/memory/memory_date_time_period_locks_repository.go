package appointment_in_memory_repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/slicex"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type DateTimePeriodLocksRepository struct {
	mu      sync.Mutex
	periods []shared.DateTimePeriod
}

func NewDateTimePeriodLocksRepository() *DateTimePeriodLocksRepository {
	return &DateTimePeriodLocksRepository{}
}

func (r *DateTimePeriodLocksRepository) Lock(ctx context.Context, period shared.DateTimePeriod) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.periods {
		if shared.DateTimePeriodApi.IsValidPeriod(
			shared.DateTimePeriodApi.IntersectPeriods(p, period),
		) {
			return fmt.Errorf("%w: %s", appointment.ErrPeriodIsLocked, period)
		}
	}
	r.periods = append(r.periods, period)
	return nil
}

func (r *DateTimePeriodLocksRepository) UnLock(ctx context.Context, period shared.DateTimePeriod) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	index := slicex.Index(r.periods, period)
	if index == -1 {
		return nil
	}
	r.periods = slicex.Delete(r.periods, index, index+1)
	return nil
}

//go:build js && wasm

package appointment_js_repository

import (
	"context"
	"syscall/js"

	"github.com/x0k/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	shared_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/shared/adapters/js"
)

type DateTimePeriodLocksRepositoryConfig struct {
	Lock   *js.Value `js:"lock"`
	UnLock *js.Value `js:"unLock"`
}

type DateTimePeriodLocksRepository struct {
	cfg DateTimePeriodLocksRepositoryConfig
}

func NewDateTimePeriodLocksRepository(cfg DateTimePeriodLocksRepositoryConfig) *DateTimePeriodLocksRepository {
	return &DateTimePeriodLocksRepository{cfg: cfg}
}

func (r *DateTimePeriodLocksRepository) Lock(ctx context.Context, period shared.DateTimePeriod) error {
	promise := r.cfg.Lock.Invoke(
		vert.ValueOf(shared_js_adapters.DateTimePeriodToDTO(period)),
	)
	_, err := js_adapters.Await(ctx, promise)
	return err
}

func (r *DateTimePeriodLocksRepository) UnLock(ctx context.Context, period shared.DateTimePeriod) error {
	promise := r.cfg.UnLock.Invoke(
		vert.ValueOf(shared_js_adapters.DateTimePeriodToDTO(period)),
	)
	_, err := js_adapters.Await(ctx, promise)
	return err
}

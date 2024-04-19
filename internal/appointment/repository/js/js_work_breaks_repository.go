//go:build js && wasm

package appointment_js_repository

import (
	"context"
	"syscall/js"

	"github.com/norunners/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

type WorkBreaksRepositoryConfig struct {
	WorkBreaks js.Value `js:"loadWorkBreaks"`
}

type WorkBreaksRepository struct {
	cfg WorkBreaksRepositoryConfig
}

func NewWorkBreaksRepository(
	cfg WorkBreaksRepositoryConfig,
) *WorkBreaksRepository {
	return &WorkBreaksRepository{
		cfg: cfg,
	}
}

func (r *WorkBreaksRepository) WorkBreaks(
	ctx context.Context,
) (appointment.WorkBreaks, error) {
	promise := r.cfg.WorkBreaks.Invoke()
	workBreaksJsValue, err := js_adapters.Await(ctx, promise)
	if err != nil {
		return nil, err
	}
	workBreaksDto := make([]WorkBreakDto, 0)
	if err := vert.ValueOf(workBreaksJsValue).AssignTo(&workBreaksDto); err != nil {
		return nil, err
	}
	workBreaks := make([]appointment.WorkBreak, len(workBreaksDto))
	for i, workBreakDto := range workBreaksDto {
		workBreaks[i] = WorkBreakFromDto(workBreakDto)
	}
	return workBreaks, nil
}

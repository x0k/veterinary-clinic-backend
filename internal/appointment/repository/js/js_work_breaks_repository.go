//go:build js && wasm

package appointment_js_repository

import (
	"context"
	"syscall/js"

	"github.com/norunners/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_js_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/js"
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
	workBreaksDTO := make([]appointment_js_adapters.WorkBreakDTO, 0)
	if err := vert.ValueOf(workBreaksJsValue).AssignTo(&workBreaksDTO); err != nil {
		return nil, err
	}
	workBreaks := make([]appointment.WorkBreak, len(workBreaksDTO))
	for i, workBreakDTO := range workBreaksDTO {
		workBreaks[i] = appointment_js_adapters.WorkBreakFromDTO(workBreakDTO)
	}
	return workBreaks, nil
}

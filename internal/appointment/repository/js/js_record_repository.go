//go:build js && wasm

package appointment_js_repository

import (
	"context"
	"syscall/js"

	"github.com/norunners/vert"
	js_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/js"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

type RecordRepositoryConfig struct {
	CreateRecord js.Value `js:"createRecord"`
}

type RecordRepository struct {
	cfg RecordRepositoryConfig
}

func NewRecordRepository(
	cfg RecordRepositoryConfig,
) *RecordRepository {
	return &RecordRepository{
		cfg: cfg,
	}
}

func (r *RecordRepository) CreateRecord(
	ctx context.Context,
	rec *appointment.RecordEntity,
) error {
	promise := r.cfg.CreateRecord.Invoke(
		vert.ValueOf(RecordToDto(*rec)),
	)
	recordId, err := js_adapters.Await(ctx, promise)
	if err != nil {
		return err
	}
	rec.SetId(appointment.RecordId(recordId.String()))
	return nil
}

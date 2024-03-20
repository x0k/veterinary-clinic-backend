package usecase

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type cancelAppointmentPresenter[R any] interface {
	RenderCancel() (R, error)
}

type CancelAppointmentUseCase[R any] struct {
	recordsRepo RecordsRemover
	presenter   cancelAppointmentPresenter[R]
}

func NewCancelAppointmentUseCase[R any](
	recordsRepo RecordsRemover,
	presenter cancelAppointmentPresenter[R],
) *CancelAppointmentUseCase[R] {
	return &CancelAppointmentUseCase[R]{
		recordsRepo: recordsRepo,
		presenter:   presenter,
	}
}

func (u *CancelAppointmentUseCase[R]) Cancel(ctx context.Context, recordId entity.RecordId) (R, error) {
	err := u.recordsRepo.Remove(ctx, recordId)
	if err != nil {
		return *new(R), err
	}
	return u.presenter.RenderCancel()
}

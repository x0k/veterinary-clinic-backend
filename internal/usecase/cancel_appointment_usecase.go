package usecase

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type cancelAppointmentPresenter[R any] interface {
	RenderError() (R, error)
	RenderCancel() (R, error)
}

type cancelAppointmentRecordsRepo interface {
	RecordByUserLoader
	RecordsRemover
}

type CancelAppointmentUseCase[R any] struct {
	recordsRepo cancelAppointmentRecordsRepo
	presenter   cancelAppointmentPresenter[R]
}

func NewCancelAppointmentUseCase[R any](
	recordsRepo cancelAppointmentRecordsRepo,
	presenter cancelAppointmentPresenter[R],
) *CancelAppointmentUseCase[R] {
	return &CancelAppointmentUseCase[R]{
		recordsRepo: recordsRepo,
		presenter:   presenter,
	}
}

func (u *CancelAppointmentUseCase[R]) Cancel(ctx context.Context, userId shared.UserId) (R, error) {
	rec, err := u.recordsRepo.RecordByUserId(ctx, userId)
	if err != nil || rec.Status != shared.RecordAwaits {
		return u.presenter.RenderError()
	}
	if err = u.recordsRepo.Remove(ctx, rec.Id); err != nil {
		return *new(R), err
	}
	return u.presenter.RenderCancel()
}

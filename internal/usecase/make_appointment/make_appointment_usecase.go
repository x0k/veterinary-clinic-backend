package make_appointment

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type makeAppointmentPresenter[R any] interface {
	RenderAppointmentInfo(appointment entity.Record) (R, error)
}

type MakeAppointmentUseCase[R any] struct {
	recordsRepo usecase.RecordsCreator
	presenter   makeAppointmentPresenter[R]
}

func NewMakeAppointmentUseCase[R any](
	recordsRepo usecase.RecordsCreator,
	presenter makeAppointmentPresenter[R],
) *MakeAppointmentUseCase[R] {
	return &MakeAppointmentUseCase[R]{
		recordsRepo: recordsRepo,
	}
}

func (u *MakeAppointmentUseCase[R]) MakeAppointment(
	ctx context.Context,
	user entity.User,
	serviceId entity.ServiceId,
	appointmentDateTime time.Time,
) (R, error) {
	record, err := u.recordsRepo.Create(
		ctx,
		user,
		serviceId,
		appointmentDateTime,
	)
	if err != nil {
		return *new(R), err
	}
	return u.presenter.RenderAppointmentInfo(record)
}

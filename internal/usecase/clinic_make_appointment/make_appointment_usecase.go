package clinic_make_appointment

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type clinicMakeAppointmentPresenter[R any] interface {
	RenderAppointmentInfo(appointment entity.Record) (R, error)
}

type MakeAppointmentUseCase[R any] struct {
	recordsRepo usecase.ClinicRecordsCreator
	presenter   clinicMakeAppointmentPresenter[R]
}

func NewMakeAppointmentUseCase[R any](
	recordsRepo usecase.ClinicRecordsCreator,
	presenter clinicMakeAppointmentPresenter[R],
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

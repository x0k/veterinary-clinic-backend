package clinic_make_appointment

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type appointmentConfirmationPresenter[R any] interface {
	RenderAppointmentConfirmation(service entity.Service, appointmentDateTime time.Time) (R, error)
}

type AppointmentConfirmationUseCase[R any] struct {
	clinicServicesRepo usecase.ClinicServiceLoader
	presenter          appointmentConfirmationPresenter[R]
}

func NewAppointmentConfirmationUseCase[R any](
	clinicServicesRepo usecase.ClinicServiceLoader,
	presenter appointmentConfirmationPresenter[R],
) *AppointmentConfirmationUseCase[R] {
	return &AppointmentConfirmationUseCase[R]{
		clinicServicesRepo: clinicServicesRepo,
		presenter:          presenter,
	}
}

func (u *AppointmentConfirmationUseCase[R]) AppointmentConfirmation(
	ctx context.Context,
	serviceId entity.ServiceId,
	appointmentDateTime time.Time,
) (R, error) {
	service, err := u.clinicServicesRepo.Load(ctx, serviceId)
	if err != nil {
		return *new(R), err
	}
	return u.presenter.RenderAppointmentConfirmation(service, appointmentDateTime)
}

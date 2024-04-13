package make_appointment

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"github.com/x0k/veterinary-clinic-backend/internal/usecase"
)

type appointmentConfirmationPresenter[R any] interface {
	RenderConfirmation(service shared.Service, appointmentDateTime time.Time) (R, error)
}

type AppointmentConfirmationUseCase[R any] struct {
	servicesRepo usecase.ServiceLoader
	presenter    appointmentConfirmationPresenter[R]
}

func NewAppointmentConfirmationUseCase[R any](
	servicesRepo usecase.ServiceLoader,
	presenter appointmentConfirmationPresenter[R],
) *AppointmentConfirmationUseCase[R] {
	return &AppointmentConfirmationUseCase[R]{
		servicesRepo: servicesRepo,
		presenter:    presenter,
	}
}

func (u *AppointmentConfirmationUseCase[R]) Confirmation(
	ctx context.Context,
	serviceId shared.ServiceId,
	appointmentDateTime time.Time,
) (R, error) {
	service, err := u.servicesRepo.Service(ctx, serviceId)
	if err != nil {
		return *new(R), err
	}
	return u.presenter.RenderConfirmation(service, appointmentDateTime)
}

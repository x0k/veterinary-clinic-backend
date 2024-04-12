package appointment_telegram_use_case

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

type AppointmentConfirmationUseCase[R any] struct {
	serviceLoader         appointment.ServiceLoader
	confirmationPresenter appointment.AppointmentConfirmationPresenter[R]
	errorPresenter        appointment.ErrorPresenter[R]
}

func NewAppointmentConfirmationUseCase[R any](
	serviceLoader appointment.ServiceLoader,
	confirmationPresenter appointment.AppointmentConfirmationPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *AppointmentConfirmationUseCase[R] {
	return &AppointmentConfirmationUseCase[R]{
		serviceLoader:         serviceLoader,
		confirmationPresenter: confirmationPresenter,
		errorPresenter:        errorPresenter,
	}
}

func (u *AppointmentConfirmationUseCase[R]) Confirmation(
	ctx context.Context,
	serviceId appointment.ServiceId,
	appointmentDateTime time.Time,
) (R, error) {
	service, err := u.serviceLoader.Service(ctx, serviceId)
	if err != nil {
		return u.errorPresenter.RenderError(err)
	}
	return u.confirmationPresenter.RenderConfirmation(service, appointmentDateTime)
}

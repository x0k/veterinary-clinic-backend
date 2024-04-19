package appointment_telegram_use_case

import (
	"context"
	"time"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const appointmentConfirmationUseCaseName = "appointment_telegram_use_case.AppointmentConfirmationUseCase"

type AppointmentConfirmationUseCase[R any] struct {
	log                   *logger.Logger
	serviceLoader         appointment.ServiceLoader
	confirmationPresenter appointment.AppointmentConfirmationPresenter[R]
	errorPresenter        appointment.ErrorPresenter[R]
}

func NewAppointmentConfirmationUseCase[R any](
	log *logger.Logger,
	serviceLoader appointment.ServiceLoader,
	confirmationPresenter appointment.AppointmentConfirmationPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *AppointmentConfirmationUseCase[R] {
	return &AppointmentConfirmationUseCase[R]{
		log:                   log.With(sl.Component(appointmentConfirmationUseCaseName)),
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
		u.log.Error(ctx, "failed to load service", sl.Err(err))
		return u.errorPresenter(err)
	}
	return u.confirmationPresenter.RenderConfirmation(service, appointmentDateTime)
}

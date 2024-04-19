package appointment_telegram_use_case

import (
	"context"
	"errors"
	"log/slog"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

const startMakeAppointmentDialogUseCaseName = "appointment_telegram_use_case.StartMakeAppointmentDialogUseCase"

type StartMakeAppointmentDialogUseCase[R any] struct {
	log                             *logger.Logger
	customerLoader                  appointment.CustomerByIdentityLoader
	customerActiveAppointmentLoader appointment.CustomerActiveAppointmentLoader
	servicesLoader                  appointment.ServicesLoader
	serviceLoader                   appointment.ServiceLoader
	appointmentInfoPresenter        appointment.AppointmentInfoPresenter[R]
	servicesPickerPresenter         appointment.ServicesPickerPresenter[R]
	registrationPresenter           appointment.RegistrationPresenter[R]
	errorPresenter                  appointment.ErrorPresenter[R]
}

func NewStartMakeAppointmentDialogUseCase[R any](
	log *logger.Logger,
	customerLoader appointment.CustomerByIdentityLoader,
	customerActiveAppointmentLoader appointment.CustomerActiveAppointmentLoader,
	servicesLoader appointment.ServicesLoader,
	appointmentInfoPresenter appointment.AppointmentInfoPresenter[R],
	servicesPickerPresenter appointment.ServicesPickerPresenter[R],
	registrationPresenter appointment.RegistrationPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *StartMakeAppointmentDialogUseCase[R] {
	return &StartMakeAppointmentDialogUseCase[R]{
		log:                             log.With(sl.Component(startMakeAppointmentDialogUseCaseName)),
		customerLoader:                  customerLoader,
		customerActiveAppointmentLoader: customerActiveAppointmentLoader,
		servicesLoader:                  servicesLoader,
		appointmentInfoPresenter:        appointmentInfoPresenter,
		servicesPickerPresenter:         servicesPickerPresenter,
		registrationPresenter:           registrationPresenter,
		errorPresenter:                  errorPresenter,
	}
}

func (u *StartMakeAppointmentDialogUseCase[R]) StartMakeAppointmentDialog(
	ctx context.Context,
	userId shared.TelegramUserId,
) (R, error) {
	customerIdentity := appointment.NewTelegramCustomerIdentity(userId)
	customer, err := u.customerLoader(ctx, customerIdentity)
	if errors.Is(err, shared.ErrNotFound) {
		return u.registrationPresenter.RenderRegistration(userId)
	}
	if err != nil {
		u.log.Error(ctx, "failed to find customer", slog.Int64("telegram_user_id", userId.Int()), sl.Err(err))
		return u.errorPresenter(err)
	}
	existedAppointment, err := u.customerActiveAppointmentLoader(ctx, customer.Id)
	if !errors.Is(err, shared.ErrNotFound) {
		if err != nil {
			u.log.Error(ctx, "failed to find customer active appointment", sl.Err(err))
			return u.errorPresenter(err)
		}
		service, err := u.serviceLoader.Service(ctx, existedAppointment.ServiceId)
		if err != nil {
			u.log.Error(ctx, "failed to load service", sl.Err(err))
			return u.errorPresenter(err)
		}
		return u.appointmentInfoPresenter.RenderInfo(existedAppointment, service)
	}
	services, err := u.servicesLoader.Services(ctx)
	if err != nil {
		u.log.Error(ctx, "failed to load services", sl.Err(err))
		return u.errorPresenter(err)
	}
	return u.servicesPickerPresenter.RenderServicesList(services)
}

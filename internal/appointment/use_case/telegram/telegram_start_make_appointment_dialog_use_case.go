package appointment_telegram_use_case

import (
	"context"
	"errors"
	"log/slog"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const startMakeAppointmentDialogUseCaseName = "appointment_telegram_use_case.StartMakeAppointmentDialogUseCase"

type StartMakeAppointmentDialogUseCase[R any] struct {
	log                     *logger.Logger
	customerLoader          appointment.CustomerLoader
	servicesLoader          appointment.ServicesLoader
	servicesPickerPresenter appointment.ServicesPickerPresenter[R]
	registrationPresenter   appointment.RegistrationPresenter[R]
	errorPresenter          appointment.ErrorPresenter[R]
}

func NewStartMakeAppointmentDialogUseCase[R any](
	log *logger.Logger,
	customerLoader appointment.CustomerLoader,
	servicesLoader appointment.ServicesLoader,
	servicesPickerPresenter appointment.ServicesPickerPresenter[R],
	registrationPresenter appointment.RegistrationPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *StartMakeAppointmentDialogUseCase[R] {
	return &StartMakeAppointmentDialogUseCase[R]{
		log:                     log.With(slog.String("component", startMakeAppointmentDialogUseCaseName)),
		customerLoader:          customerLoader,
		servicesLoader:          servicesLoader,
		servicesPickerPresenter: servicesPickerPresenter,
		registrationPresenter:   registrationPresenter,
		errorPresenter:          errorPresenter,
	}
}

func (u *StartMakeAppointmentDialogUseCase[R]) StartMakeAppointmentDialog(
	ctx context.Context,
	userId entity.TelegramUserId,
) (R, error) {
	customerId := appointment.TelegramUserIdToCustomerId(userId)
	_, err := u.customerLoader.Customer(ctx, customerId)
	if errors.Is(err, entity.ErrNotFound) {
		return u.registrationPresenter.RenderRegistration()
	}
	if err != nil {
		u.log.Error(ctx, "failed to find customer", slog.Int64("telegram_user_id", userId.Int()), sl.Err(err))
		return u.errorPresenter.RenderError(err)
	}
	services, err := u.servicesLoader.Services(ctx)
	if err != nil {
		u.log.Error(ctx, "failed to load services", sl.Err(err))
		return u.errorPresenter.RenderError(err)
	}
	return u.servicesPickerPresenter.RenderServicesList(services)
}

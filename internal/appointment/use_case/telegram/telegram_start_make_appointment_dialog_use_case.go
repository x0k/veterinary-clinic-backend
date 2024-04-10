package appointment_telegram_use_case

import (
	"context"
	"errors"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
)

type StartMakeAppointmentDialogUseCase[R any] struct {
	customerLoader          appointment.CustomerLoader
	servicesLoader          appointment.ServicesLoader
	servicesPickerPresenter appointment.ServicesPickerPresenter[R]
	registrationPresenter   appointment.RegistrationPresenter[R]
	errorPresenter          appointment.ErrorPresenter[R]
}

func NewStartMakeAppointmentDialogUseCase[R any]() *StartMakeAppointmentDialogUseCase[R] {
	return &StartMakeAppointmentDialogUseCase[R]{}
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
		return u.errorPresenter.RenderError(err)
	}
	services, err := u.servicesLoader.Services(ctx)
	if err != nil {
		return u.errorPresenter.RenderError(err)
	}
	return u.servicesPickerPresenter.RenderServicesList(services)
}

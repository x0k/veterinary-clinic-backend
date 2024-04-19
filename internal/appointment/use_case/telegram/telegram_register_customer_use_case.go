package appointment_telegram_use_case

import (
	"context"
	"fmt"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

const registerCustomerUseCaseName = "appointment_telegram_use_case.RegisterCustomerUseCase"

type RegisterCustomerUseCase[R any] struct {
	log                          *logger.Logger
	customerCreator              appointment.CustomerCreator
	servicesLoader               appointment.ServicesLoader
	successRegistrationPresenter appointment.SuccessRegistrationPresenter[R]
	errorPresenter               appointment.ErrorPresenter[R]
}

func NewRegisterCustomerUseCase[R any](
	log *logger.Logger,
	customerCreator appointment.CustomerCreator,
	servicesLoader appointment.ServicesLoader,
	successRegistrationPresenter appointment.SuccessRegistrationPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *RegisterCustomerUseCase[R] {
	return &RegisterCustomerUseCase[R]{
		log:                          log.With(sl.Component(registerCustomerUseCaseName)),
		customerCreator:              customerCreator,
		servicesLoader:               servicesLoader,
		successRegistrationPresenter: successRegistrationPresenter,
		errorPresenter:               errorPresenter,
	}
}

func (u *RegisterCustomerUseCase[R]) RegisterCustomer(
	ctx context.Context,
	telegramUserId shared.TelegramUserId,
	telegramUserName string,
	telegramUserFirstName string,
	telegramUserLastName string,
	telegramUserPhoneNumber string,
) (R, error) {
	customerIdentity := appointment.NewTelegramCustomerIdentity(telegramUserId)
	customer := appointment.NewCustomer(
		appointment.TemporalCustomerId,
		customerIdentity,
		fmt.Sprintf("%s %s", telegramUserFirstName, telegramUserLastName),
		telegramUserPhoneNumber,
		fmt.Sprintf("https://t.me/%s", telegramUserName),
	)
	if err := u.customerCreator.CreateCustomer(ctx, &customer); err != nil {
		u.log.Error(ctx, "failed to create customer", sl.Err(err))
		return u.errorPresenter(err)
	}
	if customer.Id == appointment.TemporalCustomerId {
		err := appointment.ErrInvalidCustomerId
		u.log.Error(ctx, "failed to create customer", sl.Err(err))
		return u.errorPresenter(err)
	}
	services, err := u.servicesLoader.Services(ctx)
	if err != nil {
		u.log.Error(ctx, "failed to load services", sl.Err(err))
		return u.errorPresenter(err)
	}
	return u.successRegistrationPresenter.RenderSuccessRegistration(services)
}

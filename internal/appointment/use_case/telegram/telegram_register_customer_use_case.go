package appointment_telegram_use_case

import (
	"context"
	"fmt"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const registerCustomerUseCaseName = "appointment_telegram_use_case.RegisterCustomerUseCase"

type RegisterCustomerUseCase[R any] struct {
	log                          *logger.Logger
	customerCreator              appointment.CustomerCreator
	successRegistrationPresenter appointment.SuccessRegistrationPresenter[R]
	errorPresenter               appointment.ErrorPresenter[R]
}

func NewRegisterCustomerUseCase[R any](
	log *logger.Logger,
	customerCreator appointment.CustomerCreator,
	successRegistrationPresenter appointment.SuccessRegistrationPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *RegisterCustomerUseCase[R] {
	return &RegisterCustomerUseCase[R]{
		log:                          log.With(sl.Component(registerCustomerUseCaseName)),
		customerCreator:              customerCreator,
		successRegistrationPresenter: successRegistrationPresenter,
		errorPresenter:               errorPresenter,
	}
}

func (u *RegisterCustomerUseCase[R]) RegisterCustomer(
	ctx context.Context,
	telegramUserId entity.TelegramUserId,
	telegramUserName string,
	telegramUserFirstName string,
	telegramUserLastName string,
	telegramUserPhoneNumber string,
) (bool, R, error) {
	customerId := appointment.TelegramUserIdToCustomerId(telegramUserId)
	customer := appointment.NewCustomer(
		customerId,
		fmt.Sprintf("%s %s", telegramUserFirstName, telegramUserLastName),
		telegramUserPhoneNumber,
		fmt.Sprintf("https://t.me/%s", telegramUserName),
	)
	if err := u.customerCreator.CreateCustomer(ctx, customer); err != nil {
		u.log.Error(ctx, "failed to create customer", sl.Err(err))
		res, err := u.errorPresenter.RenderError(err)
		return false, res, err
	}
	res, err := u.successRegistrationPresenter.RenderSuccessRegistration()
	return true, res, err
}

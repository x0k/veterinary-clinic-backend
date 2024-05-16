package appointment_js_use_case

import (
	"context"
	"errors"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

const upsertCustomerUseCaseName = "appointment_js_use_case.UpsertCustomerUseCase"

type UpsertCustomerUseCase[R any] struct {
	log                      *logger.Logger
	customerByIdentityLoader appointment.CustomerByIdentityLoader
	customerCreator          appointment.CustomerCreator
	customerUpdater          appointment.CustomerUpdater
	customerPresenter        appointment.CustomerPresenter[R]
	errorPresenter           appointment.ErrorPresenter[R]
}

func NewUpsertCustomerUseCase[R any](
	log *logger.Logger,
	customerByIdentityLoader appointment.CustomerByIdentityLoader,
	customerCreator appointment.CustomerCreator,
	customerUpdater appointment.CustomerUpdater,
	customerPresenter appointment.CustomerPresenter[R],
	errorPresenter appointment.ErrorPresenter[R],
) *UpsertCustomerUseCase[R] {
	return &UpsertCustomerUseCase[R]{
		log:                      log.With(sl.Component(upsertCustomerUseCaseName)),
		customerByIdentityLoader: customerByIdentityLoader,
		customerCreator:          customerCreator,
		customerUpdater:          customerUpdater,
		customerPresenter:        customerPresenter,
		errorPresenter:           errorPresenter,
	}
}

func (u *UpsertCustomerUseCase[R]) Upsert(
	ctx context.Context,
	customerIdentity appointment.CustomerIdentity,
	userName string,
	userPhone string,
	userEmail string,
) (R, error) {
	customer, err := u.customerByIdentityLoader(ctx, customerIdentity)
	if errors.Is(err, shared.ErrNotFound) {
		return u.createCustomer(ctx, customerIdentity, userName, userPhone, userEmail)
	}
	if err != nil {
		u.log.Debug(ctx, "failed to load customer", sl.Err(err))
		return u.errorPresenter(err)
	}
	if !customer.Update(userName, userPhone, userEmail) {
		return u.customerPresenter(customer)
	}
	if err := u.customerUpdater(ctx, customer); err != nil {
		u.log.Debug(ctx, "failed to update customer", sl.Err(err))
		return u.errorPresenter(err)
	}
	return u.customerPresenter(customer)
}

func (u *UpsertCustomerUseCase[R]) createCustomer(
	ctx context.Context,
	customerIdentity appointment.CustomerIdentity,
	userName string,
	userPhone string,
	userEmail string,
) (R, error) {
	customer := appointment.NewCustomer(
		appointment.TemporalCustomerId,
		customerIdentity,
		userName,
		userPhone,
		userEmail,
	)
	if err := u.customerCreator(ctx, &customer); err != nil {
		u.log.Debug(ctx, "failed to create customer", sl.Err(err))
		return u.errorPresenter(err)
	}
	if customer.Id == appointment.TemporalCustomerId {
		err := appointment.ErrInvalidCustomerId
		u.log.Debug(ctx, "failed to create customer", sl.Err(err))
		return u.errorPresenter(err)
	}
	return u.customerPresenter(customer)
}

package appointment_use_case

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

const sendCustomerNotificationUseCaseName = "appointment_use_case.SendCustomerNotificationUseCase"

type SendCustomerNotificationUseCase[R any] struct {
	log                         *logger.Logger
	customerLoader              appointment.CustomerByIdLoader
	serviceLoader               appointment.ServiceLoader
	sender                      shared.Sender[R]
	appointmentChangedPresenter appointment.ChangedEventPresenter[R]
}

func NewSendCustomerNotificationUseCase[R any](
	log *logger.Logger,
	customerLoader appointment.CustomerByIdLoader,
	serviceLoader appointment.ServiceLoader,
	sender shared.Sender[R],
	appointmentChangedPresenter appointment.ChangedEventPresenter[R],
) *SendCustomerNotificationUseCase[R] {
	return &SendCustomerNotificationUseCase[R]{
		log:                         log.With(sl.Component(sendCustomerNotificationUseCaseName)),
		customerLoader:              customerLoader,
		serviceLoader:               serviceLoader,
		sender:                      sender,
		appointmentChangedPresenter: appointmentChangedPresenter,
	}
}

func (u *SendCustomerNotificationUseCase[R]) SendCustomerNotification(
	ctx context.Context,
	event appointment.ChangedEvent,
) {
	customer, err := u.customerLoader(ctx, event.Record.CustomerId)
	if err != nil {
		u.log.Error(ctx, "failed to load customer", sl.Err(err))
		return
	}
	service, err := u.serviceLoader(ctx, event.Record.ServiceId)
	if err != nil {
		u.log.Error(ctx, "failed to load service", sl.Err(err))
		return
	}
	notification, err := u.appointmentChangedPresenter(event, customer, service)
	if err != nil {
		u.log.Error(ctx, "failed to render notification", sl.Err(err))
		return
	}
	if err := u.sender(ctx, notification); err != nil {
		u.log.Error(ctx, "failed to send notification", sl.Err(err))
		return
	}
}

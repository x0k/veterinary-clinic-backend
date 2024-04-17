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
	sender                      shared.Sender[R]
	appointmentChangedPresenter appointment.EventPresenter[appointment.ChangedEvent, R]
}

func NewSendCustomerNotificationUseCase[R any](
	log *logger.Logger,
	sender shared.Sender[R],
	appointmentChangedPresenter appointment.EventPresenter[appointment.ChangedEvent, R],
) *SendCustomerNotificationUseCase[R] {
	return &SendCustomerNotificationUseCase[R]{
		log:                         log.With(sl.Component(sendCustomerNotificationUseCaseName)),
		sender:                      sender,
		appointmentChangedPresenter: appointmentChangedPresenter,
	}
}

func (u *SendCustomerNotificationUseCase[R]) SendCustomerNotification(
	ctx context.Context,
	event appointment.ChangedEvent,
) {
	notification, err := u.appointmentChangedPresenter(event)
	if err != nil {
		u.log.Error(ctx, "failed to render notification", sl.Err(err))
	}
	if err := u.sender(ctx, notification); err != nil {
		u.log.Error(ctx, "failed to send notification", sl.Err(err))
	}
}

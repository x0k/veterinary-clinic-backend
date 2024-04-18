package appointment_use_case

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

const sendAdminNotificationUseCaseName = "appointment_use_case.SendAdminNotificationUseCase"

type SendAdminNotificationUseCase[R any] struct {
	log                          *logger.Logger
	sender                       shared.Sender[R]
	appointmentCreatedPresenter  appointment.EventPresenter[appointment.CreatedEvent, R]
	appointmentCanceledPresenter appointment.EventPresenter[appointment.CanceledEvent, R]
}

func NewSendAdminNotificationUseCase[R any](
	log *logger.Logger,
	sender shared.Sender[R],
	appointmentCreatedPresenter appointment.EventPresenter[appointment.CreatedEvent, R],
	appointmentCanceledPresenter appointment.EventPresenter[appointment.CanceledEvent, R],
) *SendAdminNotificationUseCase[R] {
	return &SendAdminNotificationUseCase[R]{
		log:                          log.With(sl.Component(sendAdminNotificationUseCaseName)),
		sender:                       sender,
		appointmentCreatedPresenter:  appointmentCreatedPresenter,
		appointmentCanceledPresenter: appointmentCanceledPresenter,
	}
}

func sendNotification[E appointment.Event, R any](
	ctx context.Context,
	log *logger.Logger,
	sender shared.Sender[R],
	presenter appointment.EventPresenter[E, R],
	event E,
) {
	notification, err := presenter(event)
	if err != nil {
		log.Error(ctx, "failed to render notification", sl.Err(err))
		return
	}
	if err := sender(ctx, notification); err != nil {
		log.Error(ctx, "failed to send notification", sl.Err(err))
	}
}

func (u *SendAdminNotificationUseCase[R]) SendAdminNotification(ctx context.Context, event appointment.Event) {
	switch e := event.(type) {
	case appointment.CreatedEvent:
		sendNotification(ctx, u.log, u.sender, u.appointmentCreatedPresenter, e)
	case appointment.CanceledEvent:
		sendNotification(ctx, u.log, u.sender, u.appointmentCanceledPresenter, e)
	}
}

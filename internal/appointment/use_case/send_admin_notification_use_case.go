package appointment_use_case

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
)

type SendAdminNotificationUseCase[R any] struct {
	sender                       shared.Sender[R]
	appointmentCreatedPresenter  appointment.EventPresenter[appointment.AppointmentCreatedEvent, R]
	appointmentCanceledPresenter appointment.EventPresenter[appointment.AppointmentCanceledEvent, R]
}

func NewSendAdminNotificationUseCase[R any](
	sender shared.Sender[R],
	appointmentCreatedPresenter appointment.EventPresenter[appointment.AppointmentCreatedEvent, R],
	appointmentCanceledPresenter appointment.EventPresenter[appointment.AppointmentCanceledEvent, R],
) *SendAdminNotificationUseCase[R] {
	return &SendAdminNotificationUseCase[R]{
		sender:                       sender,
		appointmentCreatedPresenter:  appointmentCreatedPresenter,
		appointmentCanceledPresenter: appointmentCanceledPresenter,
	}
}

func sendNotification[E appointment.Event, R any](
	ctx context.Context,
	presenter appointment.EventPresenter[E, R],
	event E,
	sender shared.Sender[R],
) error {
	notification, err := presenter(event)
	if err != nil {
		return err
	}
	return sender(ctx, notification)
}

func (u *SendAdminNotificationUseCase[R]) SendAdminNotification(ctx context.Context, event appointment.Event) error {
	switch e := event.(type) {
	case appointment.AppointmentCreatedEvent:
		return sendNotification(ctx, u.appointmentCreatedPresenter, e, u.sender)
	case appointment.AppointmentCanceledEvent:
		return sendNotification(ctx, u.appointmentCanceledPresenter, e, u.sender)
	default:
		return nil
	}
}

package appointment_use_case

import (
	"context"

	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
)

type SendAdminNotificationUseCase[R any] struct {
	notificationSender           appointment.NotificationSender[R]
	appointmentCreatedPresenter  appointment.EventPresenter[appointment.AppointmentCreatedEvent, R]
	appointmentCanceledPresenter appointment.EventPresenter[appointment.AppointmentCanceledEvent, R]
}

func NewSendAdminNotificationUseCase[R any](
	notificationSender appointment.NotificationSender[R],
	appointmentCreatedPresenter appointment.EventPresenter[appointment.AppointmentCreatedEvent, R],
	appointmentCanceledPresenter appointment.EventPresenter[appointment.AppointmentCanceledEvent, R],
) *SendAdminNotificationUseCase[R] {
	return &SendAdminNotificationUseCase[R]{
		notificationSender:           notificationSender,
		appointmentCreatedPresenter:  appointmentCreatedPresenter,
		appointmentCanceledPresenter: appointmentCanceledPresenter,
	}
}

func sendNotification[E appointment.Event, R any](
	ctx context.Context,
	presenter appointment.EventPresenter[E, R],
	event E,
	notificationSender appointment.NotificationSender[R],
) error {
	notification, err := presenter(event)
	if err != nil {
		return err
	}
	return notificationSender.SendNotification(ctx, notification)
}

func (u *SendAdminNotificationUseCase[R]) SendAdminNotification(ctx context.Context, event appointment.Event) error {
	switch e := event.(type) {
	case appointment.AppointmentCreatedEvent:
		return sendNotification(ctx, u.appointmentCreatedPresenter, e, u.notificationSender)
	case appointment.AppointmentCanceledEvent:
		return sendNotification(ctx, u.appointmentCanceledPresenter, e, u.notificationSender)
	default:
		return nil
	}
}

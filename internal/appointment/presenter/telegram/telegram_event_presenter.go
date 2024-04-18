package appointment_telegram_presenter

import (
	"strings"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
	appointment_presenter "github.com/x0k/veterinary-clinic-backend/internal/appointment/presenter"
	"gopkg.in/telebot.v3"
)

type AppointmentCreatedEventPresenter struct {
	recipient telebot.Recipient
}

func NewAppointmentCreatedEventPresenter(recipient telebot.Recipient) AppointmentCreatedEventPresenter {
	return AppointmentCreatedEventPresenter{
		recipient: recipient,
	}
}

func (p AppointmentCreatedEventPresenter) Present(
	created appointment.CreatedEvent,
) (telegram_adapters.Message, error) {
	sb := strings.Builder{}
	sb.WriteString("*Новая запись*:\n\n")
	writeAppointmentSummary(&sb, created.Record, created.Customer, created.Service)
	return telegram_adapters.NewTextMessages(
		p.recipient,
		telegram_adapters.NewSendableText(
			sb.String(),
			&telebot.SendOptions{
				ParseMode: telebot.ModeMarkdownV2,
			},
		),
	), nil
}

type AppointmentCanceledEventPresenter struct {
	recipient telebot.Recipient
}

func NewAppointmentCanceledEventPresenter(recipient telebot.Recipient) AppointmentCanceledEventPresenter {
	return AppointmentCanceledEventPresenter{
		recipient: recipient,
	}
}

func (p AppointmentCanceledEventPresenter) Present(
	canceled appointment.CanceledEvent,
) (telegram_adapters.Message, error) {
	sb := strings.Builder{}
	sb.WriteString("*Запись отменена*:\n\n")
	writeAppointmentSummary(&sb, canceled.Record, canceled.Customer, canceled.Service)
	return telegram_adapters.NewTextMessages(
		p.recipient,
		telegram_adapters.NewSendableText(
			sb.String(),
			&telebot.SendOptions{
				ParseMode: telebot.ModeMarkdownV2,
			},
		),
	), nil
}

func writeChangeType(
	sb *strings.Builder,
	changeType appointment.ChangeType,
) {
	switch changeType {
	case appointment.CreatedChangeType:
		sb.WriteString("*Создана запись*")
	case appointment.StatusChangeType:
		sb.WriteString("*Статус изменен*")
	case appointment.DateTimeChangeType:
		sb.WriteString("*Дата и время изменены*")
	case appointment.RemovedChangeType:
		sb.WriteString("*Запись удалена*")
	}
}

func AppointmentChangedEventPresenter(
	event appointment.ChangedEvent,
) (telegram_adapters.Message, error) {
	id, err := event.Customer.Identity.ToTelegramUserId()
	if err != nil {
		return nil, err
	}

	sb := strings.Builder{}
	writeChangeType(&sb, event.ChangeType)
	sb.WriteString(":\n\n")

	state, err := appointment_presenter.RecordState(event.Record.Status, event.Record.IsArchived)
	if err != nil {
		return nil, err
	}
	sb.WriteString(state)
	sb.WriteString("\n\n")

	writeAppointmentSummary(&sb, event.Record, event.Customer, event.Service)

	return telegram_adapters.NewTextMessages(
		&telebot.User{
			ID: id.Int(),
		},
		telegram_adapters.NewSendableText(
			sb.String(),
			&telebot.SendOptions{
				ParseMode: telebot.ModeMarkdownV2,
			},
		),
	), nil
}

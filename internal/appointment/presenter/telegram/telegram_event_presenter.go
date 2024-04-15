package appointment_telegram_presenter

import (
	"strings"

	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	"github.com/x0k/veterinary-clinic-backend/internal/appointment"
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
	writeAppointmentSummary(&sb, created.AppointmentAggregate)
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
	writeAppointmentSummary(&sb, canceled.AppointmentAggregate)
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

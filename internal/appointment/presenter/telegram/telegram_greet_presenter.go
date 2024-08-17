package appointment_telegram_presenter

import (
	telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/adapters/telegram"
	appointment_telegram_adapters "github.com/x0k/veterinary-clinic-backend/internal/appointment/adapters/telegram"
	"gopkg.in/telebot.v3"
)

type GreetingPresenter struct {
	buttons [][]telebot.InlineButton
}

func NewGreetingPresenter(
	createAppointment bool,
) *GreetingPresenter {
	buttons := [][]telebot.InlineButton{
		{*appointment_telegram_adapters.ScheduleBtn},
		{*appointment_telegram_adapters.ServicesBtn},
	}
	if createAppointment {
		buttons = append(buttons, []telebot.InlineButton{*appointment_telegram_adapters.StartMakeAppointmentDialogBtn})
	}
	return &GreetingPresenter{
		buttons: buttons,
	}
}

func (p *GreetingPresenter) RenderGreeting() (telegram_adapters.TextResponses, error) {
	return telegram_adapters.TextResponses{{
		Text: telegram_adapters.EscapeMarkdownString("Привет!"),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: p.buttons,
			},
		},
	}}, nil
}

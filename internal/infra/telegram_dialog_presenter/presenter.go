package telegram_dialog_presenter

import (
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"gopkg.in/telebot.v3"
)

type Config struct {
	CalendarWebAppUrl string
}

type TelegramDialogPresenter struct {
	calendarResponse shared.TelegramResponse
}

func New(cfg *Config) *TelegramDialogPresenter {
	calendarKeyboard := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{{
			{
				Text: "Открыть календарь",
				WebApp: &telebot.WebApp{
					URL: cfg.CalendarWebAppUrl,
				},
			},
		}},
	}
	return &TelegramDialogPresenter{
		calendarResponse: shared.TelegramResponse{
			Text: "Выберите дату",
			Options: &telebot.SendOptions{
				ReplyMarkup: calendarKeyboard,
			},
		},
	}
}

func (p *TelegramDialogPresenter) RenderGreeting() (shared.TelegramResponse, error) {
	return shared.TelegramResponse{
		Text: shared.EscapeTelegramMarkdownString("Привет!"),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
		},
	}, nil
}

func (p *TelegramDialogPresenter) RenderScheduleDialog(dialog entity.Dialog) (shared.TelegramResponse, error) {
	return p.calendarResponse, nil
}

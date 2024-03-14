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
	datePickerResponse shared.TelegramTextResponse
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
		datePickerResponse: shared.TelegramTextResponse{
			Text: "Выберите дату",
			Options: &telebot.SendOptions{
				ReplyMarkup: calendarKeyboard,
			},
		},
	}
}

func (p *TelegramDialogPresenter) RenderGreeting() (shared.TelegramResponse, error) {
	return shared.TelegramTextResponse{
		Text: shared.EscapeTelegramMarkdownString("Привет!"),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
		},
	}, nil
}

func (p *TelegramDialogPresenter) RenderDatePicker() (shared.TelegramResponse, error) {
	return p.datePickerResponse, nil
}

func (p *TelegramDialogPresenter) RenderSchedule(periods []entity.TimePeriod) (shared.TelegramResponse, error) {
	return shared.TelegramQueryResponse{}, nil
}

func (p *TelegramDialogPresenter) RenderError(err error) (shared.TelegramResponse, error) {
	return shared.TelegramTextResponse{
		Text: shared.EscapeTelegramMarkdownString(err.Error()),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
		},
	}, nil
}

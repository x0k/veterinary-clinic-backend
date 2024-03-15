package presenter

import (
	"github.com/x0k/veterinary-clinic-backend/internal/adapters"
	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"gopkg.in/telebot.v3"
)

type TelegramDialogConfig struct {
	CalendarWebAppUrl string
}

type TelegramDialogPresenter struct {
	datePickerResponse adapters.TelegramTextResponse
}

func NewTelegramDialog(cfg *TelegramDialogConfig) *TelegramDialogPresenter {
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
		datePickerResponse: adapters.TelegramTextResponse{
			Text: "Выберите дату",
			Options: &telebot.SendOptions{
				ReplyMarkup: calendarKeyboard,
			},
		},
	}
}

func (p *TelegramDialogPresenter) RenderGreeting() (adapters.TelegramResponse, error) {
	return adapters.TelegramTextResponse{
		Text: adapters.EscapeTelegramMarkdownString("Привет!"),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
		},
	}, nil
}

func (p *TelegramDialogPresenter) RenderDatePicker() (adapters.TelegramResponse, error) {
	return p.datePickerResponse, nil
}

func (p *TelegramDialogPresenter) RenderSchedule(periods []entity.TimePeriod) (adapters.TelegramResponse, error) {
	return adapters.TelegramQueryResponse{}, nil
}

func (p *TelegramDialogPresenter) RenderError(err error) (adapters.TelegramResponse, error) {
	return adapters.TelegramTextResponse{
		Text: adapters.EscapeTelegramMarkdownString(err.Error()),
		Options: &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
		},
	}, nil
}

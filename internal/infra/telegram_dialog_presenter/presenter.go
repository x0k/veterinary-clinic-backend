package telegram_dialog_presenter

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/x0k/veterinary-clinic-backend/internal/entity"
	"github.com/x0k/veterinary-clinic-backend/internal/shared"
	"gopkg.in/telebot.v3"
)

type TelegramDialogPresenterConfig struct {
	CalendarHandlerUrl string
}

type TelegramDialogPresenter struct {
	calendarResponse shared.TelegramResponse
}

type CalendarRequestOptions struct {
	Url string `json:"url"`
}

func New(cfg *TelegramDialogPresenterConfig) (*TelegramDialogPresenter, error) {
	const op = "infra.telegram_dialog_presenter.New"
	params := url.Values{}
	reqOptions, err := json.Marshal(&CalendarRequestOptions{
		Url: cfg.CalendarHandlerUrl,
	})
	if err != nil {
		return nil, fmt.Errorf("%s request options marshaling: %w", op, err)
	}
	params.Add("req", string(reqOptions))
	webAppUrl := fmt.Sprintf("%s?%s", "https://x0k.github.io/telegram-web-inputs/calendar", params.Encode())
	calendarKeyboard := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{{
			{
				Text: "Открыть календарь",
				WebApp: &telebot.WebApp{
					URL: webAppUrl,
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
	}, nil
}

func (p *TelegramDialogPresenter) RenderScheduleDialog(dialog entity.Dialog) (shared.TelegramResponse, error) {
	return p.calendarResponse, nil
}

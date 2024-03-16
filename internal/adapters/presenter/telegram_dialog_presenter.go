package presenter

import (
	"fmt"
	"strings"

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

func (p *TelegramDialogPresenter) RenderSchedule(schedule entity.Schedule) (adapters.TelegramResponse, error) {
	sb := strings.Builder{}
	sb.WriteString("График работы на ")
	sb.WriteString(adapters.EscapeTelegramMarkdownString(
		entity.DateToGoTime(schedule.Date).Format("02.01.2006")),
	)
	sb.WriteString(":\n\n")
	for _, period := range schedule.Periods {
		sb.WriteByte('*')
		sb.WriteString(period.Start.String())
		sb.WriteString(" \\- ")
		sb.WriteString(period.End.String())
		sb.WriteString("*\n")
		sb.WriteString(adapters.EscapeTelegramMarkdownString(period.Title))
		sb.WriteString("\n\n")
	}
	if len(schedule.Periods) == 0 {
		sb.WriteString("Нет записей\n\n")
	}
	return adapters.TelegramQueryResponse{
		Result: &telebot.ArticleResult{
			ResultBase: telebot.ResultBase{
				ID:        fmt.Sprintf("%p", &schedule),
				Type:      "article",
				ParseMode: telebot.ModeMarkdownV2,
			},
			Title: "График работы",
			Text:  sb.String(),
		},
	}, nil
}

func (p *TelegramDialogPresenter) RenderError(err error) (adapters.TelegramResponse, error) {
	return adapters.TelegramQueryResponse{
		Result: &telebot.ArticleResult{
			ResultBase: telebot.ResultBase{
				ID:        fmt.Sprintf("%p", &err),
				Type:      "article",
				ParseMode: telebot.ModeMarkdownV2,
			},
			Title: "Ошибка",
			Text:  adapters.EscapeTelegramMarkdownString(err.Error()),
		},
	}, nil
}
